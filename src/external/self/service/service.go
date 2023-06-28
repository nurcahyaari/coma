package service

import (
	"github.com/coma/coma/config"
	"github.com/coma/coma/src/external/self/dto"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type WSServicer interface {
	Connect() error
	Send(req dto.RequestSendMessage) error
}

type WSService struct {
	ws *websocket.Conn
}

func New() WSServicer {
	wsService := &WSService{}

	return wsService
}

func (w *WSService) Connect() error {
	conn, err := websocket.Dial(config.Get().External.Coma.Websocket.Url, "", config.Get().External.Coma.Websocket.OriginUrl)
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

func (w *WSService) Send(req dto.RequestSendMessage) error {
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
