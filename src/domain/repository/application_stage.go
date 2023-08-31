package repository

import (
	"context"

	"github.com/coma/coma/src/domain/entity"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryApplicationStageReader
type RepositoryApplicationStageReader interface {
	FindStage(ctx context.Context, filter entity.FilterApplicationStage) (entity.ApplicationStage, bool, error)
	FindStages(ctx context.Context, filter entity.FilterApplicationStage) (entity.ApplicationStages, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryApplicationStageWriter
type RepositoryApplicationStageWriter interface {
	CreateStage(ctx context.Context, data entity.ApplicationStage) error
	DeleteStage(ctx context.Context, filter entity.FilterApplicationStage) error
	UpdateStage(ctx context.Context, data entity.ApplicationStage) error
}
