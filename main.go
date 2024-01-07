//go:generate go run github.com/swaggo/swag/cmd/swag init

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/coma/coma/config"
	"github.com/coma/coma/container"
	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/infrastructure/integration/coma"
	"github.com/coma/coma/internal/graceful"
	"github.com/coma/coma/internal/logger"
	"github.com/coma/coma/internal/protocols/http"
	httprouter "github.com/coma/coma/internal/protocols/http/router"
	"github.com/coma/coma/internal/x/file"
	"github.com/coma/coma/internal/x/pubsub"
	applicationrepo "github.com/coma/coma/src/application/application/repository"
	applicationsvc "github.com/coma/coma/src/application/application/service"
	authrepo "github.com/coma/coma/src/application/auth/repository"
	authsvc "github.com/coma/coma/src/application/auth/service"
	userrepo "github.com/coma/coma/src/application/user/repository"
	usersvc "github.com/coma/coma/src/application/user/service"
	"github.com/rs/zerolog/log"

	httphandler "github.com/coma/coma/src/handlers/http"
	"github.com/coma/coma/src/handlers/localpubsub"
	websockethandler "github.com/coma/coma/src/handlers/websocket"
)

//@securityDefinitions.apikey comaStandardAuth
//@in header
//@name Authorization

func initHttpProtocol(cfg config.Config, c container.Service) *http.Http {
	handler := httphandler.NewHttpHandler(c)

	websocketHandler := websockethandler.NewWebsocketHandler(c)
	router := httprouter.NewHttpRouter(
		handler,
		websocketHandler)
	return http.New(cfg, router)
}

func initDependencies(cfg config.Config) container.Container {
	var (
		c = container.Container{
			Repository:  &container.Repository{},
			Service:     &container.Service{},
			Integration: &container.Integration{},
			Event:       &container.Event{},
		}
	)

	log.Info().Msgf("initialization database on path: %s", cfg.DB.Clover.Path)
	cloverDB := database.NewClover(database.Config{
		Path: cfg.DB.Clover.Path,
		Name: cfg.DB.Clover.Name,
	})

	pubsub := pubsub.NewPubsub(pubsub.SetCloverForBackup(cloverDB.DB))

	containerEvent := container.Event{
		LocalPubsub: pubsub,
	}
	c.Event = &containerEvent

	distributorExtSvc := coma.New(cfg)

	authRepo := authrepo.New(cloverDB)
	applicationRepo := applicationrepo.New(cloverDB)
	userRepo := userrepo.New(cloverDB)

	containerRepo := container.Repository{
		RepositoryAuthReader:                     authRepo.NewRepositoryReader(),
		RepositoryAuthWriter:                     authRepo.NewRepositoryWriter(),
		AuthRepositorier:                         authRepo,
		RepositoryApplicationWriter:              applicationRepo.NewRepositoryApplicationWriter(),
		RepositoryApplicationReader:              applicationRepo.NewRepositoryApplicationReader(),
		RepositoryApplicationStageReader:         applicationRepo.NewRepositoryApplicationStageReader(),
		RepositoryApplicationStageWriter:         applicationRepo.NewRepositoryApplicationStageWriter(),
		RepositoryApplicationKeyWriter:           applicationRepo.NewRepositoryApplicationKeyWriter(),
		RepositoryApplicationKeyReader:           applicationRepo.NewRepositoryApplicationKeyReader(),
		RepositoryApplicationConfigurationWriter: applicationRepo.NewRepositoryApplicationConfigurationWriter(),
		RepositoryApplicationConfigurationReader: applicationRepo.NewRepositoryApplicationConfigurationReader(),
		RepositoryUserWriter:                     userRepo.NewRepositoryUserWriter(),
		RepositoryUserReader:                     userRepo.NewRepositoryUserReader(),
		RepositoryUserApplicationScopeWriter:     userRepo.NewRepositoryUserApplicationScopeWriter(),
		RepositoryUserApplicationScopeReader:     userRepo.NewRepositoryUserApplicationScopeReader(),
		RepositoryUserAuthReader:                 authRepo.NewRepositoryUserAuthReader(),
		RepositoryUserAuthWriter:                 authRepo.NewRepositoryUserAuthWriter(),
	}
	if err := containerRepo.Validate(); err != nil {
		log.Fatal().Errs("error", err).Msg("container repository")
	}

	containerIntegration := container.Integration{
		Coma: distributorExtSvc,
	}

	c.Repository = &containerRepo
	c.Integration = &containerIntegration

	applicationStageSvc := applicationsvc.NewApplicationStage(&cfg, c)
	c.Service.ApplicationStageServicer = applicationStageSvc

	applicationSvc := applicationsvc.NewApplication(&cfg, c)
	c.Service.ApplicationServicer = applicationSvc

	applicationKeySvc := applicationsvc.NewApplicationKey(&cfg, c)
	c.Service.ApplicationKeyServicer = applicationKeySvc

	configurationSvc := applicationsvc.NewApplicationConfiguration(&cfg, c)
	c.Service.ApplicationConfigurationServicer = configurationSvc

	userSvc := usersvc.NewUserService(&cfg, c)
	c.Service.UserServicer = userSvc
	c.Service.InternalUserServicer = userSvc

	userApplicationScopeSvc := usersvc.NewUserApplicationScopeService(&cfg, c)
	c.Service.UserApplicationScopeServicer = userApplicationScopeSvc
	c.Service.InternalUserApplicationScopeServicer = userApplicationScopeSvc

	userAuthSvc := authsvc.NewUserAuthService(&cfg, c)
	c.Service.AuthServicer = userAuthSvc
	c.Service.LocalUserAuthServicer = userAuthSvc

	if err := c.Validate(); err != nil {
		log.Fatal().Errs("error", err).Msg("container service")
	}

	return c
}

func main() {
	logger.InitLogger()
	goos := runtime.GOOS

	log.Info().Msgf("Running on operating system: %s\n", goos)
	fmt.Println("development: ", os.Getenv("development"))

	wd, _ := os.Getwd()

	// init base dir
	if err := file.NewDir(config.GetBaseWorkingDir(wd)); err != nil {
		log.Fatal().Err(err).
			Msg("creating base directory")
	}

	// init database
	cfg := config.New(wd)

	// creating database
	if err := file.NewDir(cfg.DB.Clover.Path); err != nil {
		log.Fatal().Err(err).
			Msg("creating access database directory")
	}

	c := initDependencies(cfg)

	localPubsubHandler := localpubsub.NewLocalPubsub(&cfg, c)

	httpProtocol := initHttpProtocol(cfg, *c.Service)

	// init http protocol
	go httpProtocol.Listen()

	// init other protocols here
	go c.Integration.Coma.Connect()

	localPubsubHandler.TopicRegistry()

	// listen local pubsub
	go localPubsubHandler.Listen()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	graceful.GracefulShutdown(
		ctx,
		graceful.RequestGraceful{
			ShutdownPeriod: cfg.Application.GracefulShutdownPeriod,
			Operations: map[string]graceful.Operation{
				// place your service that need to graceful shutdown here
				"http":        httpProtocol.Shutdown,
				"localPubsub": c.Event.LocalPubsub.Shutdown,
			},
		},
	)
}
