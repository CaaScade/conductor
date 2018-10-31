package v1alpha1

import (
	"conductor/pkg/client"
	"conductor/pkg/resource"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const USER = "user"

var userClient *User

func init() {
	userClient = &User{}

	resource.ResourcesMap[USER] = &UserCodec{}
}

type UserCodec struct{}

func (u *UserCodec) Create(r interface{}) (interface{}, error) {
	var user User
	user.BaseResource = &resource.BaseResource{}

	err := json.Unmarshal(r.([]byte), &user)
	if err != nil {
		log.Errorf("could not unmarshal data")
		return nil, errors.New("internal error")
	}

	return userClient.Create(&user)
}

func (u *UserCodec) List(r *resource.ListOptions) ([]interface{}, error) {

	return nil, nil
}

func (u *UserCodec) Update(new interface{}) (interface{}, error) {
	var user User
	user.BaseResource = &resource.BaseResource{}

	err := json.Unmarshal(new.([]byte), &user)
	if err != nil {
		log.Errorf("could not unmarshal data")
		return nil, errors.New("internal error")
	}

	return userClient.Update(&user)
}

func (u *UserCodec) Delete(name string,  options *resource.DeleteOptions) (interface{}, error) {
	return userClient.Delete(name, options)
}

func (u *UserCodec) Get(name string, options *resource.GetOptions) (interface{}, error) {
	return userClient.Get(name, options)
}
func (u *UserCodec) GetBaseResource(userInt interface{}) *resource.BaseResource {
	return userInt.(*User).BaseResource
}


type User struct {
	*resource.BaseResource

	Username string `json:"username"`
	Password string `json:"password"`
	Email string `json:"email,omitempty"`
	Counter int `json:"counter,omitempty"`
}

type GetOptions struct {
	Username string `json:"username"`
}

func (u *User) Create(user *User) (interface{}, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	userInt, err := client.Create(USER, user.Username, user, user.BaseResource)

	if err != nil {
		log.Errorf("could not create user %+v", err)
		return nil, err
	}
	return userInt, nil
}

func (u *User) Update(user *User) (interface{}, error) {
	userInt, err := client.Update(USER, user.Username, user)
	if err != nil {
		return nil, err
	}
	return userInt.(*User), nil
}

func (u *User) Get(name string, options *resource.GetOptions) (interface{}, error) {

	userInt, err := client.Get(USER, name)
	if err != nil {
		return nil, err
	}
	return userInt, nil
}

func (u *User) Delete(name string, options *resource.DeleteOptions) (interface{}, error) {

	userInt, err := client.Delete(USER, name)
	if err != nil {
		return nil, err
	}
	return userInt, nil
}
