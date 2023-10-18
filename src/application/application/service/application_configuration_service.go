package service

import (
	"context"
	"errors"

	"github.com/coma/coma/config"
	"github.com/coma/coma/container"
	"github.com/coma/coma/infrastructure/integration/coma"
	"github.com/coma/coma/internal/x/pubsub"
	"github.com/coma/coma/src/application/application/dto"
	"github.com/coma/coma/src/domain/entity"
	domainrepository "github.com/coma/coma/src/domain/repository"
	"github.com/coma/coma/src/domain/service"
	"github.com/rs/zerolog/log"
)

type ApplicationConfigurationService struct {
	config            *config.Config
	pubSub            *pubsub.Pubsub
	comaClient        *coma.WebsocketClient
	applicationKeySvc service.ApplicationKeyServicer
	readerRepo        domainrepository.RepositoryApplicationConfigurationReader
	writerRepo        domainrepository.RepositoryApplicationConfigurationWriter
}

func NewApplicationConfiguration(
	cfg *config.Config, c container.Container) service.ApplicationConfigurationServicer {
	svc := &ApplicationConfigurationService{
		config:            cfg,
		pubSub:            c.LocalPubsub,
		comaClient:        c.Integration.WebsocketClient,
		readerRepo:        c.Repository.RepositoryApplicationConfigurationReader,
		writerRepo:        c.Repository.RepositoryApplicationConfigurationWriter,
		applicationKeySvc: c.Service.ApplicationKeyServicer,
	}
	return svc
}

func (s *ApplicationConfigurationService) GetConfigurationViewTypeJSON(ctx context.Context, req dto.RequestGetConfiguration) (dto.ResponseGetConfigurationViewTypeJSON, error) {
	var (
		response dto.ResponseGetConfigurationViewTypeJSON
		err      error
	)

	configurations, err := s.readerRepo.FindClientConfiguration(ctx, entity.FilterConfiguration{
		ClientKey: req.XClientKey,
	})
	if err != nil {
		log.Error().Err(err).Msg("[GetConfiguration] error FindClientConfiguration")
		return response, err
	}

	response = dto.NewResponseGetConfigurationViewTypeJSON(req.XClientKey)
	err = response.SetData(configurations)
	if err != nil {
		log.Error().Err(err).Msg("[GetConfiguration] error NewResponseGetConfiguration")
		return response, err
	}

	return response, nil
}

func (s *ApplicationConfigurationService) GetConfigurationViewTypeSchema(ctx context.Context, req dto.RequestGetConfiguration) (dto.ResponseGetConfigurationsViewTypeSchema, error) {
	var (
		response dto.ResponseGetConfigurationsViewTypeSchema
		err      error
	)

	configurations, err := s.readerRepo.FindClientConfiguration(ctx, entity.FilterConfiguration{
		ClientKey: req.XClientKey,
	})
	if err != nil {
		log.Error().Err(err).Msg("[GetConfiguration] error FindClientConfiguration")
		return response, err
	}

	response = dto.NewResponseGetConfigurationsViewTypeSchema(configurations)

	return response, nil
}

func (s *ApplicationConfigurationService) SetConfiguration(ctx context.Context, req dto.RequestSetConfiguration) (dto.ResponseSetConfiguration, error) {
	if err := req.Validate(); err != nil {
		log.Error().Err(err).Msg("[SetConfiguration] error validate dto")
		return dto.ResponseSetConfiguration{}, err
	}

	var (
		configuration       = req.Configuration()
		filterConfiguration = entity.FilterConfiguration{
			ClientKey: req.XClientKey,
			Field:     req.Field,
		}
	)

	clientConfigurations, err := s.readerRepo.FindClientConfiguration(ctx, filterConfiguration)
	if err != nil {
		log.Error().
			Err(err).
			Str("field", req.Field).
			Msg("[SetConfiguration] error on search configuration")
		return dto.ResponseSetConfiguration{}, err
	}
	if clientConfigurations.Exists() {
		log.Error().
			Err(err).
			Str("field", req.Field).
			Msg("[SetConfiguration] error duplicate field name")
		return dto.ResponseSetConfiguration{}, errors.New("err: duplicate field name")
	}

	insertedId, err := s.writerRepo.SetConfiguration(ctx, configuration)
	if err != nil {
		log.Error().Err(err).Msg("[SetConfiguration] error SetConfiguration")
		return dto.ResponseSetConfiguration{}, err
	}

	// after success writing to the db distribute to the client
	s.pubSub.Publish(s.config.Pubsub.Local.Publisher.ConfigDistributor.Topic,
		pubsub.SendString(req.XClientKey))

	return dto.ResponseSetConfiguration{
		Id: insertedId,
	}, nil
}

