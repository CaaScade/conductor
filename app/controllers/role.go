package controllers

import (
	"fmt"

	"github.com/koki/conductor/app"
	"github.com/koki/conductor/app/models"
	"github.com/revel/revel"
)

type Role struct {
	*revel.Controller
}

func (r *Role) ListRoles() revel.Result {
	roles := []models.Role{}
	app.DB.Find(&roles)
	return r.RenderJSON(roles)
}

func (r *Role) CreateRole() revel.Result {
	roleType := models.Role{}
	r.Params.BindJSON(&roleType)
	if !app.DB.Where(&roleType).Find(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("role already exists")}
	}
	app.DB.Create(&roleType)
	return r.RenderJSON(roleType)
}

func (r *Role) GetRole(role string) revel.Result {
	roleType := models.Role{Name: role}
	if app.DB.Where(&roleType).Find(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	return r.RenderJSON(roleType)
}

func (r *Role) UpdateRole(role string) revel.Result {
	roleType := models.Role{Name: role}
	if app.DB.Where(&roleType).First(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	r.Params.BindJSON(&roleType)
	app.DB.Model(&models.Role{}).Updates(&roleType)
	return r.RenderJSON(roleType)
}

func (r *Role) DeleteRole(role string) revel.Result {
	roleType := models.Role{Name: role}
	if app.DB.Where(&roleType).First(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	app.DB.Model(&models.Role{}).Delete(&roleType)
	return RenderStatus{200, ""}
}

func (r *Role) GetUsers(role string) revel.Result {
	roleType := models.Role{Name: role}
	if app.DB.Where(&roleType).First(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	users := new([]models.User)
	app.DB.Model(&roleType).Related(&users, "Users")
	return r.RenderJSON(users)
}

func (r *Role) GetPerms(role string) revel.Result {
	roleType := models.Role{}
	if app.DB.Where("name = ?", role).First(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	perms := new([]models.Permission)
	app.DB.Model(&roleType).Related(&perms, "Permissions")
	return r.RenderJSON(perms)
}

func (r *Role) SetPerms(role string) revel.Result {
	roleType := models.Role{}
	if app.DB.Where("name = ?", role).First(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	perms := new([]models.Permission)
	r.Params.BindJSON(&perms)
	revel.AppLog.Debugf("%+v", perms)
	app.DB.Model(&roleType).Association("Permissions").Clear()
	app.DB.Model(&roleType).Association("Permissions").Append(perms)
	app.DB.Model(&roleType).Related(&perms, "Permissions")
	return r.RenderJSON(perms)
}

func (r *Role) AddPerms(role string) revel.Result {
	roleType := models.Role{}
	if app.DB.Where("name = ?", role).First(&roleType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	perms := new([]models.Permission)
	r.Params.BindJSON(&perms)
	app.DB.Model(&roleType).Association("Permissions").Append(perms)
	app.DB.Model(&roleType).Related(&perms, "Permissions")
	return r.RenderJSON(perms)
}
