package websocket

import (
	"context"

	applicationsvc "github.com/coma/coma/src/domains/application/service"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type Client struct {
	Connection *websocket.Conn
	ClientKey  string
}

type ContentType string

var (
	StringContent ContentType = "string"
	JSONContent   ContentType = "JSON"
)

type WebsocketConnection struct {
	clients         map[string]Client
	client          chan Client
	clientsRemoved  chan []string
	close           chan bool
	configuratorSvc applicationsvc.ApplicationConfigurationServicer
}

type WebsocketConnectionOption func(h *WebsocketConnection)

func SetWebsocketConnectionDomains(configuratorSvc applicationsvc.ApplicationConfigurationServicer) WebsocketConnectionOption {
	return func(h *WebsocketConnection) {
		h.configuratorSvc = configuratorSvc
	}
}

func NewWebsocketConnection(opts ...WebsocketConnectionOption) *WebsocketConnection {
	websocketConnection := &WebsocketConnection{
		clients:        make(map[string]Client),
		client:         make(chan Client),
		close:          make(chan bool),
		clientsRemoved: make(chan []string),
	}

	for _, opt := range opts {
		opt(websocketConnection)
	}

	return websocketConnection
}

func (w *WebsocketConnection) establishConn() {
	for {
		select {
		case client := <-w.client:
			w.createClient(client)
		case <-w.close:
			w.removeAllClient()
		case clientsRemoved := <-w.clientsRemoved:
			w.removeClients(clientsRemoved)
		}
	}
}

func (w *WebsocketConnection) createClient(c Client) {
	clientId := uuid.New().String()
	w.clients[clientId] = c

	if c.ClientKey != "" {
		w.sendInitialData(c.ClientKey)
	}

	log.Info().
		Str("clientId", clientId).
		Msg("add client")
}

func (w *WebsocketConnection) sendInitialData(clientKey string) {
	if err := w.configuratorSvc.DistributeConfiguration(context.Background(), clientKey); err != nil {
		log.Warn().Msg("sendInitialData failed due to its data is empty")
	}
}

func (w *WebsocketConnection) removeClient(clientId string) {
	w.clients[clientId].Connection.Close()
	delete(w.clients, clientId)
}

func (w *WebsocketConnection) removeAllClient() {
	for id, _ := range w.clients {
		w.removeClient(id)
	}
}

func (w *WebsocketConnection) removeClients(clientIds []string) {
	for _, clientId := range clientIds {
		w.removeClient(clientId)
	}
}

type broadcast struct {
	clientKey string
}

type broadcastOption func(c *broadcast)

func SetClientKey(clientKey string) broadcastOption {
	return func(c *broadcast) {
		c.clientKey = clientKey
	}
}

func (w *WebsocketConnection) broadcast(message []byte, opts ...broadcastOption) ([]string, error) {
	var (
		clientIdsErr   []string
		err            error
		specificClient broadcast
	)

	for _, opt := range opts {
		opt(&specificClient)
	}

	for id, client := range w.clients {
		if specificClient.clientKey != "" && client.ClientKey != specificClient.clientKey {
			continue
		}

		err = websocket.Message.Send(client.Connection, message)
		if err != nil {
			clientIdsErr = append(clientIdsErr, id)
		}
	}

	return clientIdsErr, err
}
