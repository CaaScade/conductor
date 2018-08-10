package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Username string `gorm:"type:text`
	Password string `gorm:"type:text`
	Email    string `gorm:"type:text`
	Counter  uint64 `gorm:"type:bigint`
}
