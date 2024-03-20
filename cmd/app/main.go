package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/Tecu23/go-game/pkg/websocket"
)

func main() {
	addr := ":3000"

	server := websocket.NewServer(&addr)
	log.Info("starting websocket server")
	server.Start()
}
