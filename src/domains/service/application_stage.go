package service

import (
	"context"

	"github.com/coma/coma/src/application/application/dto"
)

type ApplicationStageServicer interface {
	FindStages(ctx context.Context, request dto.RequestFindStage) (dto.ResponseStages, error)
	CreateStage(ctx context.Context, request dto.RequestCreateStage) (dto.ResponseStage, error)
	DeleteStage(ctx context.Context, request dto.RequestFindStage) error
}
