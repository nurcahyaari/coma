package websocket

import (
	"encoding/json"

	distributiondto "github.com/coma/coma/src/domains/distributor/dto"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type WebsocketHandler struct {
	connection *WebsocketConnection
}

func (w WebsocketHandler) Router(r *chi.Mux) {
	r.Handle("/websocket", websocket.Handler(w.Websocket))
}

func NewWebsocketHandler() *WebsocketHandler {
	websocketHandler := &WebsocketHandler{
		connection: NewWebsocketConnection(),
	}

	go websocketHandler.connection.establishConn()

	return websocketHandler
}

func (w *WebsocketHandler) Close() {
	log.Warn().Msg("Clossing websocket connection")
	w.connection.close <- true
}

func (w *WebsocketHandler) Websocket(c *websocket.Conn) {
	// add client
	w.connection.client <- Client{
		Connection: c,
		ApiToken:   c.Request().URL.Query().Get("authorization"),
	}

	for {
		var (
			msg  string
			data distributiondto.RequestDistribute
		)

		err := websocket.Message.Receive(c, &msg)
		if err != nil {
			log.Error().
				Err(err).
				Msg("[Websocket] err: marshaling from message")
			continue
		}

		log.Info().
			Str("message", msg).
			Msg("[Websocket] received message")

		err = json.Unmarshal([]byte(msg), &data)
		if err != nil {
			log.Error().
				Err(err).
				Msg("[Websocket] err: unmarshaling to struct")
			continue
		}

		if errs := data.Validate(); len(errs) > 0 {
			log.Error().
				Errs("validate", errs).
				Msg("[Websocket] there is an error on the request object")
			continue
		}

		byt, err := json.Marshal(data)
		if err != nil {
			log.Error().
				Err(err).
				Msg("[Websocket] err: marshaling to byte")
			continue
		}

		clients, err := w.connection.broadcast(byt,
			SetClientId(data.ApiToken))
		if err != nil {
			w.connection.clientsRemoved <- clients
			log.Error().
				Err(err).
				Msg("[Websocket] err: send message")
			continue
		}

		log.Info().
			Str("message", string(msg)).
			Msg("[Websocket] success send message")
	}
}
