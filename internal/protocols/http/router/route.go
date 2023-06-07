package router

import (
	"github.com/coma/coma/src/handlers/http"
	websockethandler "github.com/coma/coma/src/handlers/websocket"
	"github.com/go-chi/chi/v5"
	httpswagger "github.com/swaggo/http-swagger"
)

type HttpRouterImpl struct {
	handler   *http.HttpHandlerImpl
	wsHandler *websockethandler.WebsocketHandler
}

func (h *HttpRouterImpl) Router(r *chi.Mux) {
	h.handler.Router(r)
	h.wsHandler.Router(r)

	r.Mount("/swagger", httpswagger.WrapHandler)
}

func (h *HttpRouterImpl) CloseWebsocket() {
	h.wsHandler.Close()
}

func NewHttpRouter(
	handler *http.HttpHandlerImpl,
	wsHandler *websockethandler.WebsocketHandler,
) *HttpRouterImpl {
	return &HttpRouterImpl{
		handler:   handler,
		wsHandler: wsHandler,
	}
}
