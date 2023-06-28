package websocket

import (
	"encoding/json"
	"io"
	"net/http"

	configuratorsvc "github.com/coma/coma/src/domains/configurator/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type WebsocketHandler struct {
	connection      *WebsocketConnection
	configuratorSvc configuratorsvc.Servicer
}

func (h WebsocketHandler) Router(r *chi.Mux) {
	r.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		s := websocket.Server{
			Handler: websocket.Handler(h.Websocket),
		}

		s.ServeHTTP(w, r)
	})
}

type WebsockethandlerOption func(h *WebsocketHandler)

func SetDomains(configuratorSvc configuratorsvc.Servicer) WebsockethandlerOption {
	return func(h *WebsocketHandler) {
		h.configuratorSvc = configuratorSvc
	}
}

func NewWebsocketHandler(opts ...WebsockethandlerOption) *WebsocketHandler {
	websocketHandler := &WebsocketHandler{}

	for _, opt := range opts {
		opt(websocketHandler)
	}

	websocketHandler.connection = NewWebsocketConnection(
		SetWebsocketConnectionDomains(websocketHandler.configuratorSvc))

	go websocketHandler.connection.establishConn()

	return websocketHandler
}

func (w *WebsocketHandler) Close() {
	log.Warn().Msg("Clossing websocket connection")
	w.connection.close <- true
}

func (w *WebsocketHandler) Websocket(c *websocket.Conn) {
	// add client
	clientKey := c.Request().URL.Query().Get("authorization")
	w.connection.client <- Client{
		Connection: c,
		ClientKey:  clientKey,
	}

	for {
		var (
			msg  string
			data RequestDistribute
		)

		err := websocket.Message.Receive(c, &msg)
		if err != nil {
			if err == io.EOF {
				log.Error().
					Err(err).
					Msg("[Websocket] connection is closed")
				break
			}
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
			SetClientKey(data.ClientKey))
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
