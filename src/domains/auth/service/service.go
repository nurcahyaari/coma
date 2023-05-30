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
	repo     repository.RepositoryReader
	authSvcs map[dto.Method]AuthServicer
}

func New(repo repository.RepositoryReader, svcs map[dto.Method]AuthServicer) Servicer {
	return &Service{
		repo:     repo,
		authSvcs: svcs,
	}
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
