package websockets

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	baseRoute string
	port      string

	Upgrader websocket.Upgrader

	httpServer     *http.Server
	defaultHandler map[string]http.Handler
	WSConnHandlers map[string]Handler

	activeClientsMu sync.RWMutex
	activeClients   map[*Client]struct{}
	counter         AtomicCounter

	debug bool
}

/*
New creates the Web Socket Server Instance

# Example

Creating and Starting the Websockets Server with handlers to echo and broadcast messages with a health check.

	wsServer := server.NewWsServer("8080").EnableAll().Insecure().Start() // blocks
	defer wsServer.Stop()

Creating the Websockets Server and passing it to an HTTP Server

	wsServer := server.New("8080").EnableAll().Insecure()
	defer wsServer.Stop()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: wsServer,
	}

	_ = httpServer.ListenAndServe()


*/

func New(port string) *WsServer {
	port, defaultPath := parsePortAndPath(port)

	log.Printf("Port: %s , Path: %s", port, defaultPath)

	s := &WsServer{
		port:            port,
		baseRoute:       "/" + defaultPath,
		Upgrader:        websocket.Upgrader{},
		defaultHandler:  make(map[string]http.Handler),
		WSConnHandlers:  make(map[string]Handler),
		activeClientsMu: sync.RWMutex{},
		activeClients:   make(map[*Client]struct{}),
		counter:         AtomicCounter{},
	}

	s.defaultHandler[s.baseRoute] = http.HandlerFunc(s.RootSocketHandler)

	return s
}

/*
Entry Point for WebSockets Server.

# Incoming Requests at the root path will be handled by this function

This enables us to achieve this following pattern

	wsServer := server.New("8080").EnableAll().Insecure()

	_ := http.ListenAndServe(":8080", wsServer)

	This is possible because:

	- ServeHTTP - named ServeHTTP
	- http.Handler - Satisfies this interface



	- We can pass any Struct that has a method with the signature:
		ServeHTTP(http.ResponseWriter, *http.Request)
*/
func (s *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.defaultHandler[s.baseRoute].ServeHTTP(w, r)
}

/*
Start starts the Web Socket Server on the specified port.

# The Web Socket Server is passed to the HTTP Server as the Handler

Example

	wsServer , err := server.New("8080").EnableAll().Insecure().Start()
	defer wsServer.Stop() // Start() blocks

Optionally - to return the Web Socket Server Instance

	// Create WebSockets Config
	wsServer := server.New("8080").EnableAll().Insecure()

	// Pass to net/http
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: wsServer,
	}

	err := httpServer.ListenAndServe()
*/
func (s *WsServer) Start() (*WsServer, error) {
	s.httpServer = &http.Server{
		Addr:    ":" + s.port, // ":8080"
		Handler: s,
	}
	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting server: %v", err)
		}
	}()

	return s, nil
}

/*
BlockUntilReady - Blocks until the Web Socket Server is ready to accept connections

Example

	wsServer , err := server.NewWsServer().EnableAll().Insecure().Start()
	defer wsServer.Stop() // Start() blocks
*/
func (s *WsServer) BlockUntilReady() error {
	maxRetries := 10
	dialer := websocket.DefaultDialer

	for i := 0; i < maxRetries; i++ {
		time.Sleep(1 * time.Second) // Wait a second between attempts

		// Try to establish a WebSocket connection
		conn, _, err := dialer.Dial("ws://localhost:"+s.port+s.baseRoute, nil)
		defer s.RemoveConnection(conn)

		if err != nil {
			log.Printf("WebSocket connection attempt %d failed: %v", i+1, err)
			continue
		}

		// Send a health check message
		healthCheckReq := map[string]interface{}{
			"type": "healthcheck",
		}
		err = conn.WriteJSON(healthCheckReq)
		if err != nil {
			log.Printf("Failed to send healthcheck message: %v", err)
			conn.Close()
			continue
		}

		log.Println("Health Check Successful!")
		return nil
	}

	return fmt.Errorf("server not ready after %d attempts", maxRetries)
}

/*
Stop - Gracefully stops the Web Socket Server
*/
func (s *WsServer) Stop() error {
	log.Printf("Called Stop")
	if s.httpServer == nil {
		return errors.New("server has not been started")
	}
	s.TerminateConnections()

	s.PrintStats()

	// Create a context with a timeout to allow connections to finish
	timeoutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.httpServer.Shutdown(timeoutContext)
	return err
}

/*
Insecure
Allows the Web Socket Server to accept connections from any origin
*/
func (s *WsServer) Insecure() *WsServer {
	s.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	return s
}

/*
Enable Debugging for the Websocket Server.
*/
func (s *WsServer) Debug() *WsServer {
	s.debug = true
	return s
}

/*
EnableAll enables all the routes that are available in the Web Socket Server

Routes
  - /echo
  - /broadcast
  - /healthcheck
*/
func (s *WsServer) EnableAll() *WsServer {
	for _, endpoint := range []string{"echo", "broadcast", "healthcheck"} {
		switch endpoint {
		case "healthcheck":
			s.WSConnHandlers[endpoint] = s.HealthCheckHandler
		case "echo":
			s.WSConnHandlers[endpoint] = s.EchoHandler
		case "broadcast":
			s.WSConnHandlers[endpoint] = s.BroadcastHandler
		}
	}

	return s
}

func (s *WsServer) UpgradeHTTPConntoWebSockets(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade HTTP to WS")
		return nil, err
	}
	log.Print("Successfully Upgraded Connection")
	s.AddActiveConnection(conn)

	// s.activeClientsMu.Lock()
	// s.activeClients[conn] = true
	// s.activeClientsMu.Unlock()

	s.counter.IncrementTotalConnections()

	return conn, nil
}

func (s *WsServer) TerminateConnections() {
	log.Println("Terminating Connections")

	var connectionsToRemove []*websocket.Conn

	s.activeClientsMu.Lock()
	for client := range s.activeClients {
		connectionsToRemove = append(connectionsToRemove, client.Conn)
	}
	s.activeClientsMu.Unlock()

	for client := range s.activeClients {
		s.RemoveConnection(client.Conn)
		// err := client.Conn.Close()
		// if err != nil {
		// 	log.Printf("Error closing connection: %v", err)
		// }
	}
	log.Println("Terminated Connections")
}

func (s *WsServer) PrintStats() {
	log.Printf("Total Connections: %d", s.counter.totalConnections)
	log.Printf("Total Messages Sent: %d", s.counter.totalMessagesSent)
	log.Printf("Total Messages Received: %d", s.counter.totalMessagesReceived)
}