func (s *ApplicationConfigurationService) UpdateConfiguration(ctx context.Context, req dto.RequestUpdateConfiguration) error {
	clientConfigurations, err := s.readerRepo.FindClientConfiguration(ctx, entity.FilterConfiguration{
		ClientKey: req.XClientKey,
		Id:        req.Id,
		Field:     req.Field,
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("field", req.Field).
			Msg("[UpdateConfiguration] error on search configuration")
		return err
	}
	if !clientConfigurations.Exists() {
		log.Error().
			Err(err).
			Str("field", req.Field).
			Msg("[UpdateConfiguration] error configuration is empty")
		return errors.New("err: configuration is empty")
	}

	var (
		configuration        = req.Configuration()
		configurations       = entity.Configurations{configuration}
		mapConfigurationById = configurations.MapConfigurationById()
	)

	clientConfigurations.Update(mapConfigurationById)

	for _, configuration := range clientConfigurations {
		err = s.writerRepo.UpdateConfiguration(ctx, configuration)
		if err != nil {
			log.Error().
				Err(err).
				Str("field", req.Field).
				Msg("[UpdateConfiguration] error on update configuration")
			return err
		}
	}

	// after success writing to the db distribute to the client
	s.pubSub.Publish(s.config.Pubsub.Local.Publisher.ConfigDistributor.Topic,
		pubsub.SendString(req.XClientKey))

	return nil
}

func (s *ApplicationConfigurationService) UpsertConfiguration(ctx context.Context, req dto.RequestSetConfiguration) error {
	clientConfigurations, err := s.readerRepo.FindClientConfiguration(ctx, entity.FilterConfiguration{
		ClientKey: req.XClientKey,
		Field:     req.Field,
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("field", req.Field).
			Msg("[UpdateConfiguration] error on search configuration")
		return err
	}

	switch clientConfigurations.Exists() {
	case true:
		// when true it means the client configuration already exists
		// so we need to update it
		err = s.UpdateConfiguration(ctx, dto.RequestUpdateConfiguration{
			XClientKey: clientConfigurations[0].ClientKey,
			Id:         clientConfigurations[0].Id,
			Field:      req.Field,
			Value:      req.Value,
		})
		if err != nil {
			log.Error().
				Err(err).
				Str("field", req.Field).
				Msg("[UpsertConfiguration] error on update configuration")
			return err
		}

	default:
		_, err = s.SetConfiguration(ctx, req)
		if err != nil {
			log.Error().
				Err(err).
				Str("field", req.Field).
				Msg("[UpsertConfiguration] error on insert configuration")
			return err
		}

	}

	return nil
}

func (s *ApplicationConfigurationService) DeleteConfiguration(ctx context.Context, req dto.RequestDeleteConfiguration) error {
	err := s.writerRepo.DeleteConfiguration(ctx, req.FilterConfiguration())
	if err != nil {
		log.Error().Err(err).Msg("[DeleteConfiguration] error when deleting configuration")
		return err
	}

	// after success writing to the db distribute to the client
	s.pubSub.Publish(s.config.Pubsub.Local.Publisher.ConfigDistributor.Topic,
		pubsub.SendString(req.XClientKey))

	return nil
}

func (s *ApplicationConfigurationService) DistributeConfiguration(ctx context.Context, clientKey string) error {
	clientConfiguration, err := s.GetConfigurationViewTypeJSON(ctx, dto.RequestGetConfiguration{
		XClientKey: clientKey,
	})
	if err != nil {
		log.Error().Err(err).
			Msg("[distributeConfiguration.GetConfigurationViewTypeJSON] error when get the configuration")
		return err
	}

	if clientConfiguration.Data == nil {
		err = errors.New("err: data is empty")
		log.Error().Err(err).
			Msg("[distributeConfiguration.GetConfigurationViewTypeJSON] data is empty")
		return err
	}

	err = s.comaClient.Send(coma.RequestSendMessage{
		ClientKey: clientKey,
		Data:      clientConfiguration.Data,
	})
	if err != nil {
		log.Error().Err(err).Msg("[distributeConfiguration.Send] error when distributing configuration to the client")
		return err
	}
	return nil
}
