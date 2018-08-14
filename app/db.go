package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/koki/conductor/app/models"
	"github.com/qor/auth/auth_identity"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
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
	DB.AutoMigrate(&models.Role{})
	DB.AutoMigrate(&models.Permission{})

	var globalConfig models.Global
	if revel.Config.BoolDefault(AUTHENTICATED_CONF, false) {
		globalConfig.AuthenticationMode = int(models.AuthenticationModePassword)
		revel.Config.SetOption(AUTHENTICATED_CONF, "true")
	} else {
		globalConfig.AuthenticationMode = int(models.AuthenticationModeUnsafe)
	}
	if DB.Where(&globalConfig).First(&globalConfig).RecordNotFound() {
		DB.Create(&globalConfig)
	}

	permission := models.Permission{
		Name:     "all",
		Resource: "*",
		Create:   true,
		Read:     true,
		Update:   true,
		Delete:   true,
	}
	DB.LogMode(true)
	role := models.Role{
		Name: "admin",
	}

	user := models.User{
		Username: "admin",
		Password: "",
		Counter:  1,
	}
	if DB.Where(&permission).First(&permission).RecordNotFound() {
		DB.Create(&permission)
	}

	if DB.Where(&role).First(&role).RecordNotFound() {
		DB.Create(&role)
	}

	if DB.Where(&user).First(&user).RecordNotFound() {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			revel.AppLog.Fatalf("could not create admin user: %+v", err)
		}
		user.Password = string(hashedPassword)
		DB.Create(&user)
	}
	DB.Model(&permission).Association("Roles").Append([]*models.Role{&role})
	DB.Model(&role).Association("Permissions").Append([]*models.Permission{&permission})
	DB.Model(&role).Association("Users").Append([]*models.User{&user})
	//	DB.Model(&user).Association("Roles").Append([]*models.Role{&role})

	AddExitEventHandler(dbShutdownHandler)
}

func dbShutdownHandler() {
	revel.AppLog.Infof("closing database connection")
	DB.Close()
}
