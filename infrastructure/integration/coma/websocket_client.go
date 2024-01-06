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
	cfg config.Config
}

func New(cfg config.Config) *WebsocketClient {
	return &WebsocketClient{
		url: cfg.External.Coma.Websocket.Url,
		cfg: cfg,
	}
}

func (w *WebsocketClient) Connect(wait chan bool) error {
	cfg := config.New("")
	conn, err := websocket.Dial(
		fmt.Sprintf("%s?self=true", w.url),
		"",
		cfg.External.Coma.Websocket.OriginUrl)
	if err != nil {
		log.Error().
			Err(err).
			Msg("error when create websocket connection")
		return err
	}
	log.Info().Msg("websocket external connected")
	w.ws = conn

	wait <- true
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
