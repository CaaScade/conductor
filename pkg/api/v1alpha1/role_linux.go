package v1alpha1

import (
	"errors"

	"conductor/pkg/client"
	"conductor/pkg/resource"

	log "github.com/sirupsen/logrus"
)

const ROLE = "role"

var roleClient *Role

func init() {
	roleClient = &Role{}

	resource.ResourcesMap[ROLE] = &RoleCodec{}
}

type RoleCodec struct{}

func (u *RoleCodec) Create(r interface{}) (interface{}, error) {
	role, ok := r.(*Role)
	if !ok {
		log.Errorf("could not create role")
		return nil, errors.New("invalid object")
	}
	return roleClient.Create(role)
}

func (u *RoleCodec) List(r *resource.ListOptions) ([]interface{}, error) {
	return nil, nil
}

func (u *RoleCodec) Update(old interface{}, new interface{}) (interface{}, error) {
	return nil, nil
}

func (u *RoleCodec) Delete(r interface{}) (interface{}, error) {
	return nil, nil
}

func (u *RoleCodec) Get(name string, options *resource.GetOptions) (interface{}, error) {
	return nil, nil
}

type Role struct {
	Name string `json:"name"`
}

func (u *Role) Create(role *Role) (*Role, error) {
	roleInt, err := client.Create(ROLE, role.Name, role)
	if err != nil {
		return nil, err
	}
	return roleInt.(*Role), nil
}

func (u *Role) Get(name string, options *resource.GetOptions) (*Role, error) {

	//roleInt, err := client.Get(ROLE, options.Name, &Role{})
	//if err != nil {
	//	return nil, err
	//}
	//return roleInt.(*Role), nil
	return nil, nil
}
