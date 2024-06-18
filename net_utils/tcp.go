package net_utils

import (
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type Server struct {
	onConn    func(conn *net.Conn)
	onMessage func(msg string)
}

func (server *Server) TcpServer() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		port := 40000 + r.Intn(10000)
		err := server.TcpServerOnPort(port)
		if err == nil {
			return port
		}
	}
}

//goland:noinspection GoUnhandledErrorResult
func (server *Server) TcpServerOnPort(port int) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Printf("Listen on port error, port=%v, err=%v", port, err)
		return err
	}
	defer listener.Close()
	log.Printf("Listening for tcp..., port=%v", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept on port error, port=%v, err=%v", port, err)
			return err
		}
		server.onConn(&conn)
		log.Printf("Accepting..., localAddr=%v， remoteAddr=%v", conn.LocalAddr(), conn.RemoteAddr())

		go func() {
			defer conn.Close()
			bytes := make([]byte, 1024)
			_, err = conn.Read(bytes)
			if err != nil {
				log.Printf("Read error, localAddr=%v， remoteAddr=%v", conn.LocalAddr(), conn.RemoteAddr())
				return
			}
			server.onMessage(string(bytes))
		}()
	}
}
