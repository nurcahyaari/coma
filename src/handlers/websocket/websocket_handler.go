package websocket

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/coma/coma/container"
	internalerrors "github.com/coma/coma/internal/x/errors"
	"github.com/coma/coma/src/application/application/dto"
	"github.com/coma/coma/src/domain/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type WebsocketHandler struct {
	connection        *WebsocketConnection
	configurationSvc  service.ApplicationConfigurationServicer
	applicationKeySvc service.ApplicationKeyServicer
}

func (h WebsocketHandler) Router(r *chi.Mux) {
	r.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		s := websocket.Server{
			Handler: websocket.Handler(h.Websocket),
		}

		s.ServeHTTP(w, r)
	})
}

func NewWebsocketHandler(c container.Service) *WebsocketHandler {
	websocketHandler := &WebsocketHandler{
		connection:        NewWebsocketConnection(c),
		configurationSvc:  c.ApplicationConfigurationServicer,
		applicationKeySvc: c.ApplicationKeyServicer,
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
	selfConnection := c.Request().URL.Query().Get("self")
	isSelfConnection, _ := strconv.ParseBool(selfConnection)
	clientKey := c.Request().URL.Query().Get("authorization")
	if !isSelfConnection {
		exists, err := w.applicationKeySvc.IsExistsApplicationKey(context.Background(), dto.RequestFindApplicationKey{
			Key: clientKey,
		})
		if err != nil {
			errCustom := err.(*internalerrors.Error)
			log.Error().
				Err(errCustom.Err).
				Msg("[Websocket.FindApplicationKey] err: search applicationKey")
			return
		}
		if !exists {
			errCustom := err.(*internalerrors.Error)
			log.Error().
				Err(errCustom.Err).
				Msg("[Websocket.FindApplicationKey] client key doesn't exists")
			return
		}
	}

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
		if err == io.EOF {
			log.Warn().
				Err(err).
				Msg("[Websocket] connection is closed")
			break
		}
		if err != nil {
			log.Error().
				Err(err).
				Msg("[Websocket] err: marshaling from message")
			break
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
