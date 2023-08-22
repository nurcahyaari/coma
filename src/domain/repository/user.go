package repository

import (
	"context"

	"github.com/coma/coma/src/domain/entity"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryUserReader
type RepositoryUserReader interface {
	FindUser(ctx context.Context, filter entity.FilterUser) (entity.User, error)
	FindUsers(ctx context.Context, filter entity.FilterUser) (entity.Users, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryUserWriter
type RepositoryUserWriter interface {
	SaveUser(ctx context.Context, user entity.User) error
	DeleteUser(ctx context.Context, filter entity.FilterUser) error
	UpdateUser(ctx context.Context, user entity.User) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryUserAccessScopeReader
type RepositoryUserAccessScopeReader interface {
	FindUserAccessScope(ctx context.Context, filter entity.FilterUserAccessScope) (entity.UserAccessScope, error)
	FindUserAccessesScope(ctx context.Context, filter entity.FilterUserAccessScope) (entity.UserAccessesScope, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryUserAccessScopeWriter
type RepositoryUserAccessScopeWriter interface {
	SaveUserAccessScope(ctx context.Context, userAccessScope entity.UserAccessScope) error
	UpdateUserAccessScope(ctx context.Context, userAccessScope entity.UserAccessScope) error
	DeleteUserAccess(ctx context.Context, filter entity.FilterUserAccessScope) error
}
