package websocket

import "github.com/gorilla/websocket"

func handleBestMove(conn *websocket.Conn, bestMove string) {
	Write(conn, "info string cmd bm not implement yet")
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

func handleGo(conn *websocket.Conn, toEng chan bool, words []string) {
	Write(conn, "info string cmd go not implement yet")
}

func handlePonderhit(conn *websocket.Conn) {
	Write(conn, "info string cmd ponderhit not implement yet")
}

func handleStop(conn *websocket.Conn) {
	Write(conn, "info string cmd stop not implement yet")
}

func handleQuit(conn *websocket.Conn) {
	Write(conn, "info string cmd quit not implement yet")
}
