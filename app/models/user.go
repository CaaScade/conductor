package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Username string `gorm:"type:text`
	Password string `json:"-" gorm:"type:text`
	Email    string `gorm:"type:text`
	Counter  uint64 `json:"-" gorm:"type:bigint`

	Roles []Role `gorm:"many_to_many:user_roles;"`
}

/*
func (user *User) Scan(data interface{}) (err error) {
	switch values := data.(type) {
	case []byte:
		if string(values) != "" {
			return json.Unmarshal(values, user)
		}
	case string:
		return user.Scan([]byte(values))
	case []string:
		for _, str := range values {
			if err := user.Scan(str); err != nil {
				return err
			}
		}
	default:
		err = errors.New("unsupported driver -> Scan for user")
	}
	return
}

// Value return struct's Value
func (user *User) Value() (driver.Value, error) {
	results, err := json.Marshal(user)
	return string(results), err
}
*/
