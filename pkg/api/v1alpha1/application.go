package v1alpha1

import (
	"encoding/json"
	"errors"

	"conductor/pkg/client"
	"conductor/pkg/resource"

	log "github.com/sirupsen/logrus"
)

const APPLICATION = "application"

var applicationClient *Application

func init() {
	applicationClient = &Application{}

	resource.ResourcesMap[APPLICATION] = &ApplicationCodec{}
}

type ApplicationCodec struct{}

func (u *ApplicationCodec) Create(r interface{}) (interface{}, error) {
	var application Application
	application.BaseResource = &resource.BaseResource{}

	err := json.Unmarshal(r.([]byte), &application)
	if err != nil {
		log.Errorf("could not unmarshal data")
		return nil, errors.New("internal error")
	}

	return applicationClient.Create(&application)
}

func (u *ApplicationCodec) List(r *resource.ListOptions) ([]interface{}, error) {
	return nil, nil
}

func (u *ApplicationCodec) Update(new interface{}) (interface{}, error) {
	return nil, nil
}

func (u *ApplicationCodec) Delete(name string,  options *resource.DeleteOptions) (interface{}, error) {
	return applicationClient.Delete(name, options)
}

func (u *ApplicationCodec) Get(name string, options *resource.GetOptions) (interface{}, error) {
	return applicationClient.Get(name, options)
}
func (u *ApplicationCodec) GetBaseResource(applicationInt interface{}) *resource.BaseResource {
	return applicationInt.(*Application).BaseResource
}

type Application struct {
	*resource.BaseResource

	Name string `json:"name"`
	PodName string `json:"podName"`
	Description string `json:"description"`
	Price float32 `json:"price"`
	PerMonth bool `json:"perMonth"`
	PerYear bool `json:"perYear"`
	URL string `json:"url"`
	ArchitectureURL string `json:"architectureUrl"`
	AdditionalReferencesURL string `json:"addinalReferenceUrl"`
	Discount float32 `json:"discount"`
	IsConfig bool `json:"isConfig"`
	IsReadOnly bool `json:"isReadOnly"`
}

func (u *Application) Create(application *Application) (*Application, error) {
	userInt, err := client.Create(APPLICATION, application.Name, application, application.BaseResource)
	if err != nil {
		return nil, err
	}
	return userInt.(*Application), nil
}

func (u *Application) Get(name string, options *resource.GetOptions) (interface{}, error) {

	userInt, err := client.Get(APPLICATION, name)
	if err != nil {
		return nil, err
	}
	return userInt, nil
}

func (u *Application) Delete(name string, options *resource.DeleteOptions) (interface{}, error) {

	userInt, err := client.Delete(APPLICATION, name)
	if err != nil {
		return nil, err
	}
	return userInt, nil
}

