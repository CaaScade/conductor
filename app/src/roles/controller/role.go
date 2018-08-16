package controller

import (
	"fmt"

	"github.com/koki/conductor/app"
	role_model "github.com/koki/conductor/app/src/roles/models"
	"github.com/revel/revel"
	"github.com/koki/conductor/app/src/util"
	user_model "github.com/koki/conductor/app/src/user/models"
)

type Role struct {
	*revel.Controller
}

func (r *Role) ListRoles() revel.Result {
	roles := []role_model.Role{}
	app.DB.Find(&roles)
	return r.RenderJSON(roles)
}

func (r *Role) CreateRole() revel.Result {
	roleType := role_model.Role{}
	r.Params.BindJSON(&roleType)
	if !app.DB.Where(&roleType).Find(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("role already exists")}
	}
	app.DB.Create(&roleType)
	return r.RenderJSON(roleType)
}

func (r *Role) GetRole(role string) revel.Result {
	roleType := role_model.Role{Name: role}
	if app.DB.Where(&roleType).Find(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	return r.RenderJSON(roleType)
}

func (r *Role) UpdateRole(role string) revel.Result {
	roleType := role_model.Role{Name: role}
	if app.DB.Where(&roleType).First(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	r.Params.BindJSON(&roleType)
	app.DB.Model(&role_model.Role{}).Updates(&roleType)
	return r.RenderJSON(roleType)
}

func (r *Role) DeleteRole(role string) revel.Result {
	roleType := role_model.Role{Name: role}
	if app.DB.Where(&roleType).First(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	app.DB.Model(&role_model.Role{}).Delete(&roleType)
	return util.AppResponse{200, "", nil}
}

func (r *Role) GetUsers(role string) revel.Result {
	roleType := role_model.Role{Name: role}
	if app.DB.Where(&roleType).First(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	users := new([]user_model.User)
	app.DB.Model(&roleType).Related(&users, "Users")
	return r.RenderJSON(users)
}

func (r *Role) GetPerms(role string) revel.Result {
	roleType := role_model.Role{}
	if app.DB.Where("name = ?", role).First(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	perms := new([]role_model.Permission)
	app.DB.Model(&roleType).Related(&perms, "Permissions")
	return r.RenderJSON(perms)
}

func (r *Role) SetPerms(role string) revel.Result {
	roleType := role_model.Role{}
	if app.DB.Where("name = ?", role).First(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	perms := new([]role_model.Permission)
	r.Params.BindJSON(&perms)
	revel.AppLog.Debugf("%+v", perms)
	app.DB.Model(&roleType).Association("Permissions").Clear()
	app.DB.Model(&roleType).Association("Permissions").Append(perms)
	app.DB.Model(&roleType).Related(&perms, "Permissions")
	return r.RenderJSON(perms)
}

func (r *Role) AddPerms(role string) revel.Result {
	roleType := role_model.Role{}
	if app.DB.Where("name = ?", role).First(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	perms := new([]role_model.Permission)
	r.Params.BindJSON(&perms)
	app.DB.Model(&roleType).Association("Permissions").Append(perms)
	app.DB.Model(&roleType).Related(&perms, "Permissions")
	return r.RenderJSON(perms)
}

func (r *Role) AddUsers(role string) revel.Result {
	roleType := role_model.Role{}

	if app.DB.Where("name = ?", role).First(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}

	users := new([]user_model.User)
	r.Params.BindJSON(&users)
	app.DB.Model(&roleType).Association("Users").Append(users)
	app.DB.Model(&roleType).Related(&users, "Users")
	return r.RenderJSON(users)
}

func (r *Role) SetUsers(role string) revel.Result {
	roleType := role_model.Role{}
	if app.DB.Where("name = ?", role).First(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	users := new([]user_model.User)
	r.Params.BindJSON(&users)
	app.DB.Model(&roleType).Association("Users").Clear()
	app.DB.Model(&roleType).Association("Users").Append(users)
	app.DB.Model(&roleType).Related(&users, "Users")
	return r.RenderJSON(users)
}
