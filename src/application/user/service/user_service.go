package service

import "github.com/coma/coma/src/domain/service"

type UserService struct{}

func NewUserRepository() service.UserServicer {
	return &UserService{}
}
