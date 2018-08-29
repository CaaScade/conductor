package controller

import (
	"github.com/koki/conductor/app"
	"github.com/koki/conductor/app/src/user/models"
	"github.com/koki/conductor/app/src/util"
	"github.com/revel/revel"
)

type Permission struct {
	*revel.Controller
}

// Fetch list of all the permission we have in the system
// TODO: provide list only if the permission is ADMIN; need to verify
func (p *Permission) ListPermissions() revel.Result {
	perms := []models.Permission{}
	app.DB.Find(&perms)
	return util.AppResponse{200, "success", perms}
}

// Create new permission
// TODO: create only if the permission is ADMIN; need to verify
func (p *Permission) CreatePermission() revel.Result {
	permType := models.Permission{}
	p.Params.BindJSON(&permType)
	if !app.DB.Where("name = ?", permType.Name).First(&permType).RecordNotFound() {
		return util.AppResponse{400, "perm already exists", nil}
	}
	app.DB.Create(&permType)
	return util.AppResponse{200, "success", permType}
}

// Get details of particular permission
// TODO: provide list only if the role is ADMIN; need to verify
func (p *Permission) GetPermission(perm string) revel.Result {
	permType := models.Permission{Name: perm}
	if app.DB.Where(&permType).First(&permType).RecordNotFound() {
		return util.AppResponse{400, "unknown perm", nil}
	}
	return util.AppResponse{200, "success", permType}
}

// Update the details of particular permission
// TODO: update only if the permission is ADMIN; need to verify
func (p *Permission) UpdatePermission(perm string) revel.Result {
	permType := models.Permission{Name: perm}
	if app.DB.Where(&permType).First(&permType).RecordNotFound() {
		return util.AppResponse{400, "unknown perm", nil}
	}
	p.Params.BindJSON(&permType)
	app.DB.Model(&models.Permission{}).Updates(&permType)
	return util.AppResponse{200, "success", permType}
}

// Delete the particular permission
// TODO: delete only if the permission is ADMIN; need to verify
func (p *Permission) DeletePermission(perm string) revel.Result {
	permType := models.Permission{Name: perm}
	if app.DB.Where(&permType).First(&permType).RecordNotFound() {
		return util.AppResponse{400, "unknown perm", nil}
	}
	app.DB.Model(&models.Permission{}).Delete(&permType)
	return util.AppResponse{200, "success", permType}
}

// Get all roles belongs to the particular permission
func (p *Permission) GetRoles(perm string) revel.Result {
	permType := models.Permission{Name: perm}
	if app.DB.Where(&permType).First(&permType).RecordNotFound() {
		return util.AppResponse{400, "unknown perm", nil}
	}
	roles := new([]models.Role)
	app.DB.Model(&permType).Related(&roles, "Roles")
	return util.AppResponse{200, "success", roles}
}

// Set new roles to the particular permission
func (p *Permission) SetRoles(perm string) revel.Result {
	permType := models.Permission{}
	if app.DB.Where("name = ?", perm).First(&permType).RecordNotFound() {
		return util.AppResponse{400, "unknown role", nil}
	}
	roles := new([]models.Role)
	p.Params.BindJSON(&roles)
	app.DB.Model(&permType).Association("Roles").Clear()
	app.DB.Model(&permType).Association("Roles").Append(roles)
	app.DB.Model(&permType).Related(&roles, "Roles")
	return util.AppResponse{200, "success", roles}
}

// Add/append new role(s) to the particular permission
func (p *Permission) AddRoles(perm string) revel.Result {
	permType := models.Permission{}
	if app.DB.Where("name = ?", perm).First(&permType).RecordNotFound() {
		return util.AppResponse{400, "unknown role", nil}
	}
	roles := new([]models.Role)
	p.Params.BindJSON(&roles)
	app.DB.Model(&permType).Association("Roles").Append(roles)
	app.DB.Model(&permType).Related(&roles, "Roles")
	//return p.RenderJSON(roles)
	return util.AppResponse{200, "success", roles}
}
