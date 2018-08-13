package controllers

import (
	"fmt"

	"github.com/koki/conductor/app"
	"github.com/koki/conductor/app/models"
	"github.com/revel/revel"
)

type Permission struct {
	*revel.Controller
}

func (p *Permission) ListPermissions() revel.Result {
	perms := []models.Permission{}
	app.DB.Find(&perms)
	return p.RenderJSON(perms)
}

func (p *Permission) CreatePermission() revel.Result {
	permType := models.Permission{}
	p.Params.BindJSON(&permType)
	if !app.DB.Where("name = ?", permType.Name).First(&permType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("perm already exists")}
	}
	app.DB.Create(&permType)
	return p.RenderJSON(permType)
}

func (p *Permission) GetPermission(perm string) revel.Result {
	permType := models.Permission{Name: perm}
	if app.DB.Where(&permType).First(&permType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown perm")}
	}
	return p.RenderJSON(permType)
}

func (p *Permission) UpdatePermission(perm string) revel.Result {
	permType := models.Permission{Name: perm}
	if app.DB.Where(&permType).First(&permType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown perm")}
	}
	p.Params.BindJSON(&permType)
	app.DB.Model(&models.Permission{}).Updates(&permType)
	return p.RenderJSON(permType)
}

func (p *Permission) DeletePermission(perm string) revel.Result {
	permType := models.Permission{Name: perm}
	if app.DB.Where(&permType).First(&permType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown perm")}
	}
	app.DB.Model(&models.Permission{}).Delete(&permType)
	return RenderStatus{200, ""}
}

func (p *Permission) GetRoles(perm string) revel.Result {
	permType := models.Permission{Name: perm}
	if app.DB.Where(&permType).First(&permType).RecordNotFound() {
		return revel.PlaintextErrorResult{Error: fmt.Errorf("unknown perm")}
	}
	roles := new([]models.Role)
	app.DB.Model(&permType).Related(&roles, "Roles")
	return p.RenderJSON(roles)
}
