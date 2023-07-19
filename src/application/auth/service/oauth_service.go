package service

import (
	"context"
	"fmt"

	"github.com/coma/coma/config"
	"github.com/coma/coma/container"
	"github.com/coma/coma/src/application/auth/dto"
	"github.com/coma/coma/src/domains/repository"
	"github.com/coma/coma/src/domains/service"
)

type OauthService struct {
	config           *config.Config
	repositoryReader repository.RepositoryAuthReader
	repositoryWriter repository.RepositoryAuthWriter
}

func NewOauth(config *config.Config, c container.Container) service.AuthServicer {
	svc := &OauthService{
		config:           config,
		repositoryReader: c.Repository.RepositoryAuthReader,
		repositoryWriter: c.Repository.RepositoryAuthWriter,
	}
	return svc
}

func (s *OauthService) ValidateToken(ctx context.Context, request dto.RequestValidateToken) (dto.ResponseValidateKey, error) {
	return dto.ResponseValidateKey{}, fmt.Errorf("err: method is not implemented yet")
}
