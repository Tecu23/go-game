package websocket

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"

	. "github.com/Tecu23/go-game/pkg/chess/constants"
	"github.com/Tecu23/go-game/pkg/chess/engine"
	"github.com/Tecu23/go-game/pkg/chess/moves"
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
	if len(words) < 5 {
		Write(
			conn,
			fmt.Sprintf("info string don't have this option %s", strings.Join(words[:], " ")),
		)
	}

	if strings.ToLower(strings.TrimSpace(words[1])) != "name" {
		Write(
			conn,
			fmt.Sprintf(
				"info string 'name' is missing in this option %s",
				strings.Join(words[:], " "),
			),
		)
	}

	switch strings.ToLower(strings.TrimSpace(words[2])) {
	case "hash":
		if strings.TrimSpace(strings.ToLower(words[3])) != "value" {
			Write(
				conn,
				fmt.Sprintf(
					"info string 'value' is missing in this option %s",
					strings.Join(words[:], " "),
				),
			)
		}

		if val, err := strconv.Atoi(strings.TrimSpace(words[4])); err == nil {
			if err = position.Trans.New(val); err != nil {
				Write(
					conn,
					fmt.Sprintf(
						"info string %s ",
						err.Error(),
					),
				)
			}
		} else {
			Write(
				conn,
				fmt.Sprintf(
					"info string the Hash value is not numeric %s",
					strings.Join(words[:], " "),
				),
			)
		}
	default:
		Write(
			conn,
			fmt.Sprintf(
				"info string don't have this option %s",
				strings.Join(words[:], " "),
			),
		)

	}
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
		position.ParseMvs(parts[1])
	}
}

func handleDebug(conn *websocket.Conn, words []string) {
	Write(conn, "info string cmd debug not implement yet")
}

func handleRegister(conn *websocket.Conn, words []string) {
	Write(conn, "info string cmd register not implement yet")
}

