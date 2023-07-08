package service

import (
	"context"

	"github.com/coma/coma/internal/utils/routine"
	"github.com/coma/coma/src/domains/application/dto"
	"github.com/coma/coma/src/domains/application/model"
	"github.com/coma/coma/src/domains/application/repository"
	"github.com/rs/zerolog/log"
)

type ApplicationKeyServicer interface {
	FindApplicationKey(ctx context.Context, request dto.RequestFindApplicationKey) (dto.ResponseFindApplicationKey, error)
	GenerateOrUpdateApplicationKey(ctx context.Context, request dto.RequestCreateApplicationKey) (dto.ResponseCreateApplicationKey, error)
}

type ApplicationKeyService struct {
	reader            repository.RepositoryApplicationKeyReader
	writer            repository.RepositoryApplicationKeyWriter
	applicationReader repository.RepositoryApplicationReader
	applicationWriter repository.RepositoryApplicationWriter
	stageReader       repository.RepositoryApplicationStageReader
	stageWriter       repository.RepositoryApplicationStageWriter
}

type ApplicationKeyServiceOptions func(s *ApplicationKeyService)

func SetApplicationKeyRepository(reader repository.RepositoryApplicationKeyReader, writer repository.RepositoryApplicationKeyWriter) ApplicationKeyServiceOptions {
	return func(s *ApplicationKeyService) {
		s.writer = writer
		s.reader = reader
	}
}

func SetApplicationKeyApplicationRepository(reader repository.RepositoryApplicationReader, writer repository.RepositoryApplicationWriter) ApplicationKeyServiceOptions {
	return func(s *ApplicationKeyService) {
		s.applicationWriter = writer
		s.applicationReader = reader
	}
}

func SetApplicationKeyStageRepository(reader repository.RepositoryApplicationStageReader, writer repository.RepositoryApplicationStageWriter) ApplicationKeyServiceOptions {
	return func(s *ApplicationKeyService) {
		s.stageWriter = writer
		s.stageReader = reader
	}
}

func NewApplicationKey(opts ...ApplicationKeyServiceOptions) ApplicationKeyServicer {
	svc := &ApplicationKeyService{}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

func (s *ApplicationKeyService) FindApplicationKey(ctx context.Context, request dto.RequestFindApplicationKey) (dto.ResponseFindApplicationKey, error) {
	var (
		response         dto.ResponseFindApplicationKey
		filter           = request.FilterApplicationKey()
		application      model.Application
		applicationStage model.ApplicationStage
		applicationKey   model.ApplicationKey
	)

	rtn := routine.New()

	rtn.Add("findApplication", &application, func(params ...any) (any, error) {
		applicationId := params[0].(string)
		resp, err := s.applicationReader.FindApplication(ctx, model.FilterApplication{
			Id: applicationId,
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("[GenerateOrUpdateApplicationKey.FindApplications] error find application")
			return nil, err
		}
		return &resp, nil
	}, request.ApplicationId)

	rtn.Add("findStage", &applicationStage, func(params ...any) (any, error) {
		resp, err := s.stageReader.FindStage(ctx, model.FilterApplicationStage{
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
			return nil, nil
		}
		return &resp, nil
	})

	rtn.Start()
	if rtn.IsError() {
		log.Error().
			Errs("routine error", rtn.Errors()).
			Msg("[GenerateOrUpdateApplicationKey] eror on goroutine")
		return response, rtn.Error()
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
		application      model.Application
		applicationStage model.ApplicationStage
	)

	// generate application key
	applicationKey.GenerateSalt(12)
	err := applicationKey.GenerateKey()
	if err != nil {
		log.Error().
			Err(err).
			Msg("[GenerateOrUpdateApplicationKey.GenerateKey] error generating key")
		return response, err
	}

	rtn := routine.New()

	rtn.Add("findApplication", &application, func(params ...any) (any, error) {
		applicationId := params[0].(string)
		resp, err := s.applicationReader.FindApplication(ctx, model.FilterApplication{
			Id: applicationId,
		})
		if err != nil {
			log.Error().
				Err(err).
				Msg("[GenerateOrUpdateApplicationKey.FindApplications] error find application")
			return nil, err
		}
		return &resp, nil
	}, request.ApplicationId)

	rtn.Add("findStage", &applicationStage, func(params ...any) (any, error) {
		resp, err := s.stageReader.FindStage(ctx, model.FilterApplicationStage{
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

	rtn.Start()
	if rtn.IsError() {
		log.Error().
			Errs("routine error", rtn.Errors()).
			Msg("[GenerateOrUpdateApplicationKey] eror on goroutine")
		return response, rtn.Error()
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
