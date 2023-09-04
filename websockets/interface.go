package websockets

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type SocketHandlerFunc func(w http.ResponseWriter, r *http.Request)

type Handler func(conn *websocket.Conn, payload interface{})

type WebsocketServer interface {
	New(port string) *WebsocketServer
	Start() (*WebsocketServer, error)
	Stop() error
	BlockUntilReady() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Insecure() *WebsocketServer
	EnableAll() *WebsocketServer
	UpgradeHTTPConntoWebSockets(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error)
	RemoveClient(conn *websocket.Conn)
	TerminateConnections()
}
