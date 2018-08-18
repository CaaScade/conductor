package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Username string `gorm:"index"`
	Password string `json:"-"`
	Email    string
	Counter  uint64 `json:"-"`

	Roles []Role `gorm:"many2many:user_roles;"`
}
