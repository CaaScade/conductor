package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/koki/conductor/app/models"
	"github.com/qor/auth/auth_identity"
	"github.com/revel/revel"
)

var DB *gorm.DB

func InitDB() {
	localPath := revel.Config.StringDefault("koki.db.location", "/tmp/koki.db")
	db, err := gorm.Open("sqlite3", localPath)
	if err != nil {
		revel.AppLog.Fatalf("could not connect to database: %+v", err)
	}
	DB = db

	DB.AutoMigrate(&models.Global{})
	DB.AutoMigrate(&auth_identity.AuthIdentity{})
	DB.AutoMigrate(&models.User{})

	var globalConfig models.Global
	if revel.Config.BoolDefault(AUTHENTICATED_CONF, false) {
		globalConfig.AuthenticationMode = int(models.AuthenticationModePassword)
		revel.Config.SetOption(AUTHENTICATED_CONF, "true")
	} else {
		globalConfig.AuthenticationMode = int(models.AuthenticationModeUnsafe)
	}
	DB.Create(&globalConfig)

	AddExitEventHandler(dbShutdownHandler)
}

func dbShutdownHandler() {
	revel.AppLog.Infof("closing database connection")
	DB.Close()
}
