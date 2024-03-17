package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/Tecu23/go-game/pkg/ws"
)

func (app *application) serve() error {
	server := ws.NewServer()

	http.Handle("/ws", websocket.Handler(server.HandleWS))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	srv.ListenAndServe()

	return nil
}
