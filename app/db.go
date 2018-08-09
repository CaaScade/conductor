package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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

	revel.AddInitEventHandler(dbShutdownHandler)
}

func dbShutdownHandler(EVENT int, _ interface{}) int {
	if EVENT == revel.ENGINE_SHUTDOWN {
		revel.AppLog.Infof("closing database connection")
		DB.Close()
	}
	return 0
}
