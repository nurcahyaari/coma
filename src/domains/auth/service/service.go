package service

import (
	"context"

	"github.com/coma/coma/src/domains/auth/dto"
	"github.com/coma/coma/src/domains/auth/repository"
	"github.com/rs/zerolog/log"
)

type Servicer interface {
	ValidateToken(context.Context, dto.RequestValidateToken) (dto.ResponseValidateKey, error)
}

type Service struct {
	repositoryReader repository.RepositoryReader
	repositoryWriter repository.RepositoryWriter
	authSvcs         map[dto.Method]AuthServicer
}

type ServiceOption func(s *Service)

func SetRepository(repositoryReader repository.RepositoryReader, repositoryWriter repository.RepositoryWriter) ServiceOption {
	return func(s *Service) {
		s.repositoryReader = repositoryReader
		s.repositoryWriter = repositoryWriter
	}
}

func SetAuthSvc(authSvcs map[dto.Method]AuthServicer) ServiceOption {
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
