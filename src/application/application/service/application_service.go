package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/nurcahyaari/coma/config"
	"github.com/nurcahyaari/coma/container"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/application/application/dto"
	"github.com/nurcahyaari/coma/src/domain/entity"
	domainrepository "github.com/nurcahyaari/coma/src/domain/repository"
	"github.com/nurcahyaari/coma/src/domain/service"
	domainservice "github.com/nurcahyaari/coma/src/domain/service"
	"github.com/rs/zerolog/log"
)

type ApplicationService struct {
	config            *config.Config
	reader            domainrepository.RepositoryApplicationReader
	writer            domainrepository.RepositoryApplicationWriter
	applicationKeySvc domainservice.ApplicationKeyServicer
}

func NewApplication(config *config.Config, c container.Container) service.ApplicationServicer {
	svc := &ApplicationService{
		config:            config,
		reader:            c.Repository.RepositoryApplicationReader,
		writer:            c.Repository.RepositoryApplicationWriter,
		applicationKeySvc: c.ApplicationKeyServicer,
	}
	return svc
}

func (s *ApplicationService) FindApplications(ctx context.Context, request dto.RequestFindApplication) (dto.ResponseApplications, error) {
	var (
		response = dto.ResponseApplications{}
	)

	applications, err := s.reader.FindApplications(ctx, entity.FilterApplication{
		Name: request.Name,
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
		return response, internalerrors.New(
			err,
			internalerrors.SetErrorSource(internalerrors.OZZO_VALIDATION_ERR))
	}

	_, exist, err := s.reader.FindApplication(ctx, entity.FilterApplication{
		Name: request.Name,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[CreateApplication.FindApplication] error finding")
		return response, internalerrors.New(err)
	}
	if exist {
		err = errors.New("err: application already exists")
		log.Error().
			Err(err).
			Msg("[CreateApplication.FindApplication] error finding")
		return response, internalerrors.New(err,
			internalerrors.SetErrorCode(http.StatusConflict))
	}

	err = s.writer.CreateApplication(ctx, application)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[CreateApplication] error creating new environment")
		return response, internalerrors.New(err)
	}

	applicationKey, err := s.applicationKeySvc.GenerateOrUpdateApplicationKey(ctx, dto.RequestCreateApplicationKey{
		ApplicationId: application.Id,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[CreateApplication.GenerateOrUpdateApplicationKey] error generating key")
		return response, internalerrors.New(err)
	}

	response = dto.NewResponseApplication(application)
	response.AttachApplicationKey(dto.ResponseFindApplicationKey{
		Key: applicationKey.Key,
	})

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
		return internalerrors.New(err)
	}
	return nil
}
