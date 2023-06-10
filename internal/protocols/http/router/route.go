package router

import (
	"net/http"

	_ "github.com/coma/coma/docs"
	httphandler "github.com/coma/coma/src/handlers/http"
	websockethandler "github.com/coma/coma/src/handlers/websocket"
	"github.com/go-chi/chi/v5"
	httpswagger "github.com/swaggo/http-swagger"
)

type HttpRouterImpl struct {
	handler   *httphandler.HttpHandle
	wsHandler *websockethandler.WebsocketHandler
}

func (h *HttpRouterImpl) Router(r *chi.Mux) {
	h.handler.Router(r)
	h.wsHandler.Router(r)

	r.Mount("/swagger", httpswagger.WrapHandler)
	r.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
}

func (h *HttpRouterImpl) CloseWebsocket() {
	h.wsHandler.Close()
}

func NewHttpRouter(
	handler *httphandler.HttpHandle,
	wsHandler *websockethandler.WebsocketHandler,
) *HttpRouterImpl {
	return &HttpRouterImpl{
		handler:   handler,
		wsHandler: wsHandler,
	}
}
