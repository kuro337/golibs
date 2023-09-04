package websockets

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

// All client requests are sent to this Handler - then this distributes the request to the appropriate handler
func (s *WsServer) RootSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.UpgradeHTTPConntoWebSockets(w, r)
	if err != nil {
		log.Printf("Failed to upgrade: %v", err)
		return
	}
	defer s.RemoveConnection(conn)

	// Keep Reading Messages from the Client connection
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil { // If Error reading message - handle Error
			if _, expected := handleSocketError(err, &websocket.CloseError{}, &net.OpError{}); expected {
				log.Printf("Successfully Closed the Connection")
				break
			}
		}

		log.Printf("+1 Sent.")
		s.counter.IncrementTotalMessagesReceived()
		var messageData map[string]interface{}
		if err := json.Unmarshal(msg, &messageData); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		// Map message type to appropriate Handler
		if handlerFunc, exists := s.WSConnHandlers[messageData["type"].(string)]; exists {
			handlerFunc(conn, messageData["payload"])
		} else {
			log.Printf("Unsupported message type: %s", messageData["type"])
		}
	}
}

func (s *WsServer) EchoHandler(conn *websocket.Conn, payload interface{}) {
	// Placeholder for the echo functionality
	message, ok := payload.(string)
	if !ok {
		log.Println("Failed to cast payload to string")
		return
	}
	// Convert the message into JSON format
	response := map[string]interface{}{
		"type":    "echo",
		"payload": message,
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, responseData)
	if err != nil {
		log.Printf("Failed to send echo message: %v", err)
	}
	s.counter.IncrementTotalMessagesSent()
}

func (s *WsServer) BroadcastHandler(conn *websocket.Conn, payload interface{}) {
	message, ok := payload.(string)
	if !ok {
		log.Println("Failed to cast payload to string")
		return
	}

	s.activeClientsMu.RLock() // Ensure to use the Read Lock since you're only reading from the map
	defer s.activeClientsMu.RUnlock()

	for client := range s.activeClients {
		client.Mu.Lock()
		err := client.Conn.WriteJSON(map[string]interface{}{
			"type":    "broadcast",
			"payload": message,
		})
		if err != nil {
			log.Printf("Failed to broadcast message to a client: %v", err)
		}
		client.Mu.Unlock()

		s.counter.IncrementTotalMessagesSent()
	}
}

func (s *WsServer) HealthCheckHandler(conn *websocket.Conn, payload interface{}) {
	response := map[string]interface{}{
		"type":    "healthcheck",
		"payload": "Server is running",
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		log.Printf("Failed to marshal healthcheck response: %v", err)
		return
	}
	err = conn.WriteMessage(websocket.TextMessage, responseData)
	if err != nil {
		log.Printf("Failed to send healthcheck response: %v", err)
	}
	s.counter.IncrementTotalMessagesSent()
}
