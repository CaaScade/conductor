package models

import (
	"github.com/jinzhu/gorm"
)

// Global structure represent authentication mode
type Global struct {
	gorm.Model

	// AuthenticationMode will be number which represent type of authentication mode
	AuthenticationMode int
}

type AuthenticationMode int

const (
	AuthenticationModePassword AuthenticationMode = iota
	AuthenticationModeUnsafe
)
