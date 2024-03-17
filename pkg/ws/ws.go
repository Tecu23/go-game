package ws

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

// Server is used to maintain connections by using a map
type Server struct {
	conns map[*websocket.Conn]bool
}

// NewServer creates a new server with no connections
func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) HandleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())

	// TODO: use mutex to protect agains race conditions
	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			// If Connection has ended from the client, we need to exit
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			// Handle wrong input by the client to the connection
			continue
		}

		msg := buf[:n]

		s.broadcast(msg)
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error:", err)
			}
		}(ws)
	}
}
