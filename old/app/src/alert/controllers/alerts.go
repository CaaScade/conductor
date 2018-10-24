package controllers

import (
	"github.com/koki/conductor/app"
	"github.com/koki/conductor/app/src/user/models"
	"github.com/koki/conductor/app/src/util"
	"github.com/revel/revel"
)

type Alerts struct {
	*revel.Controller
}

// get details of particular slack channel alert
func (sl *Alerts) GetAlertData(username string) revel.Result {
	user := models.User{Username: username}
	if app.DB.Where(&user).Find(&user).RecordNotFound() {
		return util.AppResponse{400, "unknown user name", nil}
	}

	slackAlert := new([]models.Alerts)
	app.DB.Model(&user).Related(&slackAlert, "Alerts")

	return util.AppResponse{200, "success", slackAlert}
}

// add/update slack alert
func (sl *Alerts) AddSlackAlert(username string) revel.Result {

	user := models.User{Username: username}
	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return util.AppResponse{400, "user not found", nil}
	}

	slackAlert := new(models.Alerts)
	sl.Params.BindJSON(&slackAlert)
	app.DB.Create(&slackAlert)
	app.DB.Model(&slackAlert).Association("Users").Append(&user)
	app.DB.Model(&slackAlert).Related(&user, "Users")

	return util.AppResponse{200, "Success", slackAlert}
}

// add/update slack alert
func (sl *Alerts) UpdateSlackAlert(username string) revel.Result {

	user := models.User{Username: username}
	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return util.AppResponse{400, "user not found", nil}
	}

	slackAlert := new(models.Alerts)
	sl.Params.BindJSON(&slackAlert)
	app.DB.Model(&models.Alerts{}).Updates(&slackAlert)

	return util.AppResponse{200, "Success", slackAlert}
}
