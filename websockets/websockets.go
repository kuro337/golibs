package websockets

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type WsServer struct {
	Upgrader              websocket.Upgrader
	activeClients         map[*websocket.Conn]bool
	activeClientsMu       sync.RWMutex
	routeHandlers         map[string]func(w http.ResponseWriter, r *http.Request)
	totalConnections      int64
	totalMessagesSent     int64
	totalMessagesReceived int64
	port                  string
	httpServer            *http.Server
}

/*
NewWsServer creates the Web Socket Server Instance

# Example

Creating the Websockets Server and passing it to an HTTP Server

	wsServer := server.NewWsServer("8080").EnableAll().Insecure().Start() // blocks
	defer wsServer.Stop()





*/

func NewWsServer(port string) *WsServer {
	return &WsServer{
		Upgrader:      websocket.Upgrader{},
		activeClients: make(map[*websocket.Conn]bool),
		routeHandlers: make(map[string]func(w http.ResponseWriter, r *http.Request)),
		port:          port,
	}
}

/*
Start starts the Web Socket Server on the specified port.

# The Web Socket Server is passed to the HTTP Server as the Handler

Example

	wsServer , err := server.NewWsServer().EnableAll().Insecure().Start()
	defer wsServer.Stop() // Start() blocks

Optionally - to return the Web Socket Server Instance

	// Create WebSockets Config
	wsServer := server.NewWsServer().EnableAll().Insecure()

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
	for i := 0; i < maxRetries; i++ {
		time.Sleep(1 * time.Second) // Wait a second between attempts

		resp, err := http.Get("http://localhost:" + s.port + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			return nil
		}
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

	fmt.Println("Terminated Connections")
	log.Printf("Total Connections: %d", s.totalConnections)
	log.Printf("Total Messages Sent: %d", s.totalMessagesSent)
	log.Printf("Total Messages Received: %d", s.totalMessagesReceived)

	// Create a context with a timeout to allow connections to finish
	timeoutContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.httpServer.Shutdown(timeoutContext)
	return err
}

/*
Entry Point for WebSockets Server.

# Incoming Requests at the root path will be handled by this function

This enables us to achieve this following pattern

	wsServer := server.NewWsServer().EnableAll().Insecure()

	_ := http.ListenAndServe(":8080", wsServer)

	This is possible because:

	- ServeHTTP - named ServeHTTP
	- http.Handler - Satisfies this interface



	- We can pass any Struct that has a method with the signature:
		ServeHTTP(http.ResponseWriter, *http.Request)
*/
func (s *WsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Map incoming Request to the correct Path.
	if handler, exists := s.routeHandlers[r.URL.Path]; exists {
		handler(w, r)
	} else {
		http.NotFound(w, r)
	}
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
EnableAll enables all the routes that are available in the Web Socket Server

Routes
  - /echo
  - /broadcast
*/
func (s *WsServer) EnableAll() *WsServer {
	s.routeHandlers["/echo"] = s.EchoHandler
	s.routeHandlers["/broadcast"] = s.BroadcastHandler
	return s
}

func (s *WsServer) Echo(route string) *WsServer {
	s.routeHandlers[route] = s.EchoHandler
	return s
}

func (s *WsServer) Broadcast(route string) *WsServer {
	s.routeHandlers[route] = s.GreetingEchoAllClients
	return s
}

func (s *WsServer) UpgradeHTTPtoWS(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade HTTP to WS")
		return nil, err
	}
	log.Print("Successfully Upgraded Connection")
	s.activeClientsMu.Lock()
	s.activeClients[conn] = true
	s.activeClientsMu.Unlock()

	atomic.AddInt64(&s.totalConnections, 1)

	return conn, nil
}

func (s *WsServer) RemoveClient(conn *websocket.Conn) {
	s.activeClientsMu.Lock()
	delete(s.activeClients, conn)
	s.activeClientsMu.Unlock()
}

func (s *WsServer) TerminateConnections() {
	for conn := range s.activeClients {
		conn.Close()
	}
}
