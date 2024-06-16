package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/nurcahyaari/coma/config"
	"github.com/nurcahyaari/coma/container"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/application/user/dto"
	"github.com/nurcahyaari/coma/src/domain/entity"
	domainrepository "github.com/nurcahyaari/coma/src/domain/repository"
	"github.com/nurcahyaari/coma/src/domain/service"
	"github.com/rs/zerolog/log"
)

type UserService struct {
	config *config.Config
	reader domainrepository.RepositoryUserReader
	writer domainrepository.RepositoryUserWriter
}

func NewUserService(config *config.Config, c container.Container) service.UserServicer {
	return &UserService{
		config: config,
		reader: c.Repository.RepositoryUserReader,
		writer: c.Repository.RepositoryUserWriter,
	}
}

func (s *UserService) InternalFindUser(ctx context.Context, req dto.RequestUser) (entity.User, error) {
	user, err := s.reader.FindUser(ctx, entity.FilterUser{
		Id:       req.Id,
		Username: req.Username,
		UserType: req.UserType,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[FindUser.FindUser] err: failed to find user")
		return entity.User{}, internalerrors.New(err)
	}

	return user, nil
}

func (s *UserService) InternalFindUsers(ctx context.Context, req dto.RequestUsers) (entity.Users, error) {
	users, err := s.reader.FindUsers(ctx, entity.FilterUser{})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[FindUsers.FindUsers] err: failed to find users")
		return entity.Users{}, internalerrors.New(err)
	}

	return users, nil
}

func (s *UserService) CreateRootUser(ctx context.Context, req dto.RequestCreateUser) (dto.ResponseUser, error) {
	var (
		resp dto.ResponseUser
		user = req.UserRoot()
	)

	existingUser, err := s.InternalFindUser(ctx, dto.RequestUser{
		UserType: entity.UserTypeRoot,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[CreateRootUser.FindUser] err: failed to find user")
		return resp, internalerrors.New(err)
	}
	if !existingUser.Empty() {
		log.Warn().
			Msg("[CreateRootUser] user root has already exists")
		return resp, internalerrors.New(
			errors.New("err: user root has already exists"),
			internalerrors.SetErrorCode(http.StatusConflict))
	}

	if err := user.HashPassword(); err != nil {
		log.Error().
			Err(err).
			Msg("[CreateRootUser.HashPassword] err: failed to hash password")
		return resp, internalerrors.New(err)
	}

	if err := s.writer.SaveUser(ctx, user); err != nil {
		log.Error().
			Err(err).
			Msg("[CreateRootUser.SaveUser] err: failed to save user")
		return resp, internalerrors.New(err)
	}

	resp = dto.NewResponseUser(user)

	return resp, nil
}

func (s *UserService) CreateUser(ctx context.Context, req dto.RequestCreateUserNonRoot) (dto.ResponseUser, error) {
	var (
		resp dto.ResponseUser
		user = req.User()
	)

	existingUser, err := s.InternalFindUser(ctx, dto.RequestUser{
		Username: req.Username,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[CreateUser.FindUser] err: failed to find user")
		return resp, internalerrors.New(err)
	}
	if !existingUser.Empty() {
		log.Warn().
			Msg("[CreateUser] user has already exists")
		return resp, internalerrors.New(
			errors.New("err: user has already exists"),
			internalerrors.SetErrorCode(http.StatusConflict))
	}

	if err := user.HashPassword(); err != nil {
		log.Error().
			Err(err).
			Msg("[CreateUser.HashPassword] err: failed to hash password")
		return resp, internalerrors.New(err)
	}

	if err := s.writer.SaveUser(ctx, user); err != nil {
		log.Error().
			Err(err).
			Msg("[CreateUser.SaveUser] err: failed to save user")
		return resp, internalerrors.New(err)
	}

	resp = dto.NewResponseUser(user)

	return resp, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req dto.RequestUser) error {
	user, err := s.reader.FindUser(ctx, entity.FilterUser{
		Id: req.Id,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[DeleteUser.FindUser] err: failed to get user")
		return internalerrors.New(err)
	}
	if user.Empty() {
		log.Warn().
			Msg("[DeleteUser.FindUser] warn: User not found")
		return internalerrors.New(
			errors.New("err: User not found"),
			internalerrors.SetErrorCode(http.StatusNotFound),
		)
	}

	err = s.writer.DeleteUser(ctx, entity.FilterUser{
		Id: req.Id,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[DeleteUser.DeleteUser] err: failed to delete user")
		return internalerrors.New(err)
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
		return resp, internalerrors.New(err)
	}
	if user.Empty() {
		log.Warn().
			Msg("[UpdateUser.FindUser] warn: User not found")
		return resp, internalerrors.New(
			errors.New("err: User not found"),
			internalerrors.SetErrorCode(http.StatusNotFound),
		)
	}

	userOld.Update(user)

	if err := s.writer.UpdateUser(ctx, userOld); err != nil {
		log.Error().
			Err(err).
			Msg("[UpdateUser.UpdateUser] err: failed to update user")
		return resp, internalerrors.New(err)
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
		return internalerrors.New(err)
	}
	if user.Empty() {
		log.Warn().
			Msg("[PatchUserPassword.FindUser] warn: User not found")
		return internalerrors.New(errors.New("err: User not found"), internalerrors.SetErrorCode(http.StatusNotFound))
	}

	if err := user.PatchUserPassword(req.Passowrd); err != nil {
		log.Error().
			Err(err).
			Msg("[PatchUserPassword.PatchUserPassword] err: failed to patch user password")
		return internalerrors.New(err)
	}

	if err := s.writer.UpdateUser(ctx, user); err != nil {
		log.Error().
			Err(err).
			Msg("[PatchUserPassword.UpdateUser] err: failed to update user")
		return internalerrors.New(err)
	}

	return nil
}

func (s *UserService) FindUser(ctx context.Context, req dto.RequestUser) (dto.ResponseUser, error) {
	var (
		resp dto.ResponseUser
	)

	user, err := s.InternalFindUser(ctx, req)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[FindUser.FindUser] err: failed to find user")
		return resp, internalerrors.New(err)
	}
	if user.Empty() {
		log.Warn().
			Msg("[FindUser.FindUser] warn: User not found")
		return resp, internalerrors.New(
			errors.New("err: User not found"),
			internalerrors.SetErrorCode(http.StatusNotFound))
	}

	resp = dto.NewResponseUser(user)

	return resp, nil
}

func (s *UserService) FindUsers(ctx context.Context, req dto.RequestUsers) (dto.ResponseUsers, error) {
	var (
		resp dto.ResponseUsers
	)

	users, err := s.InternalFindUsers(ctx, req)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[FindUsers.FindUsers] err: failed to find users")
	}

	resp = dto.NewResponseUsers(users)

	return resp, nil
}
