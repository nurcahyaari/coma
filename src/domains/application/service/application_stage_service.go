package service

import (
	"context"

	internalerrors "github.com/coma/coma/internal/utils/errors"
	"github.com/coma/coma/src/domains/application/dto"
	"github.com/coma/coma/src/domains/application/model"
	"github.com/coma/coma/src/domains/application/repository"
	"github.com/rs/zerolog/log"
)

type ApplicationStageServicer interface {
	FindStages(ctx context.Context, request dto.RequestFindStage) (dto.ResponseStages, error)
	CreateStage(ctx context.Context, request dto.RequestCreateStage) (dto.ResponseStage, error)
	DeleteStage(ctx context.Context, request dto.RequestFindStage) error
}

type ApplicationStageService struct {
	reader repository.RepositoryApplicationStageReader
	writer repository.RepositoryApplicationStageWriter
}

type ApplicationStageServiceOptions func(s *ApplicationStageService)

func SetApplicationStageRepository(reader repository.RepositoryApplicationStageReader, writer repository.RepositoryApplicationStageWriter) ApplicationStageServiceOptions {
	return func(s *ApplicationStageService) {
		s.writer = writer
		s.reader = reader
	}
}

func NewApplicationStage(opts ...ApplicationStageServiceOptions) ApplicationStageServicer {
	svc := &ApplicationStageService{}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

func (s *ApplicationStageService) FindStages(ctx context.Context, request dto.RequestFindStage) (dto.ResponseStages, error) {
	var (
		response = dto.ResponseStages{}
	)
	applicationStages, err := s.reader.FindStages(ctx, model.FilterApplicationStage{
		Name: request.Name,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[FindStages] error find stages")
		return response, internalerrors.NewError(err)
	}

	response = dto.NewResponseStages(applicationStages)

	return response, nil
}

func (s *ApplicationStageService) CreateStage(ctx context.Context, request dto.RequestCreateStage) (dto.ResponseStage, error) {
	var (
		applicationEnv = request.NewApplicationStage()
		response       = dto.ResponseStage{}
	)
	err := s.writer.CreateOrSaveStage(ctx, applicationEnv)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[CreateEnvirontment] error creating new environment")
		return response, internalerrors.NewError(err)
	}

	response = dto.NewResponseStage(applicationEnv)

	return response, nil
}

func (s *ApplicationStageService) DeleteStage(ctx context.Context, request dto.RequestFindStage) error {
	err := s.writer.DeleteStage(ctx, model.FilterApplicationStage{
		Name: request.Name,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[DeleteStage] error deleting stage")
		return internalerrors.NewError(err)
	}
	return nil
}
