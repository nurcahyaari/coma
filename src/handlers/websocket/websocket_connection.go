package websocket

import (
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
	clients        map[string]Client
	client         chan Client
	clientsRemoved chan []string
	close          chan bool
}

func NewWebsocketConnection() *WebsocketConnection {
	websocketConnection := &WebsocketConnection{
		clients:        make(map[string]Client),
		client:         make(chan Client),
		close:          make(chan bool),
		clientsRemoved: make(chan []string),
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
	clientId string
}

type broadcastOption func(c *broadcast)

func SetClientId(clientId string) broadcastOption {
	return func(c *broadcast) {
		c.clientId = clientId
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
		// if specificClient.clientId != "" && client.ApiToken != specificClient.clientId {
		// 	continue
		// }

		err = websocket.Message.Send(client.Connection, message)
		if err != nil {
			clientIdsErr = append(clientIdsErr, id)
		}
	}

	return clientIdsErr, err
}
