package websocket

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"

	. "github.com/Tecu23/go-game/pkg/chess/constants"
	"github.com/Tecu23/go-game/pkg/chess/engine"
	"github.com/Tecu23/go-game/pkg/chess/position"
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
	// position [fen <fenstring> | startpos ]  moves <move1> .... <movei>

	fen := ""

	position.Board.NewGame()

	cmd = strings.TrimSpace(strings.TrimPrefix(cmd, "position"))

	// dividing the cmd in 2 parts, the position and the moves
	parts := strings.Split(cmd, "moves")

	if len(cmd) == 0 || len(parts) > 2 {
		Write(conn, fmt.Sprintf("error string %v wrong length=%v", parts, len(parts)))
		return
	}

	// splitting the position
	alt := strings.Split(parts[0], " ")

	alt[0] = strings.TrimSpace(alt[0])

	if alt[0] == "startpos" {
		fen = Startpos
	} else if alt[0] == "fen" {
		fen = strings.TrimSpace(strings.TrimPrefix(parts[0], "fen"))
	} else {
		Write(conn, fmt.Sprintf("%#v must be %#v or %#v", alt[0], "fen", "startpos"))
		return
	}

	// Now parsing the FEN string
	position.ParseFEN(fen)

	if len(parts) == 2 {
		parts[1] = strings.ToLower(strings.TrimSpace(parts[1]))
		position.parseMvs(parts[1])
	}
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
