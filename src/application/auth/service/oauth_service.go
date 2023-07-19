package service

import (
	"context"
	"fmt"

	"github.com/coma/coma/src/application/auth/dto"
	"github.com/coma/coma/src/domains/repository"
	"github.com/coma/coma/src/domains/service"
)

type OauthService struct {
	repositoryReader repository.RepositoryAuthReader
	repositoryWriter repository.RepositoryAuthWriter
}

type OauthServiceOption func(s *OauthService)

func SetOauthRepository(repositoryReader repository.RepositoryAuthReader, repositoryWriter repository.RepositoryAuthWriter) OauthServiceOption {
	return func(s *OauthService) {
		s.repositoryWriter = repositoryWriter
		s.repositoryReader = repositoryReader
	}
}

func NewOauth(opts ...OauthServiceOption) service.AuthServicer {
	svc := &OauthService{}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

func (s *OauthService) ValidateToken(ctx context.Context, request dto.RequestValidateToken) (dto.ResponseValidateKey, error) {
	return dto.ResponseValidateKey{}, fmt.Errorf("err: method is not implemented yet")
}
