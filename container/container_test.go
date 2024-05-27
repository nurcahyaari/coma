package container_test

import (
	"testing"

	"github.com/nurcahyaari/coma/container"
	"github.com/nurcahyaari/coma/infrastructure/integration/coma"
	"github.com/nurcahyaari/coma/internal/x/pubsub"
	applicationsvc "github.com/nurcahyaari/coma/src/application/application/service"
	authsvc "github.com/nurcahyaari/coma/src/application/auth/service"
	"github.com/nurcahyaari/coma/src/domain/repository/repositoryfakes"
	"github.com/stretchr/testify/assert"
)

func TestRepositoryValidate(t *testing.T) {
	t.Run("test null", func(t *testing.T) {
		r := container.Repository{}
		err := r.Validate()
		assert.NotEmpty(t, err)
	})

	t.Run("test no error", func(t *testing.T) {
		r := container.Repository{
			RepositoryAuthReader:                     &repositoryfakes.FakeRepositoryAuthReader{},
			RepositoryAuthWriter:                     &repositoryfakes.FakeAuthRepositorier{},
			RepositoryApplicationWriter:              &repositoryfakes.FakeRepositoryApplicationWriter{},
			RepositoryApplicationReader:              &repositoryfakes.FakeRepositoryApplicationReader{},
			RepositoryApplicationKeyWriter:           &repositoryfakes.FakeRepositoryApplicationKeyWriter{},
			RepositoryApplicationKeyReader:           &repositoryfakes.FakeRepositoryApplicationKeyReader{},
			RepositoryApplicationConfigurationWriter: &repositoryfakes.FakeRepositoryApplicationConfigurationWriter{},
			RepositoryApplicationConfigurationReader: &repositoryfakes.FakeRepositoryApplicationConfigurationReader{},
			AuthRepositorier:                         &repositoryfakes.FakeAuthRepositorier{},
		}
		err := r.Validate()
		assert.Equal(t, 0, len(err))
		assert.Nil(t, err)
	})
}

func TestIntegrationValidate(t *testing.T) {
	t.Run("test null", func(t *testing.T) {
		r := container.Integration{}
		err := r.Validate()
		assert.NotEmpty(t, err)
	})

	t.Run("test no error", func(t *testing.T) {
		r := container.Integration{
			Coma: &coma.WebsocketClient{},
		}
		err := r.Validate()
		assert.Equal(t, 0, len(err))
		assert.Nil(t, err)
	})
}

func TestContainerValidate(t *testing.T) {
	t.Run("test null", func(t *testing.T) {
		r := container.Container{}
		err := r.Validate()
		assert.NotEmpty(t, err)
	})

	t.Run("test no error", func(t *testing.T) {
		r := container.Container{
			Repository: &container.Repository{
				RepositoryAuthReader:                     &repositoryfakes.FakeRepositoryAuthReader{},
				RepositoryAuthWriter:                     &repositoryfakes.FakeAuthRepositorier{},
				RepositoryApplicationWriter:              &repositoryfakes.FakeRepositoryApplicationWriter{},
				RepositoryApplicationReader:              &repositoryfakes.FakeRepositoryApplicationReader{},
				RepositoryApplicationKeyWriter:           &repositoryfakes.FakeRepositoryApplicationKeyWriter{},
				RepositoryApplicationKeyReader:           &repositoryfakes.FakeRepositoryApplicationKeyReader{},
				RepositoryApplicationConfigurationWriter: &repositoryfakes.FakeRepositoryApplicationConfigurationWriter{},
				RepositoryApplicationConfigurationReader: &repositoryfakes.FakeRepositoryApplicationConfigurationReader{},
				AuthRepositorier:                         &repositoryfakes.FakeAuthRepositorier{},
			},
			Service: &container.Service{
				ApplicationConfigurationServicer: &applicationsvc.ApplicationConfigurationService{},
				ApplicationKeyServicer:           &applicationsvc.ApplicationKeyService{},
				ApplicationServicer:              &applicationsvc.ApplicationService{},
				AuthServicer:                     &authsvc.UserAuthService{},
			},
			Integration: &container.Integration{
				Coma: &coma.WebsocketClient{},
			},
			Event: &container.Event{
				LocalPubsub: &pubsub.Pubsub{},
			},
		}
		err := r.Validate()
		assert.Equal(t, 0, len(err))
		assert.Nil(t, err)
	})
}
