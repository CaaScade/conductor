package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	global_model "github.com/koki/conductor/app/src/app/models"
	user_model "github.com/koki/conductor/app/src/user/models"
	role_model "github.com/koki/conductor/app/src/roles/models"
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

	DB.AutoMigrate(&global_model.Global{})
	DB.AutoMigrate(&auth_identity.AuthIdentity{})
	DB.AutoMigrate(&user_model.User{})
	DB.AutoMigrate(&role_model.Role{})
	DB.AutoMigrate(&role_model.Permission{})

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

	permission := role_model.Permission{
		Name:     "all",
		Resource: "*",
		Create:   true,
		Read:     true,
		Update:   true,
		Delete:   true,
	}
	DB.LogMode(true)
	role := role_model.Role{
		Name: "admin",
	}

	user := user_model.User{
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
	DB.Model(&permission).Association("Roles").Append([]*role_model.Role{&role})
	DB.Model(&role).Association("Permissions").Append([]*role_model.Permission{&permission})
	DB.Model(&role).Association("Users").Append([]*user_model.User{&user})
	//	DB.Model(&user).Association("Roles").Append([]*role_model.Role{&role})

	AddExitEventHandler(dbShutdownHandler)
}

func dbShutdownHandler() {
	revel.AppLog.Infof("closing database connection")
	DB.Close()
}
