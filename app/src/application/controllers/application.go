//https://blog.kowalczyk.info/article/JyRZ/generating-good-unique-ids-in-go.html
package controllers

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/koki/conductor/app"
	"github.com/koki/conductor/app/src/user/models"
	"github.com/koki/conductor/app/src/util"
	"github.com/revel/revel"
	xid "github.com/rs/xid"
	appsv1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type Application struct {
	*revel.Controller
}

func (d *Application) GetApplicationList(username string) revel.Result {
	// check for the role
	applications := []models.Application{}
	app.DB.Find(&applications)
	return util.AppResponse{200, "success", applications}
}

func (d *Application) GetApplication(username string) revel.Result {
	user := models.User{Username: username}
	if app.DB.Where(&user).Find(&user).RecordNotFound() {
		return util.AppResponse{400, "unknown user name", nil}
	}

	Applications := new([]models.Application)
	app.DB.Model(&user).Related(&Applications, "Applications")

	return util.AppResponse{200, "success", Applications}
}

func (d *Application) CreateApplication() revel.Result {
	applicationType := models.Application{}
	d.Params.BindJSON(&applicationType)
	if !app.DB.Where(&applicationType).Find(&applicationType).RecordNotFound() {
		return util.AppResponse{400, "application already exists", nil}
	}
	app.DB.Create(&applicationType)
	return util.AppResponse{200, "success", nil}
}

func (d *Application) LaunchApplication(username string) revel.Result {
	var err error
	user := models.User{Username: username}

	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return util.AppResponse{400, "user not found", nil}
	}

	application := new(models.Application)
	d.Params.BindJSON(&application)


	if application.PodName == "redis" {
		err = d.LaunchRedisStack(username)
	} else if application.PodName == "minio" {
		err = d.DeployApplicationStack("neel", application.PodName, application.ConfigData)
	} else if application.PodName == "zk" {
		err = d.DeployApplicationStack("neel", application.PodName, application.ConfigData)
	} else {
		return util.AppResponse{400, "Sorry application installation not found", nil}
	}

	if err != nil {
		glog.Errorf("Error in launching the pod: %v", err.Error())
		return util.AppResponse{500, "Application Launch Error", err.Error()}
	}

	app.DB.Model(&application).Association("Users").Append(&user)
	app.DB.Model(&application).Related(&user, "Users")

	app.DB.Model(&application).Association("ConfigData").Append(&application)
	app.DB.Model(&application).Related(&application, "ConfigData")

	return util.AppResponse{200, "Success", application}
}

func (d *Application) UpdateApplication(username string) revel.Result {
	user := models.User{Username: username}
	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return util.AppResponse{400, "user not found", nil}
	}

	Application := new(models.Application)
	d.Params.BindJSON(&Application)
	app.DB.Model(&models.Application{}).Updates(&Application)

	return util.AppResponse{200, "Success", Application}
}

func (d *Application) DeleteApplication(username string) revel.Result {
	user := models.User{Username: username}
	if app.DB.Where(&user).First(&user).RecordNotFound() {
		return util.AppResponse{400, "user not found", nil}
	}

	Application := new(models.Application)
	d.Params.BindJSON(&Application)
	app.DB.Model(&models.Application{}).Delete(&Application)

	return util.AppResponse{200, "Success", Application}
}

func (d *Application) LaunchRedisStack(username string) error {
	pod := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis-stack",
			Namespace: "neel",
		},
		Spec: corev1.PodSpec{
			NodeName: "ip-172-20-59-132.ec2.internal",
			Containers: []corev1.Container{
				{
					Name:  "redis-stack",
					Image: "neelshah23/redis-operator",
				},
			},
		},
	}
	_, err := app.Client.CoreV1().Pods("neel").Create(&pod)
	if err != nil {
		return err
	}
	return nil

}

func (d *Application) DeployApplicationStack(username string, applicationName string, config []models.ApplicationConfig ) error {
	replicaCount := int32(3)
	uniqueID := xid.New()

	config_Args := []string{}

	for _,element := range config {
		if element.Type == "text" ||  element.Type == "password" {
			config_Args = append(config_Args, fmt.Sprintf("--%s %s", element.Name, element.Value))
		}
	}
	revel.AppLog.Debugf("the config args are %+v", strings.Join(config_Args," "))

	var args []string
	if applicationName == "zk" {
		args = []string{"-e", "http://18.217.71.51:2379", "--id", fmt.Sprintf("/test/%s/%s/", uniqueID.String(), applicationName, "-adminusername","[data from config]")}
	} else if applicationName == "minio" {
		args = []string{"-s", "http://18.217.71.51:2379", "--id", fmt.Sprintf("/test/%s/%s/", uniqueID.String(), applicationName)}
	} else {
		return nil
	}

	pod := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1beta1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s",applicationName, uniqueID.String()),
			Namespace: username,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicaCount,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"koki.io/selector.name": fmt.Sprintf("%s-%s", applicationName, uniqueID.String()),
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: "RollingUpdate",
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"koki.io/selector.name": fmt.Sprintf("%s-%s", applicationName, uniqueID.String()),
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Args: args,
							Env: []corev1.EnvVar{
								{
									Name: "IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.podIP",
										},
									},
								}, {
									Name: "NS",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
							Image:           fmt.Sprintf("wlan0/%s", applicationName),
							ImagePullPolicy: "Always",
							Name:            fmt.Sprintf("%s", applicationName),
							Stdin:           true,
							TTY:             true,
						},
					},
				},
			},
		},
	}
	pod_launch, err := app.Client.AppsV1beta1().Deployments(username).Create(&pod)

	revel.AppLog.Debugf("Pod Launched %+v", pod_launch	)


	if err != nil {
		return err
	}
	return nil

}

func (d *Application) LaunchConductor() revel.Result{
	pod := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "conductor",
			Namespace: "neel",
		},
		Spec: corev1.PodSpec{
			//NodeName: "ip-172-20-59-132.ec2.internal",
			Containers: []corev1.Container{
				{
					Name:  "conductor",
					Image: "wlan0/conductor",
				},
			},
		},
	}
	pod_data, err := app.Client.CoreV1().Pods("neel").Create(&pod)
	if err != nil {
		return util.AppResponse{200, "failed", err}

	}
	return util.AppResponse{200, "Success", pod_data}

}
