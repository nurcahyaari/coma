package service

import (
	"context"

	"github.com/coma/coma/config"
	"github.com/coma/coma/container"
	"github.com/coma/coma/src/application/auth/dto"
	"github.com/coma/coma/src/domains/repository"
	"github.com/coma/coma/src/domains/service"

	"github.com/rs/zerolog/log"
)

type Service struct {
	config           *config.Config
	repositoryReader repository.RepositoryAuthReader
	repositoryWriter repository.RepositoryAuthWriter
	authSvcs         map[dto.Method]service.AuthServicer
}

func New(config *config.Config, c container.Container) service.AuthServicer {
	svc := &Service{
		config:           config,
		repositoryReader: c.Repository.RepositoryAuthReader,
		repositoryWriter: c.Repository.RepositoryAuthWriter,
		authSvcs: map[dto.Method]service.AuthServicer{
			dto.Apikey: c.ApiKeyServicer,
			dto.Oauth:  c.AuthServicer,
		},
	}
	return svc
}

func (s *Service) ValidateToken(ctx context.Context, req dto.RequestValidateToken) (dto.ResponseValidateKey, error) {
	method, err := req.Method.String()
	if err != nil {
		log.Error().
			Err(err).
			Str("Method", method).
			Msg("[ValidateToken] method is not found")
		return dto.ResponseValidateKey{}, err
	}

	res, err := s.authSvcs[req.Method].ValidateToken(ctx, dto.RequestValidateToken{
		Token: req.Token,
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("Method", method).
			Msg("[ValidateToken] err on validate token")
		return dto.ResponseValidateKey{}, err
	}

	return res, nil
}
