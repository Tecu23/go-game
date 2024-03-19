package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"

	"github.com/Tecu23/go-game/pkg/chess/engine"
)

var savedBestMove = ""

func handleBestMove(conn *websocket.Conn, bestMove string) {
	if engine.Limits.Infinite {
		savedBestMove = bestMove
		return
	}

	Write(conn, fmt.Sprintf("info string bestMove: %s", savedBestMove))
}

func handleUci(conn *websocket.Conn) {
	Write(conn, "id name GoEng")
	Write(conn, "id author Tecu23")

	Write(conn, "uciok")
}

func handleSetOption(conn *websocket.Conn, words []string) {
	Write(conn, "info string cmd setoption not implement yet")
}

func handleIsReady(conn *websocket.Conn) {
	Write(conn, "readyok")
}

func handleNewgame(conn *websocket.Conn) {
	Write(conn, "info string cmd ucinewgame not implement yet")
}

func handlePosition(conn *websocket.Conn, cmd string) {
	Write(conn, "info string cmd position not implement yet")
}

func handleDebug(conn *websocket.Conn, words []string) {
	Write(conn, "info string cmd debug not implement yet")
}

func handleRegister(conn *websocket.Conn, words []string) {
	Write(conn, "info string cmd register not implement yet")
}

func handleGo(conn *websocket.Conn, toEng chan string, words []string) {
	Write(conn, "info string cmd go not implement yet")
}

func handlePonderhit(conn *websocket.Conn) {
	Write(conn, "info string cmd ponderhit not implement yet")
}

func handleStop(conn *websocket.Conn) {
	if engine.Limits.Infinite {
		if savedBestMove != "" {
			Write(conn, fmt.Sprintf("info string bestMove: %s", savedBestMove))
			savedBestMove = ""
		}

		engine.Limits.SetInfinite(false)
	}

	engine.Limits.SetStop(true)
}

// Possibly not necessary, the client will close the connection by just closing the socket
func handleQuit(conn *websocket.Conn) {
	Write(conn, "info string cmd quit not implement yet")
}
