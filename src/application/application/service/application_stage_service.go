package service

import (
	"context"

	"github.com/coma/coma/config"
	internalerrors "github.com/coma/coma/internal/utils/errors"
	"github.com/coma/coma/src/application/application/dto"
	"github.com/coma/coma/src/application/application/repository"
	"github.com/coma/coma/src/domains/entity"
	domainrepository "github.com/coma/coma/src/domains/repository"
	"github.com/coma/coma/src/domains/service"
	"github.com/rs/zerolog/log"
)

type ApplicationStageService struct {
	config *config.Config
	reader domainrepository.RepositoryApplicationStageReader
	writer domainrepository.RepositoryApplicationStageWriter
}

type ApplicationStageServiceOptions func(s *ApplicationStageService)

func SetApplicationStageRepository(applicationRepo *repository.Repository) ApplicationStageServiceOptions {
	return func(s *ApplicationStageService) {
		s.writer = applicationRepo.NewRepositoryApplicationStageWriter()
		s.reader = applicationRepo.NewRepositoryApplicationStageReader()
	}
}

func NewApplicationStage(config *config.Config, opts ...ApplicationStageServiceOptions) service.ApplicationStageServicer {
	svc := &ApplicationStageService{
		config: config,
	}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

func (s *ApplicationStageService) FindStages(ctx context.Context, request dto.RequestFindStage) (dto.ResponseStages, error) {
	var (
		response = dto.ResponseStages{}
	)
	applicationStages, err := s.reader.FindStages(ctx, entity.FilterApplicationStage{
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
	err := s.writer.DeleteStage(ctx, entity.FilterApplicationStage{
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
