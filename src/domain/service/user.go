package service

import (
	"context"

	"github.com/coma/coma/src/application/user/dto"
)

type UserServicer interface {
	CreateUser(context.Context, dto.RequestCreateUser) (dto.ResponseUser, error)
	DeleteUser(context.Context, dto.RequestUser) error
	UpdateUser(context.Context, dto.RequestUser) (dto.ResponseUser, error)
	PatchUserPassword(context.Context, dto.RequestPatchUserPassword) error
	FindUser(context.Context, dto.RequestUser) (dto.ResponseUser, error)
	FindUsers(context.Context, dto.RequestUsers) (dto.ResponseUsers, error)
}
