package castlings

import (
	"strings"

	"github.com/Tecu23/go-game/pkg/chess/bitboard"
	. "github.com/Tecu23/go-game/pkg/chess/constants"
)

type Castlings uint

const (
	ShortW = uint(0x1) // white can castle short
	LongW  = uint(0x2) // white can castle long
	ShortB = uint(0x4) // black can castle short
	LongB  = uint(0x8) // black can castle short
)

type castlOptions struct {
	Short                                uint // flag
	Long                                 uint // flag
	Rook                                 int  // rook pc (wR/bR)
	KingPos                              int  // king pos
	RookSh                               uint // rook pos short
	RookL                                uint // rook pos long
	BetweenSh                            bitboard.BitBoard
	BetweenL                             bitboard.BitBoard
	PawnsSh, PawnsL, KnightsSh, KnightsL bitboard.BitBoard
}

var Castl = [2]castlOptions{
	{ShortW, LongW, WR, E1, H1, A1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	{ShortB, LongB, BR, E8, H8, A8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
}

func (c Castlings) Flags(sd Color) bool {
	return c.ShortFlag(sd) || c.LongFlag(sd)
}

func (c Castlings) ShortFlag(sd Color) bool {
	return (Castl[sd].Short & uint(c)) != 0
}

func (c Castlings) LongFlag(sd Color) bool {
	return (Castl[sd].Long & uint(c)) != 0
}

func (c *Castlings) On(val uint) {
	(*c) |= Castlings(val)
}

func (c *Castlings) Off(val uint) {
	(*c) &= Castlings(^val)
}

// ParseCastlings should parse the castlings part of a fen string
func ParseCastlings(fenCastl string) Castlings {
	c := uint(0)

	if fenCastl == "-" {
		return Castlings(0)
	}

	if strings.Index(fenCastl, "K") >= 0 {
		c |= ShortW
	}
	if strings.Index(fenCastl, "Q") >= 0 {
		c |= LongW
	}
	if strings.Index(fenCastl, "k") >= 0 {
		c |= ShortB
	}
	if strings.Index(fenCastl, "q") >= 0 {
		c |= LongB
	}

	return Castlings(c)
}

func (c Castlings) String() string {
	flags := ""
	if uint(c)&ShortW != 0 {
		flags = "K"
	}
	if uint(c)&LongW != 0 {
		flags += "Q"
	}
	if uint(c)&ShortB != 0 {
		flags += "k"
	}
	if uint(c)&LongB != 0 {
		flags += "q"
	}
	if flags == "" {
		flags = "-"
	}
	return flags
}
