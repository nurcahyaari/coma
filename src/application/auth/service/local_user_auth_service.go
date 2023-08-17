package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/coma/coma/config"
	"github.com/coma/coma/container"
	internalerrors "github.com/coma/coma/internal/utils/errors"
	"github.com/coma/coma/src/application/auth/dto"
	userdto "github.com/coma/coma/src/application/user/dto"
	"github.com/coma/coma/src/domain/entity"
	"github.com/coma/coma/src/domain/repository"
	"github.com/coma/coma/src/domain/service"
	"github.com/rs/zerolog/log"
)

type UserAuthService struct {
	config  *config.Config
	userSvc service.InternalUserServicer
	reader  repository.RepositoryUserAuthReader
	writer  repository.RepositoryUserAuthWriter
}

func NewUserAuthService(config *config.Config, c container.Container) service.AuthServicer {
	return &UserAuthService{
		config:  config,
		userSvc: c.UserServicer,
		reader:  c.RepositoryUserAuthReader,
		writer:  c.RepositoryUserAuthWriter,
	}
}

func (s *UserAuthService) ValidateToken(ctx context.Context, request dto.RequestValidateToken) (dto.ResponseValidateKey, error) {
	userToken, err := s.reader.FindTokenBy(ctx, entity.FilterUserAuth{
		AccessToken: request.Token,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[ValidateToken.FindTokenByAccessToken] error find user access token")
		return dto.ResponseValidateKey{}, errors.New("err: token doesn't found")
	}
	if userToken == nil {
		return dto.ResponseValidateKey{}, errors.New("err: token doesn't found")
	}

	switch request.TokenType {
	case entity.AccessToken:
		now := time.Now()
		if userToken.AccessTokenExpiredAt.Before(now) {
			return dto.ResponseValidateKey{}, errors.New("err: token has expired")
		}
	case entity.RefreshToken:
		now := time.Now()
		if userToken.RefreshTokenExpiredAt.Before(now) {
			return dto.ResponseValidateKey{}, errors.New("err: refresh token has expired")
		}
	default:
		return dto.ResponseValidateKey{}, errors.New("err: token type doesn't valid")

	}

	return dto.ResponseValidateKey{
		Valid: true,
	}, nil
}

func (s *UserAuthService) GenerateToken(ctx context.Context, request dto.RequestGenerateToken) (dto.ResponseGenerateToken, error) {
	user, err := s.userSvc.InternalFindUser(ctx, userdto.RequestUser{
		Username: request.Key,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[ValidateToken.FindTokenByAccessToken] error find user access token")
		return dto.ResponseGenerateToken{}, internalerrors.NewError(err)
	}
	if user.Empty() {
		log.Warn().
			Err(err).
			Msg("[ValidateToken.FindTokenByAccessToken] error find user access token user not found")
		return dto.ResponseGenerateToken{}, internalerrors.NewError(errors.New("err: user not found"), internalerrors.SetErrorCode(http.StatusNotFound))
	}

	if err := user.ComparePassword(request.Secret); err != nil {
		log.Warn().
			Err(err).
			Msg("[ValidateToken.ComparePassword] error password didn't match")
		return dto.ResponseGenerateToken{}, internalerrors.NewError(errors.New("err: password didn't match"), internalerrors.SetErrorCode(http.StatusNotFound))
	}

	userToken, _ := s.reader.FindTokenBy(ctx, entity.FilterUserAuth{
		UserId: user.Id,
	})

	if userToken != nil && !userToken.AccessTokenExpired(time.Now()) {
		return dto.ResponseGenerateToken{
			AccessToken:     userToken.AccessToken,
			RefreshToken:    userToken.RefreshToken,
			AccessTokenExp:  userToken.AccessTokenExpiredAt.String(),
			RefreshTokenExp: userToken.RefreshTokenExpiredAt.String(),
		}, nil
	}

	localUserAccessToken := user.LocalUserAuthToken(entity.AccessToken, s.config.Auth.User.AccessTokenDuration)
	localUserRefreshToken := user.LocalUserAuthToken(entity.AccessToken, s.config.Auth.User.RefreshTokenDuration)

	accessToken, err := localUserAccessToken.GenerateToken(s.config.Auth.User.AccessTokenKey)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[ValidateToken.GenerateToken] error generate access token")
		return dto.ResponseGenerateToken{}, internalerrors.NewError(err)
	}

	refreshToken, err := localUserRefreshToken.GenerateToken(s.config.Auth.User.RefreshTokenKey)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[ValidateToken.GenerateToken] error generate refresh token")
		return dto.ResponseGenerateToken{}, internalerrors.NewError(err)
	}

	userAuth := entity.CreateUserAuth(user.Id)
	userAuth.AccessToken = accessToken
	userAuth.RefreshToken = refreshToken
	userAuth.AccessTokenExpiredAt = localUserAccessToken.Exp
	userAuth.RefreshTokenExpiredAt = localUserRefreshToken.Exp

	err = s.writer.CreateUserToken(ctx, userAuth)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[ValidateToken.CreateUserToken] error save user auth to db")
		return dto.ResponseGenerateToken{}, internalerrors.NewError(err)
	}

	return dto.ResponseGenerateToken{
		AccessToken:     userAuth.AccessToken,
		RefreshToken:    userAuth.RefreshToken,
		AccessTokenExp:  userAuth.AccessTokenExpiredAt.String(),
		RefreshTokenExp: userAuth.RefreshTokenExpiredAt.String(),
	}, nil
}

func (s *UserAuthService) ExtractToken(context.Context, dto.RequestValidateToken) (dto.ResponseExtractedToken, error) {
	return dto.ResponseExtractedToken{}, nil
}