func handleGo(conn *websocket.Conn, toEng chan bool, words []string) {
	// TODO: Right now can only handle one of them at a time. We need to be able to mix them
	engine.Limits.Init()
	if len(words) > 1 {
		words[1] = strings.TrimSpace(strings.ToLower(words[1]))
		switch words[1] {
		case "searchmoves":
			Write(conn, "info string go searchmoves not implemented")
		case "ponder":
			Write(conn, "info string go ponder not implemented")
		case "wtime":
			Write(conn, "info string go wtime not implemented")
		case "btime":
			Write(conn, "info string go btime not implemented")
		case "winc":
			Write(conn, "info string go winc not implemented")
		case "binc":
			Write(conn, "info string go binc not implemented")
		case "movestogo":
			Write(conn, "info string go movestogo not implemented")
		case "depth":
			d := -1
			err := error(nil)
			if len(words) >= 3 {
				d, err = strconv.Atoi(words[2])
			}
			if d < 0 || err != nil {
				Write(conn, "info string depth not numeric")
				return
			}
			engine.Limits.SetDepth(d)
			toEng <- true
		case "nodes":
			Write(conn, "info string go nodes not implemented")
		case "movetime":
			mt, err := strconv.Atoi(words[2])
			if err != nil {
				Write(conn, fmt.Sprintf("info string %s not numeric", words[2]))
				return
			}
			engine.Limits.SetMoveTime(mt)
			toEng <- true
		case "mate": // mate <x>  mate in x moves
			Write(conn, "info string go mate not implemented")
		case "infinite":
			engine.Limits.SetInfinite(true)
			toEng <- true
		case "register":
			Write(conn, "info string go register not implemented")
		default:
			Write(conn, fmt.Sprintf("info string go %s not implemented", words[1]))
		}
	} else {
		Write(conn, "info string suppose go infinite")
		engine.Limits.SetInfinite(true)
		toEng <- true
	}
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

func handlePerformanveTest(conn *websocket.Conn, words []string) {
	if len(words) > 1 {
		depth, err := strconv.Atoi(words[1])
		if err != nil {
			Write(conn, err.Error())
		} else {
			engine.StartPerft(depth, &position.Board)
		}
	}
}

func handlePrintBoard(conn *websocket.Conn) {
	position.Board.Print()
}

func handlePrintAllBitBoards(conn *websocket.Conn) {
	position.Board.PrintAllBB()
}

func handlePrintAllLegalMoves(conn *websocket.Conn) {
	position.Board.PrintAllLegals()
}

func handleEvaluatePosition(conn *websocket.Conn) {
	Write(conn, fmt.Sprintf("eval = %s", position.Evaluate(&position.Board)))
}

func handleMyPositions(conn *websocket.Conn, words []string) {
	if len(words) < 2 {
		Write(
			conn,
			fmt.Sprintf("info string not correct pos command %s", strings.Join(words[:], " ")),
		)
	}

	words[1] = strings.TrimSpace(strings.ToLower(words[1]))
	handleSetOption(conn, strings.Split("setoption name hash value 256", " "))

	switch words[1] {
	case "london": // London position
		handlePosition(
			conn,
			"position startpos moves d2d4 d7d5 c1f4 g8f6 e2e3 c7c5 b1d2 b8c6 c2c3 e7e6 f1d3 f8d6",
		)
	case "phil": // Philidor position
		handlePosition(
			conn,
			"position startpos moves e2e4 d7d6 d2d4 e7e5 d4e5 d6e5 d1d8 e8d8 g1f3 f7f6 b1c3 c7c6 f1c4",
		)
	case "english": // English position
		handlePosition(
			conn,
			"position startpos moves c2c4 e7e5 g2g3 b8c6 f1g2 g7g6 b1c3 f8g7 e2e4 d7d6 g1e2 g8f6",
		)
	case "bogo": // Bogo Indian position
		handlePosition(
			conn,
			"position fen 1rb1r1k1/2pn1ppp/1p1pqn2/p4N2/2PPp1P1/2P1B2P/P1Q1PPB1/1R3RK1 b - - 0 16",
		)
	default:
		Write(
			conn,
			fmt.Sprintf(
				"info string not correct pos command %s doesn't exist. %s",
				words[1],
				strings.Join(words[:], " "),
			),
		)
	}

	engine.History.Clear()
	position.Trans.Clear()
}

func handleMyMoves(conn *websocket.Conn, words []string) {
	mvString := strings.Join(words[1:], " ")
	position.ParseMvs(mvString)
}

func handleKey(conn *websocket.Conn) {
	Write(conn, fmt.Sprintf("key = %x, fullkey=%x\n", position.Board.Key, position.Board.FullKey()))
	index := position.Board.FullKey() & uint64(position.Trans.Mask)
	lock := position.Trans.Lock(position.Board.FullKey())
	Write(conn, fmt.Sprintf("index = %x, lock=%x\n", index, lock))
}

func handleSee(conn *websocket.Conn, words []string) {
	fr, to := Empty, Empty
	if len(words[1]) == 2 && len(words[2]) == 2 {
		fr = Fen2Sq[words[1]]
		to = Fen2Sq[words[2]]
	} else if len(words[1]) == 4 {
		fr = Fen2Sq[words[1][0:2]]
		to = Fen2Sq[words[1][2:]]
	} else {
		fmt.Println("error in fr/to")
	}

	Write(conn, fmt.Sprintln("see = ", engine.See(fr, to, &position.Board)))
}

func handleQs(conn *websocket.Conn) {
	Write(conn, fmt.Sprintln("qs =", engine.Qs(MaxEval, &position.Board)))
}

func handleHistory(conn *websocket.Conn) {
	engine.History.Print(50)
}

func handleMoveValue(conn *websocket.Conn) {
	// print all legal moves with different values
	b := &position.Board
	transMove := moves.NoMove
	transDepth := 4
	ply := 1

	var transSc, scType int
	ok := false

	transMove, transSc, scType, ok = position.Trans.Retrieve(b.FullKey(), transDepth, ply)
	_, _, _ = ok, transSc, scType

	var childPV engine.PvList
	childPV.New() // TODO? make it smaller for each depth maxDepth-ply
	// bs, score := noScore, noScore
	// bm := noMove

	genInfo := engine.GenInfoStruct{Sv: 0, Ply: 1, TransMove: transMove}
	engine.Next = engine.NextNormal
	ix := 0
	bestSc, bestMv, bestHsc, bestHmv := MinEval, moves.NoMove, MinEval, moves.NoMove
	bestHmsg, bestMsg := "", ""
	for mv, msg := engine.Next(&genInfo, b); mv != moves.NoMove; mv, msg = engine.Next(&genInfo, b) {
		if !b.Move(mv) {
			continue
		}
		b.Unmove(mv)
		seeVal := engine.See(mv.Fr(), mv.To(), &position.Board)
		sc := int(engine.History.Get(mv.Fr(), mv.To(), position.Board.Stm))
		if sc > bestHsc {
			bestHsc = sc
			bestHmv = mv
			bestHmsg = msg
		}
		if seeVal < 0 {
			sc = seeVal
		}

		if sc > bestSc {
			bestSc = sc
			bestMv = mv
			bestMsg = msg
		}
		fmt.Printf(
			"%v: %v history %v,\tsee %3v,\tdpcSqTab %v\t(%v)\n",
			ix+1,
			mv,
			engine.History.Get(mv.Fr(), mv.To(), position.Board.Stm),
			engine.See(mv.Fr(), mv.To(), &position.Board),
			position.PcSqScore(mv.Pc(), mv.To())-position.PcSqScore(mv.Pc(), mv.Fr()),
			msg,
		)
		ix++
	}

	fmt.Printf(
		"best History (%v): %v %v    \tbest hist+see (%v): %v %v  \n",
		bestHmsg,
		bestHmv,
		bestHsc,
		bestMsg,
		bestMv,
		bestSc,
	)
}
