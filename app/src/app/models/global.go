package global_model

import (
	"github.com/jinzhu/gorm"
)

type Global struct {
	gorm.Model

	AuthenticationMode int
}

type AuthenticationMode int

const (
	AuthenticationModePassword AuthenticationMode = iota
	AuthenticationModeUnsafe
)
