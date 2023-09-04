package websockets

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Mu   sync.Mutex
}

func (s *WsServer) AddActiveConnection(conn *websocket.Conn) *WsServer {
	s.activeClientsMu.Lock()
	defer s.activeClientsMu.Unlock()
	client := &Client{
		Conn: conn,
		Mu:   sync.Mutex{},
	}
	s.activeClients[client] = struct{}{} // Using struct{}{} as a placeholder value

	return s
}

func (s *WsServer) RemoveConnection(conn *websocket.Conn) {
	log.Printf("Removing Connection connection")
	s.activeClientsMu.Lock()
	defer s.activeClientsMu.Unlock()

	var connsToRemove []*Client

	for client := range s.activeClients {
		if client.Conn == conn {
			connsToRemove = append(connsToRemove, client)
		}
	}

	// Use a WaitGroup to wait for all pending writes to complete
	var wg sync.WaitGroup
	for _, client := range connsToRemove {
		wg.Add(1)
		go func(client *Client) {
			client.Mu.Lock()
			defer func() {
				client.Mu.Unlock()
				wg.Done()
			}()

			// Wait for any pending writes to complete before removing the connection
			// For example, if there are multiple goroutines writing concurrently
			// to this connection, they will wait here until the write is done.
		}(client)
	}
	wg.Wait()

	// Now that all pending writes are complete, it's safe to remove the connection
	for _, client := range connsToRemove {
		delete(s.activeClients, client)
	}

	log.Printf("Removed Connection")
}
