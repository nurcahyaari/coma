package service

import (
	"context"

	"github.com/coma/coma/src/application/auth/dto"
	"github.com/coma/coma/src/domains/repository"
	"github.com/coma/coma/src/domains/service"

	"github.com/rs/zerolog/log"
)

type Servicer interface {
	ValidateToken(context.Context, dto.RequestValidateToken) (dto.ResponseValidateKey, error)
}

type Service struct {
	repositoryReader repository.RepositoryAuthReader
	repositoryWriter repository.RepositoryAuthWriter
	authSvcs         map[dto.Method]service.AuthServicer
}

type ServiceOption func(s *Service)

func SetRepository(repositoryReader repository.RepositoryAuthReader, repositoryWriter repository.RepositoryAuthWriter) ServiceOption {
	return func(s *Service) {
		s.repositoryReader = repositoryReader
		s.repositoryWriter = repositoryWriter
	}
}

func SetAuthSvc(authSvcs map[dto.Method]service.AuthServicer) ServiceOption {
	return func(s *Service) {
		s.authSvcs = authSvcs
	}
}

func New(opts ...ServiceOption) Servicer {
	svc := &Service{}

	for _, opt := range opts {
		opt(svc)
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

	res, err := s.authSvcs[req.Method].ValidateToken(ctx, dto.RequestAuthValidate{
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
