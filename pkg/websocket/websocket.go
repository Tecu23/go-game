package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// Server should ...
type Server struct {
	upgrader websocket.Upgrader
	addr     *string
}

// NewServer should create a new Server Object
func NewServer(addr *string) *Server {
	srv := &Server{
		addr:     addr,
		upgrader: websocket.Upgrader{},
	}

	http.HandleFunc("/uci", srv.uciHandler)

	return srv
}

func (srv *Server) uciHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("upgrading to websocket connection")

	srv.upgrader.CheckOrigin = func(_ *http.Request) bool { return true }

	conn, err := srv.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade:", err)
		return
	}

	uci(srv.Uci(w, r, conn))
}

// Start should start the web server
func (srv *Server) Start() {
	log.Info(*srv.addr)
	log.Info("starting websocket server")
	err := http.ListenAndServe(*srv.addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
