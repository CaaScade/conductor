package models

import (
	"github.com/jinzhu/gorm"
)

// User Represent users in our system
type User struct {
	gorm.Model

	// Username of the particular user
	// this will be use for the login
	Username string `gorm:"index"`

	// Password of the particular user
	// this will be use for the login
	// password will must be encrypted before stored in the DB
	Password string `json:"-"`

	// Email ID of the particular user
	Email string

	// Counter will be number of time user has login to the system
	Counter uint64 `json:"-"`

	// List of the roles belongs to particular user
	Roles []Role `gorm:"many2many:user_roles;"`

	Alerts []Alerts `gorm:"many2many:alert_user;"`

	Dashboards []Dashboard `gorm:"many2many:dashboard_user;"`
}
