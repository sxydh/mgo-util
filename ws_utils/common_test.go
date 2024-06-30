package ws_utils

import (
	"github.com/gorilla/websocket"
	"log"
	"testing"
)

func TestWsServer_RandPort(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Normal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &WsServer{
				OnConn: func(conn *websocket.Conn) {
				},
				OnMessage: func(msg string) {
					log.Printf("OnMessage: msg=%v", msg)
				},
			}
			got := server.RandPort("/")
			log.Printf("RandPort: got=%v", got)
			select {}
		})
	}
}
