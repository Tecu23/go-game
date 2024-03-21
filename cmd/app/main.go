package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/Tecu23/go-game/pkg/chess/castlings"
	. "github.com/Tecu23/go-game/pkg/chess/constants"
	"github.com/Tecu23/go-game/pkg/chess/magic"
	"github.com/Tecu23/go-game/pkg/chess/position"
	"github.com/Tecu23/go-game/pkg/websocket"
)

func main() {
	addr := ":3000"

	server := websocket.NewServer(&addr)
	log.Info("starting websocket server")
	Init()
	server.Start()
}

func Init() {
	InitFen2Sq()
	magic.InitMagic()
	position.InitKeys()
	position.InitAtksKings()
	position.InitAtksKnights()
	castlings.InitCastlings()
	position.PcSqInit()
	position.Board.NewGame()

	// run setoption name hash value 32
}
