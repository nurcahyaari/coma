package service

import (
	"context"

	"github.com/coma/coma/src/application/user/dto"
	"github.com/coma/coma/src/domain/entity"
)

type InternalUserServicer interface {
	InternalFindUser(context.Context, dto.RequestUser) (entity.User, error)
	InternalFindUsers(context.Context, dto.RequestUsers) (entity.Users, error)
}

type UserServicer interface {
	InternalUserServicer
	CreateUser(context.Context, dto.RequestCreateUser) (dto.ResponseUser, error)
	CreateRootUser(ctx context.Context, req dto.RequestCreateUser) (dto.ResponseUser, error)
	DeleteUser(context.Context, dto.RequestUser) error
	UpdateUser(context.Context, dto.RequestUser) (dto.ResponseUser, error)
	PatchUserPassword(context.Context, dto.RequestPatchUserPassword) error
	FindUser(context.Context, dto.RequestUser) (dto.ResponseUser, error)
	FindUsers(context.Context, dto.RequestUsers) (dto.ResponseUsers, error)
}
