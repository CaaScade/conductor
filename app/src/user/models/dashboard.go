package models

import (
	"github.com/jinzhu/gorm"
)

// User Represent users in our system
type Dashboard struct {
	gorm.Model

	// url is webhook of the slack
	// used for the sending the message
	Data string `gorm:"index"`

	Datasource []DashboardDatasource `gorm:"many2many:dashboard_datasource;"`

	// List of the roles belongs to particular user
	Users []User `gorm:"many2many:dashboard_user;"`
}

// User Represent users in our system
type DashboardDatasource struct {
	gorm.Model

	// url is webhook of the slack
	// used for the sending the message

	Datasource string

	// List of the roles belongs to particular user
	Dashboard []Dashboard `gorm:"many2many:dashboard_datasource;"`
}
