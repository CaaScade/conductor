package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/koki/conductor/app"
	"github.com/koki/conductor/app/src/user/models"
	"github.com/koki/conductor/app/src/util"
	"github.com/revel/revel"
	"io/ioutil"
	"net/http"
)

type Dashboard struct {
	*revel.Controller
}

func (d *Dashboard) GetDashboard(username string) revel.Result {
	user := models.User{Username: username}
	if app.DB.Where(&user).Find(&user).RecordNotFound() {
		return util.AppResponse{400, "unknown user name", nil}
	}

	dashboard := new([]models.Dashboard)
	app.DB.Model(&user).Related(&dashboard, "Dashboards")

	return util.AppResponse{200, "success", dashboard}
}

func (d *Dashboard) CreateDashboard(username string) revel.Result {
	var dashboardData map[string]interface{}
	d.Params.BindJSON(&dashboardData)

	user := models.User{Username: username}

	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return util.AppResponse{400, "user not found", nil}
	}

	dashboardurl, err := d.callGrafanaCreateDashboard(dashboardData["Data"])
	if err != nil {
		return util.AppResponse{500, err.Error(), nil}
	}
	dashboard := new(models.Dashboard)
	dashboard.Data = dashboardurl
	d.Params.BindJSON(&dashboard)
	app.DB.Create(&dashboard)
	app.DB.Model(&dashboard).Association("Users").Append(&user)
	app.DB.Model(&dashboard).Related(&user, "Users")

	return util.AppResponse{200, "Success", dashboard}
}

func (d *Dashboard) UpdateDashboard(username string) revel.Result {
	user := models.User{Username: username}
	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return util.AppResponse{400, "user not found", nil}
	}

	dashboard := new(models.Dashboard)
	d.Params.BindJSON(&dashboard)
	app.DB.Model(&models.Dashboard{}).Updates(&dashboard)

	return util.AppResponse{200, "Success", dashboard}
}

func (d *Dashboard) GetDataSource(dashboardid int) revel.Result {
	dashboard := models.Dashboard{}
	if app.DB.First(&dashboard, dashboardid).RecordNotFound() {
		return util.AppResponse{400, "dashboard not found", nil}
	}

	ds := new([]models.DashboardDatasource)
	app.DB.Model(&dashboard).Related(&ds, "Datasource")

	return util.AppResponse{200, "success", ds}
}

func (d *Dashboard) CreateDataSource(dashboardid int) revel.Result {

	var datasourceData map[string]interface{}
	d.Params.BindJSON(&datasourceData)

	dashboard := models.Dashboard{}
	if app.DB.First(&dashboard, dashboardid).RecordNotFound() {
		return util.AppResponse{400, "dashboard not found", nil}
	}

	datasource, err := d.callGrafanaCreateDashboard(datasourceData["Data"])
	if err != nil {
		return util.AppResponse{500, err.Error(), nil}
	}

	ds := new(models.DashboardDatasource)
	ds.Datasource = datasource
	d.Params.BindJSON(&ds)
	app.DB.Create(&ds)
	app.DB.Model(&ds).Association("Datasource").Append(&dashboard)
	app.DB.Model(&ds).Related(&dashboard, "Datasource")

	return util.AppResponse{200, "Success", dashboard}
}

func (d *Dashboard) UpdateDataSource(dashboardid int) revel.Result {
	dashboard := models.Dashboard{}
	if app.DB.First(&dashboard, dashboardid).RecordNotFound() {
		return util.AppResponse{400, "dashboard not found", nil}
	}

	ds := new(models.DashboardDatasource)
	d.Params.BindJSON(&ds)
	app.DB.Model(&models.DashboardDatasource{}).Updates(&ds)

	return util.AppResponse{200, "Success", ds}
}

func (d *Dashboard) callGrafanaCreateDashboard(data interface{}) (string, error) {
	if data == nil {
		return "", nil
	}

	jsondata, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest("POST", "http://localhost:3000/api/dashboards/db", bytes.NewBuffer(jsondata))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer eyJrIjoiSzBhNjZWcm5KV3RQQmNsT3Y3VlNsanQ4TjdTNm5GQjciLCJuIjoiYXBpIiwiaWQiOjF9")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return "", err
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("response %s\n", string(data))

		var dat map[string]interface{}
		if err := json.Unmarshal(data, &dat); err != nil {
			revel.AppLog.Errorf("The HTTP request failed with error %s\n", err)
			return "", err
		}

		return dat["url"].(string), nil
	}
}

func (d *Dashboard) callGrafanaCreateDatasource(data interface{}) (string, error) {
	if data == nil {
		return "", nil
	}

	jsondata, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest("POST", "http://localhost:3000/api/datasources", bytes.NewBuffer(jsondata))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer eyJrIjoiSzBhNjZWcm5KV3RQQmNsT3Y3VlNsanQ4TjdTNm5GQjciLCJuIjoiYXBpIiwiaWQiOjF9")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return "", err
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("response %s\n", string(data))

		var dat map[string]interface{}
		if err := json.Unmarshal(data, &dat); err != nil {
			revel.AppLog.Errorf("The HTTP request failed with error %s\n", err)
			return "", err
		}

		return dat["url"].(string), nil
	}
}
