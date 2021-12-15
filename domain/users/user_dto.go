package users

import (
	"strings"
	"users-api/utils/errors"
)

//data transfer object

const (
	StatusActive = "active"
)

type User struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	NationalID  string `json:"national_id"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
	Password    string `json:"password"`
}

type Users []User

func (user *User) Validate() *errors.RestErr {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	user.NationalID = strings.TrimSpace(strings.ToLower(user.NationalID))
	if user.NationalID == "" || len(user.NationalID) != 10 {
		return errors.NewBadRequestError("invalid national id")
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return errors.NewBadRequestError("invalid password")
	}
	return nil
}
