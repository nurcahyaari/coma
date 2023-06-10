package service

import (
	"context"

	"github.com/coma/coma/src/domains/configurator/dto"
	"github.com/coma/coma/src/domains/configurator/repository"
	selfextsvc "github.com/coma/coma/src/external/self/service"
)

type Servicer interface {
	SetConfiguration(ctx context.Context, req dto.RequestSetConfiguration) error
}

type Service struct {
	selfExtSvc selfextsvc.WSServicer
	readerRepo repository.RepositoryReader
	writerRepo repository.RepositoryWriter
}

type ServiceOption func(svc *Service)

func SetExternalService(selfExtService selfextsvc.WSServicer) ServiceOption {
	return func(svc *Service) {
		svc.selfExtSvc = selfExtService
	}
}

func SetRepository(reader repository.RepositoryReader, writer repository.RepositoryWriter) ServiceOption {
	return func(svc *Service) {
		svc.readerRepo = reader
		svc.writerRepo = writer
	}
}

func New(opts ...ServiceOption) Servicer {
	svc := &Service{}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}

func (s *Service) SetConfiguration(ctx context.Context, req dto.RequestSetConfiguration) error {
	if err := req.Validate(); err != nil {
		return err
	}

	return nil
}
