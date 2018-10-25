package v1alpha1

import (
	"errors"

	"conductor/pkg/resource"

	log "github.com/sirupsen/logrus"
)

const USER = "user"

var userClient *User

func init() {
	userClient = &User{}

	resource.ResourcesMap[USER] = &UserCodec{}
}

type UserCodec struct{}

func (u *UserCodec) Create(r interface{}) (interface{}, error) {
	user, ok := r.(*User)
	if !ok {
		log.Errorf("could not create user")
		return nil, errors.New("invalid object")
	}
	return userClient.Create(user)
}

func (u *UserCodec) List(r *resource.ListOptions) ([]interface{}, error) {
	return nil, nil
}

func (u *UserCodec) Update(old resource.Resource, new resource.Resource) (resource.Resource, error) {
	return nil, nil
}

func (u *UserCodec) Delete(r resource.Resource) (resource.Resource, error) {
	return nil, nil
}

func (u *UserCodec) Get(name string, options *resource.GetOptions) (resource.Resource, error) {
	return nil, nil
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) Create(user *User) (*User, error) {
	//client.Create(USER, user.Username, user)
	return nil, nil
}

func (u *User) Get(name string, options *resource.GetOptions) (*User, error) {
	return nil, nil
}
