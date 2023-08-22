package service

import (
	"context"

	"github.com/coma/coma/config"
	"github.com/coma/coma/container"
	"github.com/coma/coma/src/application/user/dto"
	"github.com/coma/coma/src/domain/entity"
	"github.com/coma/coma/src/domain/repository"
	"github.com/coma/coma/src/domain/service"
	"github.com/rs/zerolog/log"
)

type UserAccessScopeService struct {
	config *config.Config
	writer repository.RepositoryUserAccessScopeWriter
	reader repository.RepositoryUserAccessScopeReader
}

func NewUserAccessScopeService(config *config.Config, c container.Container) service.UserAccessScopeServicer {
	return &UserAccessScopeService{
		config: config,
		writer: c.Repository.RepositoryUserAccessScopeWriter,
		reader: c.Repository.RepositoryUserAccessScopeReader,
	}
}

func (s *UserAccessScopeService) InternalFindUserAccessScope(ctx context.Context, req dto.RequestFindUserAccessScope) (entity.UserAccessScope, error) {
	userAccessScope, err := s.reader.FindUserAccessScope(ctx, entity.FilterUserAccessScope{
		UserId: req.UserId,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[InternalFindUserAccessScope.FindUserAccessScope] err: failed to find user access scope")
		return entity.UserAccessScope{}, err
	}

	return userAccessScope, nil
}

func (s *UserAccessScopeService) FindUserAccessScope(ctx context.Context, req dto.RequestFindUserAccessScope) (dto.ResponseUserAccessScope, error) {
	userAccessScope, err := s.InternalFindUserAccessScope(ctx, req)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[FindUserAccessScope.InternalFindUserAccessScope] err: failed to find user access scope")
		return dto.ResponseUserAccessScope{}, err
	}

	return dto.NewResponseUserAccessScope(userAccessScope), nil
}

func (s *UserAccessScopeService) CreateUserAccessScope(ctx context.Context, req dto.RequestCreateUserAccessScope) error {
	var (
		userAccessScope = req.UserAccessScope()
	)

	if err := s.writer.SaveUserAccessScope(ctx, userAccessScope); err != nil {
		log.Error().
			Err(err).
			Msg("[CreateUserAccessScope.SaveUserAccessScope] err: failed to save user access scope")
		return err
	}

	return nil
}
