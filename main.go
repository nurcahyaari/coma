//go:generate go run github.com/swaggo/swag/cmd/swag init

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/coma/coma/config"
	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/internal/graceful"
	"github.com/coma/coma/internal/logger"
	"github.com/coma/coma/internal/protocols/http"
	httprouter "github.com/coma/coma/internal/protocols/http/router"
	applicationrepo "github.com/coma/coma/src/domains/application/repository"
	applicationsvc "github.com/coma/coma/src/domains/application/service"
	"github.com/coma/coma/src/domains/auth/dto"
	authrepo "github.com/coma/coma/src/domains/auth/repository"
	authsvc "github.com/coma/coma/src/domains/auth/service"

	selfextsvc "github.com/coma/coma/src/external/self/service"
	httphandler "github.com/coma/coma/src/handlers/http"
	websockethandler "github.com/coma/coma/src/handlers/websocket"
)

func initHttpProtocol(
	authSvc authsvc.Servicer,
	configurationSvc applicationsvc.ApplicationConfigurationServicer,
	applicationStageSvc applicationsvc.ApplicationStageServicer,
	applicationSvc applicationsvc.ApplicationServicer,
	applicationKeySvc applicationsvc.ApplicationKeyServicer) *http.Http {
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

	// init database
	wd, _ := os.Getwd()
	cloverDB := database.NewClover(database.Config{
		Path: fmt.Sprintf("%s/%s", wd, config.Get().DB.Clover.Path),
		Name: config.Get().DB.Clover.Name,
	})

	distributorExtSvc := selfextsvc.New()

	authRepo := authrepo.New(cloverDB)
	authSvc := authsvc.New(authsvc.SetRepository(authRepo.NewRepositoryReader(), authRepo.NewRepositoryWriter()),
		authsvc.SetAuthSvc(map[dto.Method]authsvc.AuthServicer{
			dto.Apikey: authsvc.NewApiKey(authsvc.SetApiKeyRepository(authRepo.NewRepositoryReader(), authRepo.NewRepositoryWriter())),
			dto.Oauth:  authsvc.NewOauth(authsvc.SetOauthRepository(authRepo.NewRepositoryReader(), authRepo.NewRepositoryWriter())),
		}))

	applicationRepo := applicationrepo.New(cloverDB)

	applicationStageSvc := applicationsvc.NewApplicationStage(
		applicationsvc.SetApplicationStageRepository(applicationRepo))

	applicationSvc := applicationsvc.NewApplication(
		applicationsvc.SetApplicationRepository(applicationRepo))

	applicationKeySvc := applicationsvc.NewApplicationKey(
		applicationsvc.SetApplicationKeyRepository(applicationRepo))

	configurationSvc := applicationsvc.NewApplicationConfiguration(
		applicationsvc.SetApplicationConfigurationExternalService(distributorExtSvc),
		applicationsvc.SetApplicationConfigurationRepository(applicationRepo),
		applicationsvc.SetApplicationConfigurationInternalService(applicationKeySvc),
	)

	httpProtocol := initHttpProtocol(
		authSvc,
		configurationSvc,
		applicationStageSvc,
		applicationSvc,
		applicationKeySvc)

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
				"http": httpProtocol.Shutdown,
			},
		},
	)
}
