package services

import (
	"users-api/domain/users"
	"users-api/utils/crypto_utils"
	"users-api/utils/date_utils"
	"users-api/utils/errors"
)

var (
	UsersService usersServiceInterface = &usersService{}
)

type usersService struct {}

type usersServiceInterface interface {
	GetUser(int64) (*users.User, *errors.RestErr)
	GetAllUser(int,string) (users.Users, *errors.RestErr,int64)
	CreateUser(users.User) (*users.User, *errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, *errors.RestErr)
	DeleteUser(int64) *errors.RestErr
	SearchUser(int,string) (users.Users, *errors.RestErr,int64)
}

func (s *usersService) GetUser(userId int64) (*users.User, *errors.RestErr) {
	result := &users.User{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *usersService) GetAllUser(page int,status string) (users.Users, *errors.RestErr,int64){
	dao := &users.User{}
	return dao.GetAll(page,status)
}

func (s *usersService) CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDBFormat()
	user.Password = crypto_utils.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *usersService) UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr) {
	current := &users.User{Id: user.Id}
	if err := current.Get(); err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if user.NationalID != "" {
			current.NationalID = user.NationalID
		}
	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.NationalID = user.NationalID
	}

	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

func (s *usersService) DeleteUser(userId int64) *errors.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func (s *usersService) SearchUser(page int ,status string) (users.Users, *errors.RestErr,int64) {
	dao := &users.User{}
	return dao.FindByStatus(page,status)
}

