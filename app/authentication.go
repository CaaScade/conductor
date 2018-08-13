package app

import (
	"crypto/rsa"
	"encoding/base32"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/koki/conductor/app/models"
	"github.com/pquerna/otp/hotp"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

const AUTHORIZATION_HEADER = "Authorization"
const AUTHENTICATED_CONF = "koki.app.authenticated"
const USER_ID = "user_id"
const OTP = "otp"
const ROLES = "roles"
const PERMS = "perms"
const AUTH_LOGIN_PATH = "/auth/login"
const AUTH_REGISTER_PATH = "/auth/register"
const AUTH_LOGOUT_PATH = "/auth/logout"
const UI_PATH = "/ui"

var AuthCounter map[string]uint64

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func init() {
	AuthCounter = map[string]uint64{}
}

func AuthFilter(c *revel.Controller, fc []revel.Filter) {
	if !revel.Config.BoolDefault(AUTHENTICATED_CONF, false) {
		fc[0](c, fc[1:])
		return
	}
	//appSecret := revel.Config.StringDefault("app.secret", "dummy_secret")
	if strings.Index(strings.TrimRight(c.Request.GetPath(), "/"), "/ui") == 0 {
		fc[0](c, fc[1:])
		return
	}
	redirect := c.Params.Query.Get("redirect")
	if strings.TrimRight(c.Request.GetPath(), "/") == AUTH_LOGIN_PATH {
		var user models.User
		username := c.Params.Form.Get("username")
		if DB.Where(&models.User{Username: username}).First(&user).RecordNotFound() {
			c.Response.SetStatus(http.StatusUnauthorized)
			c.Response.GetWriter().Write([]byte("username or password incorrect"))
			return
		}
		password := c.Params.Form.Get("password")
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			c.Response.SetStatus(http.StatusUnauthorized)
			c.Response.GetWriter().Write([]byte("username or password incorrect"))
			return
		}
		roles := new([]models.Role)
		perms := new([]models.Permission)
		/*
			DB.Where(&user).Related(&roles)
			newPerms := new([]models.Permission)
			for _, r := range *roles {
				DB.Model(&r).Related(&newPerms, "Permissions")
				*perms = append(*perms, *newPerms...)
			}
		*/
		secretUsername := &strings.Builder{}
		encoder := base32.NewEncoder(base32.StdEncoding, secretUsername)
		encoder.Write([]byte(username))
		defer encoder.Close()
		otp, err := hotp.GenerateCode(secretUsername.String(), user.Counter)
		if err != nil {
			c.Response.SetStatus(http.StatusInternalServerError)
			c.Response.GetWriter().Write([]byte("error generating one time password" + err.Error()))
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			USER_ID: username,
			OTP:     otp,
			ROLES:   roles,
			PERMS:   perms,
			"nbf":   time.Now().Unix(),
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString(signKey)
		if err != nil {
			c.Response.SetStatus(http.StatusInternalServerError)
			c.Response.GetWriter().Write([]byte("error generating token " + err.Error()))
			return
		}
		c.Response.Out.Header().Set(AUTHORIZATION_HEADER, fmt.Sprintf("Bearer %s", tokenString))
		AuthCounter[username] = uint64(user.Counter)
		if redirect != "" {
			c.Redirect(redirect)
			return
		}
		c.Response.SetStatus(http.StatusOK)
		return
	}

	var username string
	var otp string
	if strings.TrimRight(c.Request.GetPath(), "/") == AUTH_REGISTER_PATH {
		var user models.User
		username := c.Params.Form.Get("username")
		if !DB.Where(&models.User{Username: username}).First(&user).RecordNotFound() {
			c.Response.SetStatus(http.StatusBadRequest)
			c.Response.GetWriter().Write([]byte("username already exists"))
			return
		}
		password := c.Params.Form.Get("password")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.Response.SetStatus(http.StatusInternalServerError)
			c.Response.GetWriter().Write([]byte(err.Error()))
			return
		}
		email := c.Params.Form.Get("email")
		user.Username = username
		user.Password = string(hashedPassword)
		user.Email = email
		user.Counter = 1
		DB.Create(&user)
		c.Response.SetStatus(http.StatusOK)
		if redirect != "" {
			c.Redirect(redirect)
			return
		}
		return
	}
	headerData := c.Request.Header.Get(AUTHORIZATION_HEADER)
	headers := strings.Split(headerData, " ")
	if len(headers) != 2 {
		c.Response.SetStatus(http.StatusUnauthorized)
		c.Response.GetWriter().Write([]byte(AUTHORIZATION_HEADER + " format invalid"))
		return
	}
	if headers[0] != "Bearer" {
		c.Response.SetStatus(http.StatusUnauthorized)
		c.Response.GetWriter().Write([]byte(AUTHORIZATION_HEADER + " format invalid; 'Bearer' missing"))
		return
	}
	tokenString := headers[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return verifyKey, nil
	})
	if err != nil {
		c.Response.SetStatus(http.StatusUnauthorized)
		c.Response.GetWriter().Write([]byte(err.Error()))
		return
	}
	if !token.Valid {
		c.Response.SetStatus(http.StatusUnauthorized)
		c.Response.GetWriter().Write([]byte("invalid token"))
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userId := claims[USER_ID]
		if uname, ok := userId.(string); !ok {
			c.Response.SetStatus(http.StatusUnauthorized)
			c.Response.GetWriter().Write([]byte("no username found in claim"))
			return
		} else {
			username = uname
		}
		otpInt := claims[OTP]
		if p, ok := otpInt.(string); !ok {
			c.Response.SetStatus(http.StatusUnauthorized)
			c.Response.GetWriter().Write([]byte("OTP not found"))
			return
		} else {
			otp = p
		}
		secretUsername := &strings.Builder{}
		encoder := base32.NewEncoder(base32.StdEncoding, secretUsername)
		encoder.Write([]byte(username))
		defer encoder.Close()
		if AuthCounter[username] == 0 {
			var user models.User
			if DB.Where(&models.User{Username: username}).First(&user).RecordNotFound() {
				c.Response.SetStatus(http.StatusUnauthorized)
				c.Response.GetWriter().Write([]byte("Username is invalid"))
			}
			AuthCounter[username] = user.Counter
		}
		if !hotp.Validate(otp, AuthCounter[username], secretUsername.String()) {
			c.Response.SetStatus(http.StatusUnauthorized)
			c.Response.GetWriter().Write([]byte("OTP is invalid"))
			return
		}
		c.Args[USER_ID] = username
		c.Args[ROLES] = claims[ROLES]
		c.Args[PERMS] = claims[PERMS]
	} else {
		c.Response.SetStatus(http.StatusUnauthorized)
		c.Response.GetWriter().Write([]byte("invalid token claim"))
		return
	}
	if strings.TrimRight(c.Request.GetPath(), "/") == AUTH_LOGOUT_PATH {
		var user models.User
		if DB.Where(&models.User{Username: username}).First(&user).RecordNotFound() {
			c.Response.SetStatus(http.StatusBadRequest)
			c.Response.GetWriter().Write([]byte("invalid username"))
			return
		}
		user.Counter = user.Counter + 1
		AuthCounter[username] = AuthCounter[username] + 1
		DB.Model(&user).Update(user)
		c.Response.SetStatus(http.StatusOK)
		return
	}

	fc[0](c, fc[1:])
}
