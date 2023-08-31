//go:generate go run github.com/swaggo/swag/cmd/swag init

package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/coma/coma/config"
	"github.com/coma/coma/container"
	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/infrastructure/integration/coma"
	"github.com/coma/coma/internal/graceful"
	"github.com/coma/coma/internal/logger"
	"github.com/coma/coma/internal/protocols/http"
	httprouter "github.com/coma/coma/internal/protocols/http/router"
	"github.com/coma/coma/internal/utils/pubsub"
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

func initHttpProtocol(c container.Service) *http.Http {
	handler := httphandler.NewHttpHandler(c)

	websocketHandler := websockethandler.NewWebsocketHandler(c)
	router := httprouter.NewHttpRouter(
		handler,
		websocketHandler)
	return http.New(router)
}

func main() {
	logger.InitLogger()

	var (
		cfg         = config.Get()
		wd          = "/var/lib/coma"
		storagePath = filepath.Join(wd, "coma", cfg.DB.Clover.Path)
		c           = container.Container{
			Repository:  &container.Repository{},
			Service:     &container.Service{},
			Integration: &container.Integration{},
			Event:       &container.Event{},
		}
	)

	// init database
	if cfg.Application.Development {
		wd, _ = os.Getwd()
	}
	storagePath = filepath.Join(wd, cfg.DB.Clover.Path)

	cmd := exec.Command("mkdir", "-m", "0777", "-p", filepath.Join(storagePath, cfg.DB.Clover.Name))
	cmd.Env = append(os.Environ(), "SUDO_COMMAND=true")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatal().Err(err).
			Str("path", storagePath).
			Msg("creating database directory")
	}

	cmd = exec.Command("chmod", "755", filepath.Join(storagePath, cfg.DB.Clover.Name))
	cmd.Env = append(os.Environ(), "SUDO_COMMAND=true")
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatal().Err(err).
			Str("path", storagePath).
			Msg("creating access database directory")
	}

	log.Info().Msgf("initialization database on path: %s", storagePath)
	cloverDB := database.NewClover(database.Config{
		Path: storagePath,
		Name: cfg.DB.Clover.Name,
	})

	pubsub := pubsub.NewPubsub(pubsub.SetCloverForBackup(cloverDB.DB))

	containerEvent := container.Event{
		LocalPubsub: pubsub,
	}
	c.Event = &containerEvent

	distributorExtSvc := coma.New(
		coma.Config{
			URL: config.Get().External.Coma.Websocket.Url,
		},
	)

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
		WebsocketClient: distributorExtSvc,
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

	if err := c.Service.Validate(); err != nil {
		log.Fatal().Errs("error", err).Msg("container service")
	}

	localPubsubHandler := localpubsub.NewLocalPubsub(&cfg, c)

	httpProtocol := initHttpProtocol(*c.Service)

	// init http protocol
	go httpProtocol.Listen()

	// init other protocols here
	go distributorExtSvc.Connect()

	localPubsubHandler.TopicRegistry()

	// listen local pubsub
	go localPubsubHandler.Listen()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	graceful.GracefulShutdown(
		ctx,
		graceful.RequestGraceful{
			ShutdownPeriod: config.Get().Application.Graceful.ShutdownPeriod,
			Operations: map[string]graceful.Operation{
				// place your service that need to graceful shutdown here
				"http":        httpProtocol.Shutdown,
				"localPubsub": pubsub.Shutdown,
			},
		},
	)
}
