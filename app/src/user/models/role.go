package models

import (
	"github.com/jinzhu/gorm"
)

// Role Represent Roles for the each user with permission belongs to them
type Role struct {
	gorm.Model

	// Name of the role
	Name        string        `gorm:"index"`

	// List of the Users belongs to particular role
	Users       []User `gorm:"many2many:user_roles;"`

	// List of the permission belongs to particular role
	Permissions []Permission  `gorm:"many2many:role_permissions;"`
}


// Permission Represent access for the each resource
type Permission struct {
	gorm.Model

	// Name of the Permission
	Name     string `gorm:index`

	// Name of the resource
	// eg: application
	Resource string

	// Create permission to resource
	// true if you want to provide create permission to resource
	Create   bool

	// Read permission to resource
	// true if you want to provide Read permission to resource
	Read     bool

	// Update permission to resource
	// true if you want to provide Update permission to resource
	Update   bool

	// Delete permission to resource
	// true if you want to provide Delete permission to resource
	Delete   bool

	// List of the Roles belongs to particular permission
	Roles []Role `gorm:"many2many:role_permissions;"`
}
