package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/nurcahyaari/coma/config"
	"github.com/nurcahyaari/coma/container"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/internal/x/routine"
	"github.com/nurcahyaari/coma/src/application/application/dto"
	"github.com/nurcahyaari/coma/src/domain/entity"
	domainrepository "github.com/nurcahyaari/coma/src/domain/repository"
	"github.com/nurcahyaari/coma/src/domain/service"
	"github.com/rs/zerolog/log"
)

type ApplicationKeyService struct {
	config            *config.Config
	reader            domainrepository.RepositoryApplicationKeyReader
	writer            domainrepository.RepositoryApplicationKeyWriter
	applicationReader domainrepository.RepositoryApplicationReader
	applicationWriter domainrepository.RepositoryApplicationWriter
}

func NewApplicationKey(config *config.Config, c container.Container) service.ApplicationKeyServicer {
	svc := &ApplicationKeyService{
		config:            config,
		reader:            c.Repository.RepositoryApplicationKeyReader,
		writer:            c.Repository.RepositoryApplicationKeyWriter,
		applicationReader: c.Repository.RepositoryApplicationReader,
		applicationWriter: c.Repository.RepositoryApplicationWriter,
	}
	return svc
}

func (s *ApplicationKeyService) IsExistsApplicationKey(ctx context.Context, request dto.RequestFindApplicationKey) (bool, error) {
	var (
		response       bool
		filter         = request.FilterApplicationKey()
		applicationKey entity.ApplicationKey
		err            error
	)

	// skip validation
	filter.SkipValidation = true

	applicationKey, err = s.reader.FindApplicationKey(ctx, filter)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[FindApplicationKey.FindApplicationKey] error find application key")
		return response, internalerrors.New(err)
	}

	if applicationKey.Id == "" {
		err = errors.New("err: application key doesn't exists")
		log.Error().
			Err(errors.New("application key not found")).
			Msg("[FindApplicationKey.FindApplicationKey] error application key not found")
		return response, internalerrors.New(err, internalerrors.SetErrorCode(http.StatusNotFound))
	}

	response = true

	return response, nil
}

func (s *ApplicationKeyService) FindApplicationKey(ctx context.Context, request dto.RequestFindApplicationKey) (dto.ResponseFindApplicationKey, error) {
	var (
		response       dto.ResponseFindApplicationKey
		filter         = request.FilterApplicationKey()
		application    entity.Application
		applicationKey entity.ApplicationKey
	)

	if err := request.Validate(); err != nil {
		return response, internalerrors.New(
			err,
			internalerrors.SetErrorSource(internalerrors.OZZO_VALIDATION_ERR))
	}

	rtn := routine.New()

	rtn.Add("findApplication", &application, func(params ...any) (any, error) {
		applicationId := params[0].(string)
		resp, exist, err := s.applicationReader.FindApplication(ctx, entity.FilterApplication{
			Id: applicationId,
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("[GenerateOrUpdateApplicationKey.FindApplications] error find application")
			return nil, internalerrors.New(err)
		}
		if !exist {
			err = errors.New("err: application not found")
			log.Error().
				Err(err).
				Msg("[GenerateOrUpdateApplicationKey.FindApplications] error: application not found")
			return nil, internalerrors.New(err, internalerrors.SetErrorCode(http.StatusNotFound))
		}

		return &resp, nil
	}, request.ApplicationId)

	rtn.Add("findKey", &applicationKey, func(params ...any) (any, error) {
		resp, err := s.reader.FindApplicationKey(ctx, filter)
		if err != nil {
			log.Error().
				Err(err).
				Msg("[FindApplicationKey.FindApplicationKey] error find application key")
			return nil, internalerrors.New(err)
		}
		return &resp, nil
	})

	rtn.Start()
	if rtn.IsError() {
		log.Error().
			Errs("routine error", rtn.Errors()).
			Msg("[GenerateOrUpdateApplicationKey] eror on goroutine")
		return response, internalerrors.New(rtn.Error())
	}

	response = dto.NewResponseFindApplicationKey(applicationKey)
	response.AttachApplication(application)

	return response, nil
}

func (s *ApplicationKeyService) GenerateOrUpdateApplicationKey(ctx context.Context, request dto.RequestCreateApplicationKey) (dto.ResponseCreateApplicationKey, error) {
	var (
		response       dto.ResponseCreateApplicationKey
		applicationKey = request.ApplicationKey()
		application    entity.Application
	)

	if err := request.Validate(); err != nil {
		return response, internalerrors.New(
			err,
			internalerrors.SetErrorSource(internalerrors.OZZO_VALIDATION_ERR))
	}

	// generate application key
	applicationKey.GenerateKey(36)

	rtn := routine.New()

	rtn.Add("findApplication", &application, func(params ...any) (any, error) {
		applicationId := params[0].(string)
		resp, exist, err := s.applicationReader.FindApplication(ctx, entity.FilterApplication{
			Id: applicationId,
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("[GenerateOrUpdateApplicationKey.FindApplications] error find application")
			return nil, internalerrors.New(err)
		}
		if !exist {
			err = errors.New("err: application not found")
			log.Error().
				Err(err).
				Msg("[GenerateOrUpdateApplicationKey.FindApplications] error: application not found")
			return nil, internalerrors.New(err, internalerrors.SetErrorCode(http.StatusNotFound))
		}

		return &resp, nil
	}, request.ApplicationId)

	rtn.Start()
	if rtn.IsError() {
		log.Error().
			Errs("routine error", rtn.Errors()).
			Msg("[GenerateOrUpdateApplicationKey] eror on goroutine")
		return response, internalerrors.New(rtn.Error())
	}

	err := s.writer.CreateOrSaveApplicationKey(ctx, applicationKey)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[GenerateOrUpdateApplicationKey] error create or save application key")
		return response, err
	}

	response = dto.ResponseCreateApplicationKey{
		ApplicationName: application.Name,
		Key:             applicationKey.Key,
	}

	return response, nil
}
