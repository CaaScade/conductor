package user

import (
	"fmt"

	"github.com/koki/conductor/app"
	user_model "github.com/koki/conductor/app/src/user/models"
	role_model "github.com/koki/conductor/app/src/roles/models"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"github.com/koki/conductor/app/src/util"
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
	users := []user_model.User{}
	app.DB.Find(&users)
	return u.RenderJSON(users)
}

func (u *User) GetUser(username string) revel.Result {
	user := user_model.User{Username: username}
	if app.DB.Where(&user).Find(&user).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown username")}
	}
	return util.AppResponse{200, "Success" ,user}
}

func (u *User) UpdateUser(username string) revel.Result {
	user := user_model.User{Username: username}
	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown username")}
	}
	counter := user.Counter
	u.Params.BindJSON(&user)
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return util.AppResponse{500, err.Error(), nil}
		}
		user.Password = string(hashedPassword)
	}
	app.AuthCounter[username] = app.AuthCounter[username] + 1
	user.Counter = counter
	app.DB.Model(&user_model.User{}).Updates(&user)
	return util.AppResponse{200, "success" ,user}
}

func (u *User) DeleteUser(username string) revel.Result {
	user := user_model.User{Username: username}
	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return util.AppResponse{400, "unknown username" ,nil}
	}
	app.DB.Model(&user_model.User{}).Delete(&user)
	app.AuthCounter[username] = app.AuthCounter[username] + 1
	return util.AppResponse{200, "" ,nil}
}

func (u *User) GetRoles(username string) revel.Result {
	user := user_model.User{Username: username}
	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return util.AppResponse{400, "unknown username" ,nil}
	}
	roles := new([]role_model.Role)
	app.DB.Model(&user).Related(&roles, "Roles")
	app.DB.Preload("Roles").First(&user)
	return  util.AppResponse{200, "Success", roles}
}

func (u *User) SetRoles(username string) revel.Result {
	user := user_model.User{Username: username}
	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return util.AppResponse{400, "user not found" ,nil}
	}
	roles := new([]role_model.Role)
	u.Params.BindJSON(&roles)
	app.DB.Model(&user).Association("Roles").Clear()
	app.DB.Model(&user).Association("Roles").Append(roles)
	app.DB.Model(&user).Related(&roles, "Roles")
	return  util.AppResponse{200, "Success", roles}
}

func (u *User) AddRole(username string) revel.Result {
	user := user_model.User{Username: username}
	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return util.AppResponse{400, "user not found", nil}
	}
	roles := new([]role_model.Role)
	u.Params.BindJSON(&roles)
	app.DB.Model(&user).Association("Roles").Append(roles)
	app.DB.Model(&user).Related(&roles, "Roles")


	return  util.AppResponse{200, "Success", roles}
}