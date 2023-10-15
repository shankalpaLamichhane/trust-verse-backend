package service

import (
	"time"
	"trust-verse-backend/app/module/user/repository"
)

type UserService struct {
}

type IUserService interface {
	GetUser() (*repository.User, error)
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) GetUser() (repository.User, error) {
	user := repository.User{
		ID:        1,
		Username:  "shankalpa",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	return user, nil
}
