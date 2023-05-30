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
	"github.com/coma/coma/src/domains/auth/dto"
	"github.com/coma/coma/src/domains/auth/repository"
	"github.com/coma/coma/src/domains/auth/service"
	httphandler "github.com/coma/coma/src/handlers/http"
	websockethandler "github.com/coma/coma/src/handlers/websocket"
)

func initHttpProtocol(svc service.Servicer) *http.HttpImpl {
	handler := httphandler.NewHttpHandler(svc)
	websocketHandler := websockethandler.NewWebsocketHandler()
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

	repo := repository.NewApiKey(cloverDB)

	authSvc := service.New(repo, map[dto.Method]service.AuthServicer{
		dto.Apikey: service.NewApiKey(repo),
		dto.Oauth:  service.NewOauth(repo),
	})

	httpProtocol := initHttpProtocol(authSvc)

	graceful.GracefulShutdown(
		context.TODO(),
		config.Get().Application.Graceful.MaxSecond,
		map[string]graceful.Operation{
			// place your service that need to graceful shutdown here
		},
	)

	// init http protocol
	httpProtocol.Listen()

	// init other protocols here
}
