// Package websocket secures the connection between the client and the engine
package websocket

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/Tecu23/go-game/pkg/chess/engine"
)

// Uci should connect the websocket to the actual uci interface
func (srv *Server) Uci(
	w http.ResponseWriter,
	r *http.Request,
	conn *websocket.Conn,
) (chan string, *websocket.Conn) {
	log.Info("websocket connection established")
	Write(conn, "info string starting Engine")

	message := make(chan string)

	go func() {
		defer Write(conn, "info string quit Engine")
		defer conn.Close()
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Error("here", err)
				return
			}

			message <- string(p)
		}
	}()

	return message, conn
}

func uci(input chan string, conn *websocket.Conn) {
	toEng, frEng := engine.Engine()
	var cmd string
	var bestMove string
	quit := false

	for !quit {
		select {
		case cmd = <-input:
			log.Info(cmd)
		case bestMove = <-frEng:
			handleBestMove(conn, bestMove)
			continue
		}

		words := strings.Split(cmd, " ")
		words[0] = strings.TrimSpace(strings.ToLower(words[0]))

		switch words[0] {
		case "uci":
			handleUci(conn)
		case "setoption":
			handleSetOption(conn, words)
		case "isready":
			handleIsReady(conn)
		case "ucinewgame":
			handleNewgame(conn) // add trans.clear() to handle new game method
		case "position":
			handlePosition(conn, cmd)
		case "debug":
			handleDebug(conn, words)
		case "register":
			handleRegister(conn, words)
		case "go":
			handleGo(conn, toEng, words)
		case "ponderhit":
			handlePonderhit(conn)
		case "stop":
			handleStop(conn)
		case "quit", "q":
			handleQuit(conn)
			quit = true
			continue

			// CUSTOM COMMANDS TO HELP WITH DEBUG AND TESTING
		case "perft":
			handlePerformanveTest(conn, words)
		case "pb": // Print current board
			handlePrintBoard(conn)
		case "pbb": // Print all bitboard
			handlePrintAllBitBoards(conn)
		case "pm": // Print all legal moves
			handlePrintAllLegalMoves(conn)
		case "eval": // Evaluate current position
			handleEvaluatePosition(conn)
		case "pos":
			handleMyPositions(conn, words)
		case "moves":
			handleMyMoves(conn, words)
		case "key":
			handleKey(conn)
		case "see":
			handleSee(conn, words)
		case "qs":
			handleQs(conn)
		case "hist":
			handleHistory(conn)
		case "moveval":
			handleMoveValue(conn)
		default:
			Write(conn, fmt.Sprintf("info string unknown cmd %s", cmd))
		}
	}
	Write(conn, "info string leaving uci()")
}
