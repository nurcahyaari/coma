package service

import (
	"context"
	"net/http"

	internalerrors "github.com/coma/coma/internal/utils/errors"
	"github.com/coma/coma/src/domains/application/dto"
	"github.com/coma/coma/src/domains/application/model"
	"github.com/coma/coma/src/domains/application/repository"
	"github.com/rs/zerolog/log"
)

type ApplicationServicer interface {
	FindApplications(ctx context.Context, request dto.RequestFindApplication) (dto.ResponseApplications, error)
	CreateApplication(ctx context.Context, request dto.RequestCreateApplication) (dto.ResponseApplication, error)
	DeleteApplication(ctx context.Context, request dto.RequestFindApplication) error
}

type ApplicationService struct {
	reader      repository.RepositoryApplicationReader
	writer      repository.RepositoryApplicationWriter
	stageReader repository.RepositoryApplicationStageReader
	stageWriter repository.RepositoryApplicationStageWriter
}

type ApplicationServiceOptions func(s *ApplicationService)

func SetApplicationRepository(
	applicationRepo *repository.Repository) ApplicationServiceOptions {
	return func(s *ApplicationService) {
		s.writer = applicationRepo.NewRepositoryApplicationWriter()
		s.reader = applicationRepo.NewRepositoryApplicationReader()
		s.stageReader = applicationRepo.NewRepositoryApplicationStageReader()
		s.stageWriter = applicationRepo.NewRepositoryApplicationStageWriter()
	}
}

func NewApplication(opts ...ApplicationServiceOptions) ApplicationServicer {
	svc := &ApplicationService{}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

func (s *ApplicationService) FindApplications(ctx context.Context, request dto.RequestFindApplication) (dto.ResponseApplications, error) {
	var (
		response = dto.ResponseApplications{}
	)

	applications, err := s.reader.FindApplications(ctx, model.FilterApplication{
		Name:    request.Name,
		StageId: request.StageId,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[FindApplications.FindApplications] error find applications")
		return response, err
	}

	response = dto.NewResponseApplications(applications)

	return response, nil
}

func (s *ApplicationService) CreateApplication(ctx context.Context, request dto.RequestCreateApplication) (dto.ResponseApplication, error) {
	var (
		application = request.NewApplication()
		response    = dto.ResponseApplication{}
	)

	if err := request.Validate(); err != nil {
		return response, internalerrors.NewError(
			err,
			internalerrors.SetErrorSource(internalerrors.OZZO_VALIDATION_ERR))
	}

	stage, err := s.stageReader.FindStage(ctx, model.FilterApplicationStage{
		Id: request.StageId,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[CreateEnvirontment.FindStage] error finding stage")
		return response, internalerrors.NewError(err)
	}
	if stage.Id == "" {
		log.Error().
			Err(err).
			Msg("[CreateEnvirontment.FindStage] error stage doesn't found")
		return response, internalerrors.NewError(err,
			internalerrors.SetErrorCode(http.StatusNotFound))
	}

	err = s.writer.CreateApplication(ctx, application)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[CreateEnvirontment] error creating new environment")
		return response, internalerrors.NewError(err)
	}

	response = dto.NewResponseApplication(application)

	return response, nil
}

func (s *ApplicationService) DeleteApplication(ctx context.Context, request dto.RequestFindApplication) error {
	err := s.writer.DeleteApplication(ctx, model.FilterApplication{
		Id:   request.Id,
		Name: request.Name,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[DeleteApplication] error deleting application")
		return internalerrors.NewError(err)
	}
	return nil
}
