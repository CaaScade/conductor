package v1alpha1

import (
	"errors"

	"conductor/pkg/client"
	"conductor/pkg/resource"

	log "github.com/sirupsen/logrus"
)

const PERMISSION = "permission"

var permissionClient *Permission

func init() {
	permissionClient = &Permission{}

	resource.ResourcesMap[PERMISSION] = &PermissionCodec{}
}

type PermissionCodec struct{}

func (u *PermissionCodec) Create(r interface{}) (interface{}, error) {
	permission, ok := r.(*Permission)
	if !ok {
		log.Errorf("could not create permission")
		return nil, errors.New("invalid object")
	}
	return permissionClient.Create(permission)
}

func (u *PermissionCodec) List(r *resource.ListOptions) ([]interface{}, error) {
	return nil, nil
}

func (u *PermissionCodec) Update(old interface{}, new interface{}) (interface{}, error) {
	return nil, nil
}

func (u *PermissionCodec) Delete(r interface{}) (interface{}, error) {
	return nil, nil
}

func (u *PermissionCodec) Get(name string, options *resource.GetOptions) (interface{}, error) {
	return nil, nil
}

type Permission struct {
	Name string `json:"name"`
	Resource string `json:"resource"`
	CreateP string `json:"create"`
	Read string `json:"read"`
	Update string `json:"update"`
	Delete string `json:"delete"`
}

func (u *Permission) Create(permission *Permission) (*Permission, error) {
	permissionInt, err := client.Create(PERMISSION, permission.Name, permission)
	if err != nil {
		return nil, err
	}
	return permissionInt.(*Permission), nil
}

func (u *Permission) Get(name string, options *resource.GetOptions) (*Permission, error) {

	//permissionInt, err := client.Get(PERMISSION, options.Name, &Permission{})
	//if err != nil {
	//	return nil, err
	//}
	//return permissionInt.(*Permission), nil
	return nil, nil
}
