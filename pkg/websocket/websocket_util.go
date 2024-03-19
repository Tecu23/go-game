package websocket

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func Write(conn *websocket.Conn, msg string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Fatal("write:", err)
	}
}
