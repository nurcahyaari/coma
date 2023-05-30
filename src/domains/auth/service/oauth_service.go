package service

import (
	"context"
	"fmt"

	"github.com/coma/coma/src/domains/auth/dto"
	"github.com/coma/coma/src/domains/auth/repository"
)

type OauthServicer interface {
	AuthServicer
}

type OauthService struct {
	repo repository.RepositoryReader
}

func NewOauth(repo repository.RepositoryReader) OauthServicer {
	return &OauthService{
		repo: repo,
	}
}

func (s *OauthService) ValidateToken(ctx context.Context, request dto.RequestAuthValidate) (dto.ResponseValidateKey, error) {
	return dto.ResponseValidateKey{}, fmt.Errorf("err: method is not implemented yet")
}
