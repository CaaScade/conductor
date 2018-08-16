package controller

import (
	"fmt"

	"github.com/koki/conductor/app"
	role_model "github.com/koki/conductor/app/src/roles/models"
	"github.com/revel/revel"
	"github.com/koki/conductor/app/src/util"
)

type Permission struct {
	*revel.Controller
}

func (p *Permission) ListPermissions() revel.Result {
	perms := []role_model.Permission{}
	app.DB.Find(&perms)
	return p.RenderJSON(perms)
}

func (p *Permission) CreatePermission() revel.Result {
	permType := role_model.Permission{}
	p.Params.BindJSON(&permType)
	if !app.DB.Where("name = ?", permType.Name).First(&permType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("perm already exists")}
	}
	app.DB.Create(&permType)
	return p.RenderJSON(permType)
}

func (p *Permission) GetPermission(perm string) revel.Result {
	permType := role_model.Permission{Name: perm}
	if app.DB.Where(&permType).First(&permType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown perm")}
	}
	return p.RenderJSON(permType)
}

func (p *Permission) UpdatePermission(perm string) revel.Result {
	permType := role_model.Permission{Name: perm}
	if app.DB.Where(&permType).First(&permType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown perm")}
	}
	p.Params.BindJSON(&permType)
	app.DB.Model(&role_model.Permission{}).Updates(&permType)
	return p.RenderJSON(permType)
}

func (p *Permission) DeletePermission(perm string) revel.Result {
	permType := role_model.Permission{Name: perm}
	if app.DB.Where(&permType).First(&permType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown perm")}
	}
	app.DB.Model(&role_model.Permission{}).Delete(&permType)
	return util.AppResponse{200, "" , nil}
}

func (p *Permission) GetRoles(perm string) revel.Result {
	permType := role_model.Permission{Name: perm}
	if app.DB.Where(&permType).First(&permType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown perm")}
	}
	roles := new([]role_model.Role)
	app.DB.Model(&permType).Related(&roles, "Roles")
	return p.RenderJSON(roles)
}

func (p *Permission) SetRoles(perm string) revel.Result {
	permType := role_model.Permission{}
	if app.DB.Where("name = ?", perm).First(&permType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	roles := new([]role_model.Role)
	p.Params.BindJSON(&roles)
	app.DB.Model(&permType).Association("Roles").Clear()
	app.DB.Model(&permType).Association("Roles").Append(roles)
	app.DB.Model(&permType).Related(&roles, "Roles")
	return p.RenderJSON(roles)
}

func (p *Permission) AddRoles(perm string) revel.Result {
	permType := role_model.Permission{}
	if app.DB.Where("name = ?", perm).First(&permType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown role")}
	}
	roles := new([]role_model.Role)
	p.Params.BindJSON(&roles)
	app.DB.Model(&permType).Association("Roles").Append(roles)
	app.DB.Model(&permType).Related(&roles, "Roles")
	return p.RenderJSON(roles)
}
