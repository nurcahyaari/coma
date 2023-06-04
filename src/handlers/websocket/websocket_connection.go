package websocket

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type Client struct {
	Connection *websocket.Conn
	ApiToken   string
}

type ContentType string

var (
	StringContent ContentType = "string"
	JSONContent   ContentType = "JSON"
)

type WebsocketConnection struct {
	clients                   map[string]Client
	client                    chan Client
	clientsRemoved            chan []string
	close                     chan bool
	mapBroadcastByContentType map[ContentType]func(message []byte, opts ...broadcastOption) ([]string, error)
}

func NewWebsocketConnection() *WebsocketConnection {
	websocketConnection := &WebsocketConnection{
		clients:                   make(map[string]Client),
		client:                    make(chan Client),
		close:                     make(chan bool),
		clientsRemoved:            make(chan []string),
		mapBroadcastByContentType: make(map[ContentType]func(message []byte, opts ...broadcastOption) ([]string, error)),
	}

	// registering the method handling for broadcasting message to the client
	websocketConnection.mapBroadcastByContentType = map[ContentType]func(message []byte, opts ...broadcastOption) ([]string, error){
		JSONContent: func(message []byte, opts ...broadcastOption) ([]string, error) {
			return websocketConnection.broadcastJSON(message, opts...)
		},
		StringContent: func(message []byte, opts ...broadcastOption) ([]string, error) {
			return websocketConnection.broadcastMessage(string(message), opts...)
		},
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

	log.Info().
		Str("clientId", clientId).
		Msg("add client")
}

func (w *WebsocketConnection) removeClient(clientId string) {
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
	clientId    string
	contentType ContentType
}

type broadcastOption func(c *broadcast)

func SetClientId(clientId string) broadcastOption {
	return func(c *broadcast) {
		c.clientId = clientId
	}
}

func SetContentType(contentType ContentType) broadcastOption {
	return func(c *broadcast) {
		c.contentType = contentType
	}
}

func (w *WebsocketConnection) broadcast(contentType ContentType) (func(message []byte, opts ...broadcastOption) ([]string, error), error) {
	var broadcastFunc func(message []byte, opts ...broadcastOption) ([]string, error)

	if w.mapBroadcastByContentType == nil {
		return nil, errors.New("function is not defined")
	}

	broadcastFunc = w.mapBroadcastByContentType[contentType]

	return broadcastFunc, nil
}

func (w *WebsocketConnection) broadcastMessage(message string, opts ...broadcastOption) ([]string, error) {
	var (
		clientIdsErr   []string
		err            error
		specificClient broadcast
	)

	for _, opt := range opts {
		opt(&specificClient)
	}

	for id, client := range w.clients {
		if specificClient.clientId != "" && client.ApiToken != specificClient.clientId {
			continue
		}

		err = websocket.Message.Send(client.Connection, message)
		if err != nil {
			clientIdsErr = append(clientIdsErr, id)
		}
	}

	return clientIdsErr, err
}

func (w *WebsocketConnection) broadcastJSON(message json.RawMessage, opts ...broadcastOption) ([]string, error) {
	var (
		clientIdsErr   []string
		err            error
		specificClient broadcast
	)

	for _, opt := range opts {
		opt(&specificClient)
	}

	for id, client := range w.clients {
		if specificClient.clientId != "" && client.ApiToken != specificClient.clientId {
			continue
		}

		err = websocket.JSON.Send(client.Connection, message)
		if err != nil {
			clientIdsErr = append(clientIdsErr, id)
		}
	}

	return clientIdsErr, err
}
