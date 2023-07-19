package service

import (
	"context"
	"errors"

	"github.com/coma/coma/config"
	internalerrors "github.com/coma/coma/internal/utils/errors"
	"github.com/coma/coma/internal/utils/routine"
	"github.com/coma/coma/src/application/application/dto"
	"github.com/coma/coma/src/application/application/repository"
	"github.com/coma/coma/src/domains/entity"
	domainrepository "github.com/coma/coma/src/domains/repository"
	"github.com/coma/coma/src/domains/service"
	"github.com/rs/zerolog/log"
)

type ApplicationKeyService struct {
	config            *config.Config
	reader            domainrepository.RepositoryApplicationKeyReader
	writer            domainrepository.RepositoryApplicationKeyWriter
	applicationReader domainrepository.RepositoryApplicationReader
	applicationWriter domainrepository.RepositoryApplicationWriter
	stageReader       domainrepository.RepositoryApplicationStageReader
	stageWriter       domainrepository.RepositoryApplicationStageWriter
}

type ApplicationKeyServiceOptions func(s *ApplicationKeyService)

func SetApplicationKeyRepository(applicationRepo *repository.Repository) ApplicationKeyServiceOptions {
	return func(s *ApplicationKeyService) {
		s.writer = applicationRepo.NewRepositoryApplicationKeyWriter()
		s.reader = applicationRepo.NewRepositoryApplicationKeyReader()
		s.applicationWriter = applicationRepo.NewRepositoryApplicationWriter()
		s.applicationReader = applicationRepo.NewRepositoryApplicationReader()
		s.stageWriter = applicationRepo.NewRepositoryApplicationStageWriter()
		s.stageReader = applicationRepo.NewRepositoryApplicationStageReader()
	}
}

func NewApplicationKey(
	config *config.Config,
	opts ...ApplicationKeyServiceOptions) service.ApplicationKeyServicer {
	svc := &ApplicationKeyService{
		config: config,
	}

	for _, opt := range opts {
		opt(svc)
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
		return response, internalerrors.NewError(err)
	}

	if applicationKey.Id == "" {
		err = errors.New("err: application key doesn't exists")
		log.Error().
			Err(errors.New("application key doesn't found")).
			Msg("[FindApplicationKey.FindApplicationKey] error application key doesn't found")
		return response, internalerrors.NewError(err)
	}

	response = true

	return response, nil
}

func (s *ApplicationKeyService) FindApplicationKey(ctx context.Context, request dto.RequestFindApplicationKey) (dto.ResponseFindApplicationKey, error) {
	var (
		response         dto.ResponseFindApplicationKey
		filter           = request.FilterApplicationKey()
		application      entity.Application
		applicationStage entity.ApplicationStage
		applicationKey   entity.ApplicationKey
	)

	if err := request.Validate(); err != nil {
		return response, internalerrors.NewError(
			err,
			internalerrors.SetErrorSource(internalerrors.OZZO_VALIDATION_ERR))
	}

	rtn := routine.New()

	rtn.Add("findApplication", &application, func(params ...any) (any, error) {
		applicationId := params[0].(string)
		resp, err := s.applicationReader.FindApplication(ctx, entity.FilterApplication{
			Id: applicationId,
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("[GenerateOrUpdateApplicationKey.FindApplications] error find application")
			return nil, internalerrors.NewError(err)
		}
		return &resp, nil
	}, request.ApplicationId)

	rtn.Add("findStage", &applicationStage, func(params ...any) (any, error) {
		resp, err := s.stageReader.FindStage(ctx, entity.FilterApplicationStage{
			Id: request.StageId,
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("[GenerateOrUpdateApplicationKey.FindStages] error find stage")
			return nil, err
		}
		return &resp, nil
	}, request.StageId)

	rtn.Add("findKey", &applicationKey, func(params ...any) (any, error) {
		resp, err := s.reader.FindApplicationKey(ctx, filter)
		if err != nil {
			log.Error().
				Err(err).
				Msg("[FindApplicationKey.FindApplicationKey] error find application key")
			return nil, internalerrors.NewError(err)
		}
		return &resp, nil
	})

	rtn.Start()
	if rtn.IsError() {
		log.Error().
			Errs("routine error", rtn.Errors()).
			Msg("[GenerateOrUpdateApplicationKey] eror on goroutine")
		return response, internalerrors.NewError(rtn.Error())
	}

	response = dto.NewResponseFindApplicationKey(applicationKey)
	response.AttachApplication(application).
		AttachApplicationStage(applicationStage)

	return response, nil
}

func (s *ApplicationKeyService) GenerateOrUpdateApplicationKey(ctx context.Context, request dto.RequestCreateApplicationKey) (dto.ResponseCreateApplicationKey, error) {
	var (
		response         dto.ResponseCreateApplicationKey
		applicationKey   = request.ApplicationKey()
		application      entity.Application
		applicationStage entity.ApplicationStage
	)

	if err := request.Validate(); err != nil {
		return response, internalerrors.NewError(
			err,
			internalerrors.SetErrorSource(internalerrors.OZZO_VALIDATION_ERR))
	}

	// generate application key
	applicationKey.GenerateSalt(12)
	err := applicationKey.GenerateKey()
	if err != nil {
		log.Error().
			Err(err).
			Msg("[GenerateOrUpdateApplicationKey.GenerateKey] error generating key")
		return response, internalerrors.NewError(err)
	}

	rtn := routine.New()

	rtn.Add("findApplication", &application, func(params ...any) (any, error) {
		applicationId := params[0].(string)
		resp, err := s.applicationReader.FindApplication(ctx, entity.FilterApplication{
			Id: applicationId,
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("[GenerateOrUpdateApplicationKey.FindApplications] error find application")
			return nil, internalerrors.NewError(err)
		}
		return &resp, nil
	}, request.ApplicationId)

	rtn.Add("findStage", &applicationStage, func(params ...any) (any, error) {
		resp, err := s.stageReader.FindStage(ctx, entity.FilterApplicationStage{
			Id: request.StageId,
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("[GenerateOrUpdateApplicationKey.FindStages] error find stage")
			return nil, internalerrors.NewError(err)
		}
		return &resp, nil
	}, request.StageId)

	rtn.Start()
	if rtn.IsError() {
		log.Error().
			Errs("routine error", rtn.Errors()).
			Msg("[GenerateOrUpdateApplicationKey] eror on goroutine")
		return response, internalerrors.NewError(rtn.Error())
	}

	err = s.writer.CreateOrSaveApplicationKey(ctx, applicationKey)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[GenerateOrUpdateApplicationKey] error create or save application key")
		return response, err
	}

	response = dto.ResponseCreateApplicationKey{
		ApplicationName: application.Name,
		StageName:       applicationStage.Name,
		Key:             applicationKey.Key,
	}

	return response, nil
}