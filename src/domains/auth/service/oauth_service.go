package service

import (
	"context"
	"fmt"

	"github.com/coma/coma/src/domains/auth/dto"
	"github.com/coma/coma/src/domains/auth/repository"
)

type OauthService struct {
	repositoryReader repository.RepositoryReader
	repositoryWriter repository.RepositoryWriter
}

type OauthServiceOption func(s *OauthService)

func SetOauthRepository(repositoryReader repository.RepositoryReader, repositoryWriter repository.RepositoryWriter) OauthServiceOption {
	return func(s *OauthService) {
		s.repositoryWriter = repositoryWriter
		s.repositoryReader = repositoryReader
	}
}

func NewOauth(opts ...OauthServiceOption) AuthServicer {
	svc := &OauthService{}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

func (s *OauthService) ValidateToken(ctx context.Context, request dto.RequestAuthValidate) (dto.ResponseValidateKey, error) {
	return dto.ResponseValidateKey{}, fmt.Errorf("err: method is not implemented yet")
}
