package service

import (
	"context"

	"github.com/nurcahyaari/coma/config"
	"github.com/nurcahyaari/coma/container"
	"github.com/nurcahyaari/coma/src/application/user/dto"
	"github.com/nurcahyaari/coma/src/domain/entity"
	"github.com/nurcahyaari/coma/src/domain/repository"
	"github.com/nurcahyaari/coma/src/domain/service"
	"github.com/rs/zerolog/log"
)

type UserApplicationScopeService struct {
	config *config.Config
	writer repository.RepositoryUserApplicationScopeWriter
	reader repository.RepositoryUserApplicationScopeReader
}

func NewUserApplicationScopeService(config *config.Config, c container.Container) service.UserApplicationScopeServicer {
	return &UserApplicationScopeService{
		config: config,
		writer: c.Repository.RepositoryUserApplicationScopeWriter,
		reader: c.Repository.RepositoryUserApplicationScopeReader,
	}
}

func (s *UserApplicationScopeService) InternalFindUserApplicationScope(ctx context.Context, req dto.RequestFindUserApplicationScope) (entity.UserApplicationScope, bool, error) {
	userApplicationScope, exists, err := s.reader.FindUserApplicationScope(ctx, entity.FilterUserApplicationScope{
		UserId: req.UserId,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[InternalFindUserApplicationScope.FindUserApplicationScope] err: failed to find user access scope")
		return entity.UserApplicationScope{}, false, err
	}

	return userApplicationScope, exists, nil
}

func (s *UserApplicationScopeService) FindUserApplicationsScope(ctx context.Context, req dto.RequestFindUserApplicationScope) (dto.ResponseUserApplicationsScope, error) {
	userApplicationScope, err := s.reader.FindUserApplicationsScope(ctx, entity.FilterUserApplicationScope{})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[FindUserApplicationScope.InternalFindUserApplicationScope] err: failed to find user access scope")
		return dto.ResponseUserApplicationsScope{}, err
	}

	return dto.NewResponseUserApplicationsScope(userApplicationScope), nil
}

func (s *UserApplicationScopeService) UpsetUserApplicationScope(ctx context.Context, req dto.RequestCreateUserApplicationScope) error {
	var (
		userApplicationScope = req.UserApplicationScope()
	)

	existingUserApplicationScope, exist, err := s.InternalFindUserApplicationScope(ctx, dto.RequestFindUserApplicationScope{
		UserId: req.UserId,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[UpsetUserApplicationScope.InternalFindUserApplicationScope] err: failed to retrieve user access scope")
		return err
	}

	if exist {
		// if user access already exist, it means you neet to update based on request resource
		existingUserApplicationScope.UpdateRbac(userApplicationScope)

		if err := s.writer.UpdateUserApplicationScope(ctx, existingUserApplicationScope); err != nil {
			log.Error().
				Err(err).
				Msg("[UpsetUserApplicationScope.UpdateUserApplicationScope] err: failed to update user access scope")
			return err
		}
		return nil
	}

	if err := s.writer.SaveUserApplicationScope(ctx, userApplicationScope); err != nil {
		log.Error().
			Err(err).
			Msg("[UpsetUserApplicationScope.SaveUserApplicationScope] err: failed to save user access scope")
		return err
	}

	return nil
}
