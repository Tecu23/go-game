package moves

import (
	. "github.com/Tecu23/go-game/pkg/chess/constants"
)

// Pc2Fen convert pc to fenString
func Pc2Fen(pc int) string {
	if pc == Empty {
		return " "
	}
	return string(PcFen[pc])
}
