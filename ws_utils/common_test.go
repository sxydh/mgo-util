package ws_utils

import (
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
			server := &WsServer{}
			got := server.RandPort()
			log.Printf("RandPort: got=%v", got)
			select {}
		})
	}
}
