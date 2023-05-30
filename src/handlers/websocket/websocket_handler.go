package websocket

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type WebsocketHandler struct {
	connection *WebsocketConnection
}

func (w WebsocketHandler) Router(r *chi.Mux) {
	r.Handle("/distribute", websocket.Handler(w.DistributeData))
}

func NewWebsocketHandler() *WebsocketHandler {
	websocketHandler := &WebsocketHandler{
		connection: NewWebsocketConnection(),
	}

	go websocketHandler.connection.establishConn()

	return websocketHandler
}

func (w *WebsocketHandler) DistributeData(c *websocket.Conn) {
	token := c.Request().URL.Query().Get("authorization")
	fmt.Println(token)
	// add client
	w.connection.client <- c

	defer func() {
		w.connection.clientRemoved <- c
	}()

	for {
		msg := ""
		err := websocket.Message.Receive(c, &msg)
		if err != nil {
			log.Error().
				Err(err).
				Msg("[DistributeData] err: marshaling")
			return
		}

		log.Info().
			Str("message", msg).
			Msg("[DistributeData] received message")

		clients, err := w.connection.broadcastMessage(msg)
		if err != nil {
			w.connection.clientsRemoved <- clients
			log.Error().
				Err(err).
				Msg("[DistributeData] err: send message")
			return
		}

		log.Info().
			Str("message", string(msg)).
			Msg("[DistributeData] success send message")
	}

}
