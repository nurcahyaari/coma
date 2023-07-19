package repository

import (
	"context"

	"github.com/coma/coma/src/domain/entity"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryApplicationConfigurationWriter
type RepositoryApplicationConfigurationWriter interface {
	SetConfiguration(ctx context.Context, data entity.Configuration) (string, error)
	DeleteConfiguration(ctx context.Context, filter entity.FilterConfiguration) error
	UpdateConfiguration(ctx context.Context, data entity.Configuration) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryApplicationConfigurationReader
type RepositoryApplicationConfigurationReader interface {
	FindClientConfiguration(ctx context.Context, filter entity.FilterConfiguration) (entity.Configurations, error)
}
