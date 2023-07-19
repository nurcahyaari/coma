package container

import (
	"github.com/coma/coma/infrastructure/integration/coma"
	"github.com/coma/coma/src/domains/repository"
	"github.com/coma/coma/src/domains/service"
)

type Repository struct {
	repository.RepositoryAuthReader
	repository.RepositoryAuthWriter
	repository.AuthRepositorier
	repository.RepositoryApplicationWriter
	repository.RepositoryApplicationReader
	repository.RepositoryApplicationStageReader
	repository.RepositoryApplicationStageWriter
	repository.RepositoryApplicationKeyWriter
	repository.RepositoryApplicationKeyReader
	repository.RepositoryApplicationConfigurationWriter
	repository.RepositoryApplicationConfigurationReader
}

type Service struct {
	service.ApplicationConfigurationServicer
	service.ApplicationKeyServicer
	service.ApplicationStageServicer
	service.ApplicationServicer
	service.ApiKeyServicer
	service.AuthServicer
}

type Integration struct {
	*coma.WebsocketClient
}

type Container struct {
	Repository
	Service
	Integration
}
