package repository

import (
	"context"

	"github.com/coma/coma/src/domains/entity"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryApplicationKeyWriter
type RepositoryApplicationKeyWriter interface {
	CreateOrSaveApplicationKey(ctx context.Context, data entity.ApplicationKey) error
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryApplicationKeyReader
type RepositoryApplicationKeyReader interface {
	FindApplicationKey(ctx context.Context, filter entity.FilterApplicationKey) (entity.ApplicationKey, error)
}
