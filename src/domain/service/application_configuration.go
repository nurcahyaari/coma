package service

import (
	"context"

	"github.com/coma/coma/src/application/application/dto"
)

type ApplicationConfigurationServicer interface {
	GetConfigurationViewTypeJSON(ctx context.Context, req dto.RequestGetConfiguration) (dto.ResponseGetConfigurationViewTypeJSON, error)
	GetConfigurationViewTypeSchema(ctx context.Context, req dto.RequestGetConfiguration) (dto.ResponseGetConfigurationsViewTypeSchema, error)
	SetConfiguration(ctx context.Context, req dto.RequestSetConfiguration) (dto.ResponseSetConfiguration, error)
	UpdateConfiguration(ctx context.Context, req dto.RequestUpdateConfiguration) error
	UpsertConfiguration(ctx context.Context, req dto.RequestSetConfiguration) error
	DeleteConfiguration(ctx context.Context, req dto.RequestDeleteConfiguration) error
	DistributeConfiguration(ctx context.Context, clientKey string) error
}
