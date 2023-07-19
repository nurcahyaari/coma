package coma

import (
	"fmt"

	"github.com/coma/coma/config"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type WebsocketClient struct {
	url string
	ws  *websocket.Conn
}

func New(config Config) *WebsocketClient {
	return &WebsocketClient{
		url: config.URL,
	}
}

func (w *WebsocketClient) Connect() error {
	conn, err := websocket.Dial(
		fmt.Sprintf("%s?self=true", w.url),
		"",
		config.Get().External.Coma.Websocket.OriginUrl)
	if err != nil {
		log.Error().
			Err(err).
			Msg("error when create websocket connection")
		return err
	}
	log.Info().Msg("websocket external connected")
	w.ws = conn
	return nil
}

func (w *WebsocketClient) Send(req RequestSendMessage) error {
	message, err := req.Message()
	if err != nil {
		log.
			Error().
			Err(err).Msg("[Send.Message] error when marshaling dto")
		return err
	}
	err = websocket.Message.Send(w.ws, message)
	if err != nil {
		log.
			Error().
			Err(err).Msg("[Send] error when send data")
		return err
	}

	return nil
}
