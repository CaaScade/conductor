package v1alpha1

import (
	"errors"

	"conductor/pkg/client"
	"conductor/pkg/resource"

	log "github.com/sirupsen/logrus"
)

const DASHBOARD = "dashboard"

var dashboardClient *Dashboard

func init() {
	dashboardClient = &Dashboard{}

	resource.ResourcesMap[DASHBOARD] = &DashboardCodec{}
}

type DashboardCodec struct{}

func (u *DashboardCodec) Create(r interface{}) (interface{}, error) {
	dashboard, ok := r.(*Dashboard)
	if !ok {
		log.Errorf("could not create dashboard")
		return nil, errors.New("invalid object")
	}
	return dashboardClient.Create(dashboard)
}

func (u *DashboardCodec) List(r *resource.ListOptions) ([]interface{}, error) {
	return nil, nil
}

func (u *DashboardCodec) Update(old interface{}, new interface{}) (interface{}, error) {
	return nil, nil
}

func (u *DashboardCodec) Delete(r interface{}) (interface{}, error) {
	return nil, nil
}

func (u *DashboardCodec) Get(name string, options *resource.GetOptions) (interface{}, error) {
	return nil, nil
}

type Dashboard struct {
	Name string `json:"name"`
}

func (u *Dashboard) Create(user *Dashboard) (*Dashboard, error) {
	dashboardInt, err := client.Create(DASHBOARD, user.Name, user)
	if err != nil {
		return nil, err
	}
	return dashboardInt.(*Dashboard), nil
}

func (u *Dashboard) Get(name string, options *resource.GetOptions) (*Dashboard, error) {

	//dashboardInt, err := client.Get(DASHBOARD, options.Name, user)
	//if err != nil {
	//	return nil, err
	//}
	//return dashboardInt.(*Dashboard), nil
	return nil, nil
}
