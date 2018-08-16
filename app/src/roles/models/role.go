package role_model

import (
	"github.com/jinzhu/gorm"
	user_model "github.com/koki/conductor/app/src/user/models"
)

type Role struct {
	gorm.Model

	Name        string        `gorm:"index"`
	Users       []user_model.User `gorm:"many2many:user_roles;"`
	Permissions []Permission  `gorm:"many2many:role_permissions;"`
}

type Permission struct {
	gorm.Model

	Name     string `gorm:index`
	Resource string
	Create   bool
	Read     bool
	Update   bool
	Delete   bool

	Roles []Role `gorm:"many2many:role_permissions;"`
}
