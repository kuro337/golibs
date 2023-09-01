package websockets

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// "broadcast" Route Handler - sends msgs to all Clients
func (s *WsServer) BroadcastHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.UpgradeHTTPtoWS(w, r)
	if err != nil {
		log.Println("Failed to Upgrade HTTP Connection to WebSocket")
		return
	}
	defer s.RemoveClient(conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		atomic.AddInt64(&s.totalMessagesReceived, 1)

		// Print the received message
		fmt.Printf("Received: %s\n", p)

		// Broadcast the message to all clients
		s.BroadcastMessage(messageType, p) // Just call the function without assigning the result
	}
}

// "echo" route - sends msgs back to the Client
func (s *WsServer) EchoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.UpgradeHTTPtoWS(w, r)
	if err != nil {
		log.Println("Failed to Upgrade HTTP Connection to WebSocket")
		return
	}
	defer s.RemoveClient(conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		atomic.AddInt64(&s.totalMessagesReceived, 1)

		// Print the received message
		fmt.Printf("Received: %s\n", p)

		// Echo the message back to the sender
		if err := conn.WriteMessage(messageType, p); err != nil {
			fmt.Println(err)
			return
		}
		atomic.AddInt64(&s.totalMessagesSent, 1)

	}
}

// Health Check Route
func (s *WsServer) EnableHealthCheck() *WsServer {
	s.routeHandlers["/health"] = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
	return s
}

/*

HELPERS

1. BroadcastMessage - sends a message to all clients
2. RemoveClient - removes a client from the activeClients map
3. UpgradeHTTPtoWS - upgrades an HTTP connection to a WebSocket connection
4. ServeHTTP - handles incoming requests at the root path
5. EnableAll - enables all routes
6. Insecure - allows all origins
7. Echo - adds an echo route
8. Broadcast - adds a broadcast route


*/

// Broadcasts message to all clients - used by "broadcast" route
func (s *WsServer) BroadcastMessage(messageType int, data []byte) {
	s.activeClientsMu.RLock()
	for conn := range s.activeClients {
		go func(c *websocket.Conn) {
			if err := c.WriteMessage(messageType, data); err != nil {
				fmt.Println("Error sending message:", err)
			} else {
				atomic.AddInt64(&s.totalMessagesSent, 1)
			}
		}(conn)
	}
	s.activeClientsMu.RUnlock()
}

// Sends a message "Hello, everyone!" to all clients
func (s *WsServer) GreetingEchoAllClients(w http.ResponseWriter, r *http.Request) {
	message := []byte("Hello, everyone!")
	messageType := websocket.TextMessage
	s.BroadcastMessage(messageType, message) // Just call the function without assigning the result
}
