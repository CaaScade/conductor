package models

import (
	"github.com/jinzhu/gorm"
)

// Application for the users
type Application struct {
	gorm.Model

	// Name of the application
	Name                    string `gorm:"index"`
	PodName                 string
	Description             string
	Price                   float32
	PerMonth                bool
	PerYear                 bool
	UpTime                  string
	URL                     string
	ArchitectureURL         string
	AdditionalReferencesURL string
	Discount                float32

	// List of the Users belongs to particular application
	Users []User `gorm:"many2many:user_apps;"`
}
