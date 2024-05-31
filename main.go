//go:generate go run github.com/swaggo/swag/cmd/swag init

package main

import (
	"context"
	"runtime"

	"github.com/nurcahyaari/coma/config"
	"github.com/nurcahyaari/coma/container"
	"github.com/nurcahyaari/coma/infrastructure/database"
	"github.com/nurcahyaari/coma/infrastructure/integration/coma"
	"github.com/nurcahyaari/coma/internal/graceful"
	"github.com/nurcahyaari/coma/internal/logger"
	"github.com/nurcahyaari/coma/internal/protocols/http"
	httprouter "github.com/nurcahyaari/coma/internal/protocols/http/router"
	"github.com/nurcahyaari/coma/internal/x/pubsub"
	applicationrepo "github.com/nurcahyaari/coma/src/application/application/repository"
	applicationsvc "github.com/nurcahyaari/coma/src/application/application/service"
	authrepo "github.com/nurcahyaari/coma/src/application/auth/repository"
	authsvc "github.com/nurcahyaari/coma/src/application/auth/service"
	userrepo "github.com/nurcahyaari/coma/src/application/user/repository"
	usersvc "github.com/nurcahyaari/coma/src/application/user/service"
	"github.com/rs/zerolog/log"

	httphandler "github.com/nurcahyaari/coma/src/handlers/http"
	"github.com/nurcahyaari/coma/src/handlers/localpubsub"
	websockethandler "github.com/nurcahyaari/coma/src/handlers/websocket"
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

	applicationKeySvc := applicationsvc.NewApplicationKey(&cfg, c)
	c.Service.ApplicationKeyServicer = applicationKeySvc

	applicationSvc := applicationsvc.NewApplication(&cfg, c)
	c.Service.ApplicationServicer = applicationSvc

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// creating configuration
	if err := config.CreateBaseDir(ctx); err != nil {
		log.Fatal().Err(err).
			Msg("creating base directory")
	}

	cfg := config.New()

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
