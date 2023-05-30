package websocket

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/net/websocket"
)

type WebsocketConnection struct {
	clients        map[string]*websocket.Conn
	client         chan *websocket.Conn
	clientRemoved  chan *websocket.Conn
	clientsRemoved chan []string
}

func NewWebsocketConnection() *WebsocketConnection {
	websocketConnection := &WebsocketConnection{
		clients:        make(map[string]*websocket.Conn),
		client:         make(chan *websocket.Conn),
		clientRemoved:  make(chan *websocket.Conn),
		clientsRemoved: make(chan []string),
	}

	return websocketConnection
}

func (w *WebsocketConnection) establishConn() {
	for {
		select {
		case client := <-w.client:
			w.createClient(client)
		case clientRemoved := <-w.clientRemoved:
			w.removeClient(clientRemoved)
		case clientsRemoved := <-w.clientsRemoved:
			w.removeClients(clientsRemoved)
		}
	}
}

func (w *WebsocketConnection) createClient(c *websocket.Conn) {
	remoteAddr := c.RemoteAddr().String()
	w.clients[remoteAddr] = c

	log.Info().
		Str("addr", remoteAddr).
		Msg("add client")
}

func (w *WebsocketConnection) removeClient(c *websocket.Conn) {
	delete(w.clients, c.RemoteAddr().Network())
}

func (w *WebsocketConnection) removeClients(addrs []string) {
	for _, addr := range addrs {
		w.removeClient(w.clients[addr])
	}
}

func (w *WebsocketConnection) broadcastMessage(message string) ([]string, error) {
	var (
		clientAddrErr []string
		err           error
	)
	for addr, client := range w.clients {
		err = websocket.Message.Send(client, message)
		if err != nil {
			clientAddrErr = append(clientAddrErr, addr)
		}
	}

	return clientAddrErr, err
}
