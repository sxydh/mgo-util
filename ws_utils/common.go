package ws_utils

import (
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

type WsServer struct {
	OnConn    func(conn *websocket.Conn)
	OnMessage func(msg string)
}

func (server *WsServer) RandPort() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		port := 40000 + r.Intn(10000)
		listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			continue
		}
		_ = listener.Close()
		go func() {
			server.Port(port)
		}()
		return port
	}
}

func (server *WsServer) Port(port int) {
	http.HandleFunc("/", server.wsHandler)
	log.Printf("ListenAndServe going: port=%v", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Printf("ListenAndServe error: err=%v", err)
	}
}

//goland:noinspection GoUnhandledErrorResult
func (server *WsServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: err=%v", err)
		return
	}
	defer wsConn.Close()
	log.Printf("Get user connection: remoteAddr=%v", r.RemoteAddr)
	server.OnConn(wsConn)
	for {
		_, p, err := wsConn.ReadMessage()
		if err != nil {
			log.Printf("ReadMessage error: err=%v", err)
			return
		}
		server.OnMessage(string(p))
	}
}

func (server *WsServer) Send(conn *websocket.Conn, msg string) error {
	return conn.WriteMessage(websocket.TextMessage, []byte(msg))
}
