package repository

import (
	"context"

	"github.com/coma/coma/src/domain/entity"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryApplicationWriter
type RepositoryApplicationWriter interface {
	CreateApplication(ctx context.Context, data entity.Application) error
	DeleteApplication(ctx context.Context, filter entity.FilterApplication) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryApplicationReader
type RepositoryApplicationReader interface {
	FindApplication(ctx context.Context, filter entity.FilterApplication) (entity.Application, error)
	FindApplications(ctx context.Context, filter entity.FilterApplication) (entity.Applications, error)
}
