//go:generate go run github.com/swaggo/swag/cmd/swag init

package main

import (
	"context"
	"fmt"
	"os"

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
	"github.com/rs/zerolog/log"

	httphandler "github.com/coma/coma/src/handlers/http"
	"github.com/coma/coma/src/handlers/localpubsub"
	websockethandler "github.com/coma/coma/src/handlers/websocket"
)

func initHttpProtocol(c container.Service) *http.Http {
	handler := httphandler.NewHttpHandler(c)

	websocketHandler := websockethandler.NewWebsocketHandler(c)
	router := httprouter.NewHttpRouter(
		handler,
		websocketHandler)
	return http.NewHttp(router)
}

func main() {
	logger.InitLogger()

	c := container.Container{
		Repository:  &container.Repository{},
		Service:     &container.Service{},
		Integration: &container.Integration{},
		Event:       &container.Event{},
	}

	cfg := config.Get()

	// init database
	wd, _ := os.Getwd()
	cloverDB := database.NewClover(database.Config{
		Path: fmt.Sprintf("%s/%s", wd, cfg.DB.Clover.Path),
		Name: cfg.DB.Clover.Name,
	})

	pubsub := pubsub.NewPubsub()

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
	}
	if err := containerRepo.Validate(); err != nil {
		log.Fatal().Errs("error", err).Msg("container repository")
	}

	containerIntegration := container.Integration{
		WebsocketClient: distributorExtSvc,
	}

	c.Repository = &containerRepo
	c.Integration = &containerIntegration

	apiKeySvc := authsvc.NewApiKey(&cfg, c)
	c.Service.ApiKeyServicer = apiKeySvc

	oauthSvc := authsvc.NewOauth(&cfg, c)
	c.Service.AuthServicer = oauthSvc

	authSvc := authsvc.New(&cfg, c)
	c.Service.AuthServicer = authSvc

	applicationStageSvc := applicationsvc.NewApplicationStage(&cfg, c)
	c.Service.ApplicationStageServicer = applicationStageSvc

	applicationSvc := applicationsvc.NewApplication(&cfg, c)
	c.Service.ApplicationServicer = applicationSvc

	applicationKeySvc := applicationsvc.NewApplicationKey(&cfg, c)
	c.Service.ApplicationKeyServicer = applicationKeySvc

	configurationSvc := applicationsvc.NewApplicationConfiguration(&cfg, c)
	c.Service.ApplicationConfigurationServicer = configurationSvc

	if err := c.Service.Validate(); err != nil {
		log.Fatal().Errs("error", err).Msg("container service")
	}

	httpProtocol := initHttpProtocol(*c.Service)

	localPubsubHandler := localpubsub.NewLocalPubsub(&cfg, c)
	localPubsubHandler.TopicRegistry()

	// listen local pubsub
	go localPubsubHandler.Listen()

	// init http protocol
	go httpProtocol.Listen()

	// init other protocols here
	go distributorExtSvc.Connect()

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
