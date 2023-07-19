//go:generate go run github.com/swaggo/swag/cmd/swag init

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coma/coma/config"
	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/infrastructure/integration/coma"
	"github.com/coma/coma/internal/graceful"
	"github.com/coma/coma/internal/logger"
	"github.com/coma/coma/internal/protocols/http"
	httprouter "github.com/coma/coma/internal/protocols/http/router"
	"github.com/coma/coma/internal/utils/pubsub"
	applicationrepo "github.com/coma/coma/src/application/application/repository"
	applicationsvc "github.com/coma/coma/src/application/application/service"
	"github.com/coma/coma/src/application/auth/dto"
	authrepo "github.com/coma/coma/src/application/auth/repository"
	authsvc "github.com/coma/coma/src/application/auth/service"
	"github.com/coma/coma/src/domains/service"

	httphandler "github.com/coma/coma/src/handlers/http"
	"github.com/coma/coma/src/handlers/localpubsub"
	websockethandler "github.com/coma/coma/src/handlers/websocket"
)

func initHttpProtocol(
	authSvc authsvc.Servicer,
	configurationSvc service.ApplicationConfigurationServicer,
	applicationStageSvc service.ApplicationStageServicer,
	applicationSvc service.ApplicationServicer,
	applicationKeySvc service.ApplicationKeyServicer) *http.Http {
	handler := httphandler.NewHttpHandler(
		httphandler.SetDomains(
			authSvc,
			configurationSvc,
			applicationStageSvc,
			applicationSvc,
			applicationKeySvc))

	websocketHandler := websockethandler.NewWebsocketHandler(websockethandler.SetDomains(
		configurationSvc,
		applicationKeySvc))
	router := httprouter.NewHttpRouter(
		handler,
		websocketHandler)
	return http.NewHttp(router)
}

func main() {
	logger.InitLogger()

	cfg := config.Get()

	// init database
	wd, _ := os.Getwd()
	cloverDB := database.NewClover(database.Config{
		Path: fmt.Sprintf("%s/%s", wd, cfg.DB.Clover.Path),
		Name: cfg.DB.Clover.Name,
	})

	pubSub := pubsub.NewPubsub()

	distributorExtSvc := coma.New()

	authRepo := authrepo.New(cloverDB)
	authSvc := authsvc.New(authsvc.SetRepository(authRepo.NewRepositoryReader(), authRepo.NewRepositoryWriter()),
		authsvc.SetAuthSvc(map[dto.Method]authsvc.AuthServicer{
			dto.Apikey: authsvc.NewApiKey(authsvc.SetApiKeyRepository(authRepo.NewRepositoryReader(), authRepo.NewRepositoryWriter())),
			dto.Oauth:  authsvc.NewOauth(authsvc.SetOauthRepository(authRepo.NewRepositoryReader(), authRepo.NewRepositoryWriter())),
		}))

	applicationRepo := applicationrepo.New(cloverDB)

	applicationStageSvc := applicationsvc.NewApplicationStage(
		&cfg,
		applicationsvc.SetApplicationStageRepository(applicationRepo))

	applicationSvc := applicationsvc.NewApplication(
		&cfg,
		applicationsvc.SetApplicationRepository(applicationRepo))

	applicationKeySvc := applicationsvc.NewApplicationKey(
		&cfg,
		applicationsvc.SetApplicationKeyRepository(applicationRepo))

	configurationSvc := applicationsvc.NewApplicationConfiguration(
		&cfg,
		applicationsvc.SetApplicationConfigurationExternalService(distributorExtSvc),
		applicationsvc.SetApplicationConfigurationRepository(applicationRepo),
		applicationsvc.SetApplicationConfigurationInternalService(applicationKeySvc),
		applicationsvc.SetApplicationConfigurationEvent(pubSub),
	)

	httpProtocol := initHttpProtocol(
		authSvc,
		configurationSvc,
		applicationStageSvc,
		applicationSvc,
		applicationKeySvc)

	localPubsubHandler := localpubsub.NewLocalPubsub(&cfg, pubSub, localpubsub.SetDomains(configurationSvc))
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
				"localPubsub": pubSub.Shutdown,
			},
		},
	)
}
