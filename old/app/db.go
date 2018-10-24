package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	global_model "github.com/koki/conductor/app/src/app/models"
	"github.com/koki/conductor/app/src/user/models"
	"github.com/qor/auth/auth_identity"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

var DB *gorm.DB

// init db
// also create some initial roles and users with it
func InitDB() {
	localPath := revel.Config.StringDefault("koki.db.location", "/tmp/koki.db")
	db, err := gorm.Open("sqlite3", localPath)
	if err != nil {
		revel.AppLog.Fatalf("could not connect to database: %+v", err)
	}
	DB = db

	DB.AutoMigrate(&global_model.Global{})
	DB.AutoMigrate(&auth_identity.AuthIdentity{})
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Role{})
	DB.AutoMigrate(&models.Permission{})
	DB.AutoMigrate(&models.Alerts{})
	DB.AutoMigrate(&models.Dashboard{})
	DB.AutoMigrate(&models.DashboardDatasource{})
	DB.AutoMigrate(&models.Application{})
	DB.AutoMigrate(&models.ApplicationConfig{})

	var globalConfig global_model.Global
	if revel.Config.BoolDefault(AUTHENTICATED_CONF, false) {
		globalConfig.AuthenticationMode = int(global_model.AuthenticationModePassword)
		revel.Config.SetOption(AUTHENTICATED_CONF, "true")
	} else {
		globalConfig.AuthenticationMode = int(global_model.AuthenticationModeUnsafe)
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
	svcRole := models.Role{
		Name: "service",
	}

	user := models.User{
		Username: "admin",
		Password: "",
		Counter:  1,
	}

	svcUser := models.User{
		Username: "service",
		Password: "",
		Counter:  1,
	}

	alerts := models.Alerts{}
	if DB.Where(&permission).First(&permission).RecordNotFound() {
		DB.Create(&permission)
	}

	if DB.Where(&role).First(&role).RecordNotFound() {
		DB.Create(&role)
	}

	if DB.Where(&svcRole).First(&svcRole).RecordNotFound() {
		DB.Create(&svcRole)
	}

	if DB.Where(&user).First(&user).RecordNotFound() {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		if err != nil {
			revel.AppLog.Fatalf("could not create admin user: %+v", err)
		}
		user.Password = string(hashedPassword)
		DB.Create(&user)
	}
	if DB.Where(&svcUser).First(&svcUser).RecordNotFound() {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("service"), bcrypt.DefaultCost)
		if err != nil {
			revel.AppLog.Fatalf("could not create service user: %+v", err)
		}
		svcUser.Password = string(hashedPassword)
		DB.Create(&svcUser)
	}
	DB.Model(&permission).Association("Roles").Append([]*models.Role{&role})
	DB.Model(&role).Association("Permissions").Append([]*models.Permission{&permission})
	DB.Model(&svcRole).Association("Permissions").Append([]*models.Permission{&permission})
	DB.Model(&role).Association("Users").Append([]*models.User{&user})
	DB.Model(&svcRole).Association("Users").Append([]*models.User{&svcUser})
	DB.Model(&alerts).Association("Users").Append([]*models.User{&user})

	AddExitEventHandler(dbShutdownHandler)
}

// close the db connections
func dbShutdownHandler() {
	revel.AppLog.Infof("closing database connection")
	DB.Close()
}
