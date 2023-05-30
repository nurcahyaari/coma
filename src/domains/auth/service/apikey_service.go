package service

import (
	"context"
	"fmt"

	"github.com/coma/coma/src/domains/auth/dto"
	"github.com/coma/coma/src/domains/auth/model"
	"github.com/coma/coma/src/domains/auth/repository"
	"github.com/rs/zerolog/log"
)

type ApiKeyServicer interface {
	AuthServicer
}

type ApiKeyService struct {
	repo repository.RepositoryReader
}

func NewApiKey(repo repository.RepositoryReader) ApiKeyServicer {
	return &ApiKeyService{
		repo: repo,
	}
}

func (s *ApiKeyService) ValidateToken(ctx context.Context, request dto.RequestAuthValidate) (dto.ResponseValidateKey, error) {
	var (
		resp   = dto.ResponseValidateKey{}
		apikey model.Apikey
		err    error
	)

	apikey, err = s.repo.FindTokenByToken(ctx, request.Token)
	if err != nil {
		log.Error().Err(err).Msg("[ApiKeyService][ValidateToken] err on FindTokenById")
		return resp, err
	}

	if apikey.Id == 0 {
		log.Error().Err(err).Msg("[ApiKeyService][ValidateToken] err token is not found")
		return resp, fmt.Errorf("error: token is not found")
	}

	return dto.ResponseValidateKey{
		Valid: true,
	}, nil
}
