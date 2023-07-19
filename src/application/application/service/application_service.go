package service

import (
	"context"
	"net/http"

	"github.com/coma/coma/config"
	"github.com/coma/coma/container"
	internalerrors "github.com/coma/coma/internal/utils/errors"
	"github.com/coma/coma/src/application/application/dto"
	"github.com/coma/coma/src/domains/entity"
	domainrepository "github.com/coma/coma/src/domains/repository"
	"github.com/coma/coma/src/domains/service"
	"github.com/rs/zerolog/log"
)

type ApplicationService struct {
	config      *config.Config
	reader      domainrepository.RepositoryApplicationReader
	writer      domainrepository.RepositoryApplicationWriter
	stageReader domainrepository.RepositoryApplicationStageReader
	stageWriter domainrepository.RepositoryApplicationStageWriter
}

func NewApplication(config *config.Config, c container.Container) service.ApplicationServicer {
	svc := &ApplicationService{
		config:      config,
		reader:      c.Repository.RepositoryApplicationReader,
		writer:      c.Repository.RepositoryApplicationWriter,
		stageReader: c.Repository.RepositoryApplicationStageReader,
		stageWriter: c.Repository.RepositoryApplicationStageWriter,
	}
	return svc
}

func (s *ApplicationService) FindApplications(ctx context.Context, request dto.RequestFindApplication) (dto.ResponseApplications, error) {
	var (
		response = dto.ResponseApplications{}
	)

	applications, err := s.reader.FindApplications(ctx, entity.FilterApplication{
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

	stage, err := s.stageReader.FindStage(ctx, entity.FilterApplicationStage{
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
	err := s.writer.DeleteApplication(ctx, entity.FilterApplication{
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
