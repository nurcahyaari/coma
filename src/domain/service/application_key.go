package service

import (
	"context"

	"github.com/coma/coma/src/application/application/dto"
)

type ApplicationKeyServicer interface {
	IsExistsApplicationKey(ctx context.Context, request dto.RequestFindApplicationKey) (bool, error)
	FindApplicationKey(ctx context.Context, request dto.RequestFindApplicationKey) (dto.ResponseFindApplicationKey, error)
	GenerateOrUpdateApplicationKey(ctx context.Context, request dto.RequestCreateApplicationKey) (dto.ResponseCreateApplicationKey, error)
}
