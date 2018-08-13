package controllers

import (
	"fmt"

	"github.com/koki/conductor/app"
	"github.com/koki/conductor/app/models"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	*revel.Controller
}

type RenderStatus struct {
	Status int
	Text   string
}

func (r RenderStatus) Apply(req *revel.Request, resp *revel.Response) {
	resp.SetStatus(r.Status)
	resp.GetWriter().Write([]byte(r.Text))
}

func (u *User) ListUsers() revel.Result {
	users := []models.User{}
	app.DB.Find(&users)
	return u.RenderJSON(users)
}

func (u *User) GetUser(username string) revel.Result {
	user := models.User{Username: username}
	if app.DB.Where(&user).Find(&user).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown username")}
	}
	return u.RenderJSON(user)
}

func (u *User) UpdateUser(username string) revel.Result {
	user := models.User{Username: username}
	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown username")}
	}
	counter := user.Counter
	u.Params.BindJSON(&user)
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return RenderStatus{500, err.Error()}
		}
		user.Password = string(hashedPassword)
	}
	app.AuthCounter[username] = app.AuthCounter[username] + 1
	user.Counter = counter
	app.DB.Model(&models.User{}).Updates(&user)
	return u.RenderJSON(user)
}

func (u *User) DeleteUser(username string) revel.Result {
	user := models.User{Username: username}
	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown username")}
	}
	app.DB.Model(&models.User{}).Delete(&user)
	app.AuthCounter[username] = app.AuthCounter[username] + 1
	return RenderStatus{200, ""}
}

func (u *User) GetRoles(username string) revel.Result {
	user := models.User{Username: username}
	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown username")}
	}
	roles := new([]models.Role)
	//app.DB.Preload("Roles").First(&user)
	app.DB.Model(&user).Related(&roles, "Roles")
	return u.RenderJSON(roles)
}
