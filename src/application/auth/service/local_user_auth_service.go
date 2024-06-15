package service

import (
	"context"
	"crypto/rsa"
	"errors"
	"net/http"
	"time"

	"github.com/nurcahyaari/coma/config"
	"github.com/nurcahyaari/coma/container"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/application/auth/dto"
	userdto "github.com/nurcahyaari/coma/src/application/user/dto"
	"github.com/nurcahyaari/coma/src/domain/entity"
	"github.com/nurcahyaari/coma/src/domain/repository"
	"github.com/nurcahyaari/coma/src/domain/service"
	"github.com/rs/zerolog/log"
)

type UserAuthService struct {
	config             *config.Config
	userSvc            service.InternalUserServicer
	userApplicationSvc service.InternalUserApplicationScopeServicer
	reader             repository.RepositoryUserAuthReader
	writer             repository.RepositoryUserAuthWriter
}

func NewUserAuthService(config *config.Config, c container.Container) service.LocalUserAuthServicer {
	return &UserAuthService{
		config:             config,
		userSvc:            c.UserServicer,
		userApplicationSvc: c.InternalUserApplicationScopeServicer,
		reader:             c.RepositoryUserAuthReader,
		writer:             c.RepositoryUserAuthWriter,
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
		if userToken.AccessTokenExpired(now) {
			return dto.ResponseValidateKey{}, errors.New("err: token has expired")
		}
	case entity.RefreshToken:
		now := time.Now()
		if userToken.RefreshTokenExpired(now) {
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

	accessToken, err := localUserAccessToken.GenerateJWTToken(s.config.Auth.User.PrivateKey)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[ValidateToken.GenerateToken] error generate access token")
		return dto.ResponseGenerateToken{}, internalerrors.NewError(err)
	}

	refreshToken := localUserRefreshToken.GenerateUuidToken()
	userAuth := entity.CreateUserAuth(user.Id)
	userAuth.AccessToken = accessToken
	userAuth.RefreshToken = refreshToken
	userAuth.AccessTokenExpiredAt = localUserAccessToken.ExpiresAt.Time
	userAuth.RefreshTokenExpiredAt = localUserRefreshToken.ExpiresAt.Time

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

func (s *UserAuthService) ExtractToken(ctx context.Context, req dto.RequestValidateToken) (dto.ResponseExtractedToken, error) {
	var key *rsa.PublicKey

	switch req.TokenType {
	case entity.AccessToken,
		entity.RefreshToken:
		key = s.config.Auth.User.PublicKey
	default:
		return dto.ResponseExtractedToken{}, errors.New("err: token type is not valid")
	}

	localUserAuthToken, err := entity.NewLocalUserAuthTokenFromToken(req.Token, key)
	if err != nil {
		log.Error().
			Err(err).
			Msg("[ExtractToken.NewLocalUserAuthTokenFromToken] error token is not valid")
		return dto.ResponseExtractedToken{}, internalerrors.NewError(err)
	}

	if !localUserAuthToken.ValidTokenType(req.TokenType) {
		log.Warn().
			Msg("[ExtractToken.ValidTokenType] error token is mismatch")
		return dto.ResponseExtractedToken{}, internalerrors.NewError(errors.New("err: token type is mismatch"))
	}

	return dto.ResponseExtractedToken{
		UserId:    localUserAuthToken.Id,
		ExpiredAt: localUserAuthToken.ExpiresAt.Time,
		UserType:  localUserAuthToken.UserType,
	}, nil
}

func (s *UserAuthService) ValidateUserScope(ctx context.Context, req dto.RequestUserScopeValidation) (dto.ResponseValidateKey, error) {
	user, err := s.userSvc.InternalFindUser(ctx, userdto.RequestUser{
		Id: req.UserId,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msg("[ValidateUserScope.InternalFindUser] error user id is not found")
		return dto.ResponseValidateKey{}, internalerrors.NewError(err)
	}

	if user.UserAdmin() {
		return dto.ResponseValidateKey{
			Valid: true,
		}, nil
	}

	if user.Rbac == nil {
		return dto.ResponseValidateKey{}, nil
	}

	if !user.HasRbacAccess(req.Method) {
		return dto.ResponseValidateKey{}, nil
	}

	return dto.ResponseValidateKey{
		Valid: true,
	}, nil
}

func (s *UserAuthService) ValidateUserApplicationScope(ctx context.Context, req dto.RequestUserApplicationScopeValidation) (dto.ResponseValidateKey, error) {
	userApplicationScope, exist, err := s.userApplicationSvc.InternalFindUserApplicationScope(ctx, userdto.RequestFindUserApplicationScope{
		UserId: req.UserId,
	})

	if err != nil {
		log.Error().
			Err(err).
			Msg("[ValidateUserApplicationScope.InternalFindUserApplicationScope] error")
		return dto.ResponseValidateKey{}, internalerrors.NewError(err)
	}
	if !exist {
		err := errors.New("err: user application scope is not exist")
		log.Error().
			Err(err).
			Msg("[ValidateUserApplicationScope.InternalFindUserApplicationScope] error user application scope is not found")
		return dto.ResponseValidateKey{}, internalerrors.NewError(err)
	}

	if userApplicationScope.Rbac == nil {
		return dto.ResponseValidateKey{}, nil
	}

	if !userApplicationScope.HasRbacAccess(req.Method) {
		return dto.ResponseValidateKey{}, nil
	}

	return dto.ResponseValidateKey{}, nil
}
