package role

import (
	"github.com/koki/conductor/app"
	"github.com/koki/conductor/app/src/user/models"
	"github.com/koki/conductor/app/src/util"
	"github.com/revel/revel"
)

type Role struct {
	*revel.Controller
}

// Get List of all roles
// TODO: provide list only if the role is ADMIN; need to verify
func (r *Role) ListRoles() revel.Result {
	roles := []models.Role{}
	app.DB.Find(&roles)
	return util.AppResponse{200, "success", roles}
}

// Create new role
// TODO: create only if the role is ADMIN; need to verify
func (r *Role) CreateRole() revel.Result {
	roleType := models.Role{}
	r.Params.BindJSON(&roleType)
	if !app.DB.Where(&roleType).Find(&roleType).RecordNotFound() {
		return util.AppResponse{400, "role already exists", nil}
	}
	app.DB.Create(&roleType)
	return util.AppResponse{200, "success", roleType}
}

// Get the details of the role
func (r *Role) GetRole(role string) revel.Result {
	roleType := models.Role{Name: role}
	if app.DB.Where(&roleType).Find(&roleType).RecordNotFound() {
		return util.AppResponse{400, "unknown role", nil}
	}
	return util.AppResponse{200, "success", roleType}
}

// Update role
// TODO: update only if the role is ADMIN; need to verify
func (r *Role) UpdateRole(role string) revel.Result {
	roleType := models.Role{Name: role}
	if app.DB.Where(&roleType).First(&roleType).RecordNotFound() {
		return util.AppResponse{400, "unknown role", nil}
	}
	r.Params.BindJSON(&roleType)
	app.DB.Model(&models.Role{}).Updates(&roleType)
	return util.AppResponse{200, "success", roleType}
}

// Delete role
// TODO: delete only if the role is ADMIN; need to verify
func (r *Role) DeleteRole(role string) revel.Result {
	roleType := models.Role{Name: role}
	if app.DB.Where(&roleType).First(&roleType).RecordNotFound() {
		return util.AppResponse{400, "unknown role", nil}
	}
	app.DB.Model(&models.Role{}).Delete(&roleType)
	return util.AppResponse{200, "success", nil}
}

// Get users belongs to particular role
func (r *Role) GetUsers(role string) revel.Result {
	roleType := models.Role{Name: role}
	if app.DB.Where(&roleType).First(&roleType).RecordNotFound() {
		return util.AppResponse{400, "unknown role", nil}
	}
	users := new([]models.User)
	app.DB.Model(&roleType).Related(&users, "Users")
	return util.AppResponse{200, "success", users}
}

// Get permissions belongs to particular role
func (r *Role) GetPerms(role string) revel.Result {
	roleType := models.Role{}
	if app.DB.Where("name = ?", role).First(&roleType).RecordNotFound() {
		return util.AppResponse{400, "unknown role", nil}
	}
	perms := new([]models.Permission)
	app.DB.Model(&roleType).Related(&perms, "Permissions")
	return util.AppResponse{200, "success", perms}
}

// Set new permissions belongs to particular role
func (r *Role) SetPerms(role string) revel.Result {
	roleType := models.Role{}
	if app.DB.Where("name = ?", role).First(&roleType).RecordNotFound() {
		return util.AppResponse{400, "unknown role", nil}
	}
	perms := new([]models.Permission)
	r.Params.BindJSON(&perms)
	revel.AppLog.Debugf("%+v", perms)
	app.DB.Model(&roleType).Association("Permissions").Clear()
	app.DB.Model(&roleType).Association("Permissions").Append(perms)
	app.DB.Model(&roleType).Related(&perms, "Permissions")
	return util.AppResponse{200, "success", perms}
}

// Add/Append new permission(s) to particular role
func (r *Role) AddPerms(role string) revel.Result {
	roleType := models.Role{}
	if app.DB.Where("name = ?", role).First(&roleType).RecordNotFound() {
		return util.AppResponse{400, "unknown role", nil}
	}
	perms := new([]models.Permission)
	r.Params.BindJSON(&perms)
	app.DB.Model(&roleType).Association("Permissions").Append(perms)
	app.DB.Model(&roleType).Related(&perms, "Permissions")
	return util.AppResponse{200, "success", perms}
}

// Add/Append new user(s) to particular role
func (r *Role) AddUsers(role string) revel.Result {
	roleType := models.Role{}

	if app.DB.Where("name = ?", role).First(&roleType).RecordNotFound() {
		return util.AppResponse{400, "unknown role", nil}
	}

	users := new([]models.User)
	r.Params.BindJSON(&users)
	app.DB.Model(&roleType).Association("Users").Append(users)
	app.DB.Model(&roleType).Related(&users, "Users")
	return util.AppResponse{200, "success", users}
}

// Add new user(s) to particular role
func (r *Role) SetUsers(role string) revel.Result {
	roleType := models.Role{}
	if app.DB.Where("name = ?", role).First(&roleType).RecordNotFound() {
		return util.AppResponse{400, "unknown role", nil}
	}
	users := new([]models.User)
	r.Params.BindJSON(&users)
	app.DB.Model(&roleType).Association("Users").Clear()
	app.DB.Model(&roleType).Association("Users").Append(users)
	app.DB.Model(&roleType).Related(&users, "Users")
	return util.AppResponse{200, "success", users}
}
