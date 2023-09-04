package websockets

import (
	"log"
	"reflect"
	"strings"

	"github.com/gorilla/websocket"
)

// parsePortAndPath parses the port and path from the given input string.
func parsePortAndPath(input string) (string, string) {
	parts := strings.SplitN(input, "/", 2)
	port := parts[0]
	path := "" // Default path is "/ws"
	if len(parts) > 1 && parts[1] != "" {
		path = parts[1]
	}
	return port, path
}

func handleSocketError(err error, expectedErrors ...interface{}) (error, bool) {
	if err == nil {
		return err, false
	}
	errType := reflect.TypeOf(err)
	for _, expectedErr := range expectedErrors {
		if errType == reflect.TypeOf(expectedErr) {
			return err, true
		}
	}

	log.Print(errType.String())

	if websocket.IsUnexpectedCloseError(err,
		websocket.CloseGoingAway,
		websocket.CloseNormalClosure,
		websocket.CloseNoStatusReceived) {

		switch {
		case err == websocket.ErrCloseSent:
			log.Println("Close message sent by peer")
		case err == websocket.ErrReadLimit:
			log.Println("Read limit exceeded")
		default:
			log.Printf("Unexpected close error: %v", err)
		}
	} else {
		log.Printf("Error reading message: %v", err)
	}

	return err, false
}
