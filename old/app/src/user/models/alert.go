package models

import (
	"github.com/jinzhu/gorm"
)

// User Represent users in our system
type Alerts struct {
	gorm.Model

	// url is webhook of the slack
	// used for the sending the message
	Url string `gorm:"index"`

	// you can mention particular recipient for the alert
	Recipient string

	// mention person or the tag for the alert
	Mention string

	// token is secret used for sending message
	Token string

	// List of the roles belongs to particular user
	Users []User `gorm:"many2many:alert_user;"`
}
