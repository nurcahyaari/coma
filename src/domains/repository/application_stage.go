package repository

import (
	"context"

	"github.com/coma/coma/src/domains/entity"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryApplicationStageReader
type RepositoryApplicationStageReader interface {
	FindStage(ctx context.Context, filter entity.FilterApplicationStage) (entity.ApplicationStage, error)
	FindStages(ctx context.Context, filter entity.FilterApplicationStage) (entity.ApplicationStages, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryApplicationStageWriter
type RepositoryApplicationStageWriter interface {
	CreateOrSaveStage(ctx context.Context, data entity.ApplicationStage) error
	DeleteStage(ctx context.Context, filter entity.FilterApplicationStage) error
}
