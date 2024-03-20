package engine

import (
	. "github.com/Tecu23/go-game/pkg/chess/constants"
	"github.com/Tecu23/go-game/pkg/chess/moves"
)

// ///////////////  Killers ///////////////////////////////////////////////
// killerStruct holds the killer moves per ply
type KillerStruct [MaxPly]struct {
	K1 moves.Move
	K2 moves.Move
}

// Clear killer moves
func (k *KillerStruct) Clear() {
	for ply := 0; ply < MaxPly; ply++ {
		k[ply].K1 = moves.NoMove
		k[ply].K2 = moves.NoMove
	}
}

// add killer 1 and 2 (Not inCheck, caaptures and promotions)
func (k *KillerStruct) Add(mv moves.Move, ply int) {
	if !k[ply].K1.CmpFrTo(mv) {
		k[ply].K2 = k[ply].K1
		k[ply].K1 = mv
	}
}

var Killers KillerStruct
