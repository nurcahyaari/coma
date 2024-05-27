package service

import (
	"context"

	"github.com/nurcahyaari/coma/src/application/application/dto"
)

type ApplicationServicer interface {
	FindApplications(ctx context.Context, request dto.RequestFindApplication) (dto.ResponseApplications, error)
	CreateApplication(ctx context.Context, request dto.RequestCreateApplication) (dto.ResponseApplication, error)
	DeleteApplication(ctx context.Context, request dto.RequestFindApplication) error
}
