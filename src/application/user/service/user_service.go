package service

import (
	"context"

	"github.com/coma/coma/config"
	"github.com/coma/coma/container"
	"github.com/coma/coma/src/application/user/dto"
	"github.com/coma/coma/src/domain/entity"
	domainrepository "github.com/coma/coma/src/domain/repository"
	"github.com/coma/coma/src/domain/service"
	"github.com/rs/zerolog/log"
)

type UserService struct {
	config *config.Config
	reader domainrepository.RepositoryUserReader
	writer domainrepository.RepositoryUserWriter
}

func NewUserRepository(config *config.Config, c container.Container) service.UserServicer {
	return &UserService{
		config: config,
		reader: c.Repository.RepositoryUserReader,
		writer: c.Repository.RepositoryUserWriter,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.RequestCreateUser) (dto.ResponseUser, error) {
	var (
		resp dto.ResponseUser
		user = req.User()
	)

	if err := user.HashPassword(); err != nil {
		log.Error().
			Err(err).
			Msg("[CreateUser.HashPassword] err: failed to hash password")
		return resp, err
	}

	if err := s.writer.SaveUser(ctx, user); err != nil {
		log.Error().
			Err(err).
			Msg("[CreateUser.SaveUser] err: failed to save user")
		return resp, err
	}

	resp = dto.NewResponseUser(user)

	return resp, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req dto.RequestUser) error {
	_, err := s.reader.FindUser(ctx, entity.FilterUser{
		Id: req.Id,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[DeleteUser.FindUser] err: failed to get user")
		return err
	}

	err = s.writer.DeleteUser(ctx, entity.FilterUser{
		Id: req.Id,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[DeleteUser.DeleteUser] err: failed to delete user")
		return err
	}

	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, req dto.RequestUser) (dto.ResponseUser, error) {
	var (
		resp dto.ResponseUser
		user = req.User()
	)

	userOld, err := s.reader.FindUser(ctx, entity.FilterUser{
		Id: req.Id,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[UpdateUser.FindUser] err: failed to get user")
		return resp, err
	}

	userOld.Update(user)

	if err := s.writer.UpdateUser(ctx, userOld); err != nil {
		log.Error().
			Err(err).
			Msg("[UpdateUser.UpdateUser] err: failed to update user")
		return resp, err
	}

	resp = dto.NewResponseUser(userOld)

	return resp, nil
}

func (s *UserService) PatchUserPassword(ctx context.Context, req dto.RequestPatchUserPassword) error {
	user, err := s.reader.FindUser(ctx, entity.FilterUser{
		Id: req.Id,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[PatchUserPassword.FindUser] err: failed to get user")
		return err
	}

	if err := user.PatchUserPassword(req.Passowrd); err != nil {
		log.Error().
			Err(err).
			Msg("[PatchUserPassword.PatchUserPassword] err: failed to patch user password")
		return err
	}

	if err := s.writer.UpdateUser(ctx, user); err != nil {
		log.Error().
			Err(err).
			Msg("[PatchUserPassword.UpdateUser] err: failed to update user")
		return err
	}

	return nil
}

func (s *UserService) FindUser(ctx context.Context, req dto.RequestUser) (dto.ResponseUser, error) {
	var (
		resp dto.ResponseUser
	)

	user, err := s.reader.FindUser(ctx, entity.FilterUser{
		Id: req.Id,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[FindUser.FindUser] err: failed to find user")
	}

	resp = dto.NewResponseUser(user)

	return resp, nil
}

func (s *UserService) FindUsers(ctx context.Context, req dto.RequestUsers) (dto.ResponseUsers, error) {
	var (
		resp dto.ResponseUsers
	)

	users, err := s.reader.FindUsers(ctx, entity.FilterUser{})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[FindUsers.FindUsers] err: failed to find users")
	}

	resp = dto.NewResponseUsers(users)

	return resp, nil
}
