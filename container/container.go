package container

import (
	"fmt"
	"reflect"

	"github.com/coma/coma/infrastructure/integration/coma"
	"github.com/coma/coma/internal/utils/pubsub"
	"github.com/coma/coma/src/domain/repository"
	"github.com/coma/coma/src/domain/service"
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
	repository.RepositoryUserWriter
	repository.RepositoryUserReader
	repository.RepositoryUserAuthReader
	repository.RepositoryUserAuthWriter
	repository.RepositoryUserApplicationScopeWriter
	repository.RepositoryUserApplicationScopeReader
}

func (c Repository) Validate() []error {
	errs := []error{}
	value := reflect.ValueOf(c)
	types := value.Type()
	for i := 0; i < types.NumField(); i++ {
		if value.Field(i).IsNil() {
			errs = append(errs, fmt.Errorf("%s: must not be empty", types.Field(i).Name))
		}
	}
	if len(errs) > 0 {
		return errs
	}

	return nil
}

type Service struct {
	service.ApplicationConfigurationServicer
	service.ApplicationKeyServicer
	service.ApplicationStageServicer
	service.ApplicationServicer
	service.AuthServicer
	service.LocalUserAuthServicer
	service.UserServicer
	service.InternalUserServicer
	service.UserApplicationScopeServicer
	service.InternalUserApplicationScopeServicer
}

func (c Service) Validate() []error {
	errs := []error{}
	value := reflect.ValueOf(c)
	types := value.Type()
	for i := 0; i < types.NumField(); i++ {
		if value.Field(i).IsNil() {
			errs = append(errs, fmt.Errorf("%s: must not be empty", types.Field(i).Name))
		}
	}
	if len(errs) > 0 {
		return errs
	}

	return nil
}

type Integration struct {
	*coma.WebsocketClient
}

func (c Integration) Validate() []error {
	errs := []error{}
	value := reflect.ValueOf(c)
	types := value.Type()
	for i := 0; i < types.NumField(); i++ {
		if value.Field(i).IsNil() {
			errs = append(errs, fmt.Errorf("%s: must not be empty", types.Field(i).Name))
		}
	}
	if len(errs) > 0 {
		return errs
	}

	return nil
}

type Event struct {
	LocalPubsub *pubsub.Pubsub
}

func (c Event) Validate() []error {
	errs := []error{}
	value := reflect.ValueOf(c)
	types := value.Type()
	for i := 0; i < types.NumField(); i++ {
		if value.Field(i).IsNil() {
			errs = append(errs, fmt.Errorf("%s: must not be empty", types.Field(i).Name))
		}
	}
	if len(errs) > 0 {
		return errs
	}

	return nil
}

type Container struct {
	*Repository
	*Service
	*Integration
	*Event
}

func (c Container) Validate() []error {
	errs := []error{}
	value := reflect.ValueOf(c)
	types := value.Type()
	for i := 0; i < types.NumField(); i++ {
		if value.Field(i).IsNil() {
			errs = append(errs, fmt.Errorf("%s: must not be empty", types.Field(i).Name))
		}
	}
	if len(errs) > 0 {
		return errs
	}

	if errs := c.Repository.Validate(); len(errs) > 0 {
		return errs
	}

	if errs := c.Service.Validate(); len(errs) > 0 {
		return errs
	}

	if errs := c.Integration.Validate(); len(errs) > 0 {
		return errs
	}

	if errs := c.Event.Validate(); len(errs) > 0 {
		return errs
	}

	return nil
}
