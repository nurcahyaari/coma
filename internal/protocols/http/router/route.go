package router

import (
	"github.com/go-chi/chi/v5"
	httphandler "github.com/nurcahyaari/coma/src/handlers/http"
	websockethandler "github.com/nurcahyaari/coma/src/handlers/websocket"
)

type HttpRoute struct {
	handler   *httphandler.HttpHandle
	wsHandler *websockethandler.WebsocketHandler
}

func (h *HttpRoute) Router(r *chi.Mux) {
	h.handler.Router(r)
	h.wsHandler.Router(r)

}

func (h *HttpRoute) CloseWebsocket() {
	h.wsHandler.Close()
}

func NewHttpRouter(
	handler *httphandler.HttpHandle,
	wsHandler *websockethandler.WebsocketHandler,
) *HttpRoute {
	return &HttpRoute{
		handler:   handler,
		wsHandler: wsHandler,
	}
}
