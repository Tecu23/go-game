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

func InitCastlings() {
	// squares between K and R short castling
	Castl[WHITE].BetweenSh.SetBit(F1)
	Castl[WHITE].BetweenSh.SetBit(G1)
	Castl[BLACK].BetweenSh.SetBit(F8)
	Castl[BLACK].BetweenSh.SetBit(G8)

	// squares between K and R long castling
	Castl[WHITE].BetweenL.SetBit(B1)
	Castl[WHITE].BetweenL.SetBit(C1)
	Castl[WHITE].BetweenL.SetBit(D1)
	Castl[BLACK].BetweenL.SetBit(B8)
	Castl[BLACK].BetweenL.SetBit(C8)
	Castl[BLACK].BetweenL.SetBit(D8)

	// pawns stop short castling W
	Castl[WHITE].PawnsSh.SetBit(D2)
	Castl[WHITE].PawnsSh.SetBit(E2)
	Castl[WHITE].PawnsSh.SetBit(F2)
	Castl[WHITE].PawnsSh.SetBit(G2)
	Castl[WHITE].PawnsSh.SetBit(H2)
	// pawns stop long castling W
	Castl[WHITE].PawnsL.SetBit(B2)
	Castl[WHITE].PawnsL.SetBit(C2)
	Castl[WHITE].PawnsL.SetBit(D2)
	Castl[WHITE].PawnsL.SetBit(E2)
	Castl[WHITE].PawnsL.SetBit(F2)

	// pawns stop short castling B
	Castl[BLACK].PawnsSh.SetBit(D7)
	Castl[BLACK].PawnsSh.SetBit(E7)
	Castl[BLACK].PawnsSh.SetBit(F7)
	Castl[BLACK].PawnsSh.SetBit(G7)
	Castl[BLACK].PawnsSh.SetBit(H7)
	// pawns stop long castling B
	Castl[BLACK].PawnsL.SetBit(B7)
	Castl[BLACK].PawnsL.SetBit(C7)
	Castl[BLACK].PawnsL.SetBit(D7)
	Castl[BLACK].PawnsL.SetBit(E7)
	Castl[BLACK].PawnsL.SetBit(F7)

	// knights stop short castling W
	Castl[WHITE].KnightsSh.SetBit(C2)
	Castl[WHITE].KnightsSh.SetBit(D2)
	Castl[WHITE].KnightsSh.SetBit(E2)
	Castl[WHITE].KnightsSh.SetBit(G2)
	Castl[WHITE].KnightsSh.SetBit(H2)
	Castl[WHITE].KnightsSh.SetBit(D3)
	Castl[WHITE].KnightsSh.SetBit(E3)
	Castl[WHITE].KnightsSh.SetBit(F3)
	Castl[WHITE].KnightsSh.SetBit(G3)
	Castl[WHITE].KnightsSh.SetBit(H3)
	// knights stop long castling W
	Castl[WHITE].KnightsL.SetBit(A2)
	Castl[WHITE].KnightsL.SetBit(B2)
	Castl[WHITE].KnightsL.SetBit(C2)
	Castl[WHITE].KnightsL.SetBit(E2)
	Castl[WHITE].KnightsL.SetBit(F2)
	Castl[WHITE].KnightsL.SetBit(G2)
	Castl[WHITE].KnightsL.SetBit(B3)
	Castl[WHITE].KnightsL.SetBit(C3)
	Castl[WHITE].KnightsL.SetBit(D3)
	Castl[WHITE].KnightsL.SetBit(E3)
	Castl[WHITE].KnightsL.SetBit(F3)

	// knights stop short castling B
	Castl[BLACK].KnightsSh.SetBit(C7)
	Castl[BLACK].KnightsSh.SetBit(D7)
	Castl[BLACK].KnightsSh.SetBit(E7)
	Castl[BLACK].KnightsSh.SetBit(G7)
	Castl[BLACK].KnightsSh.SetBit(H7)
	Castl[BLACK].KnightsSh.SetBit(D6)
	Castl[BLACK].KnightsSh.SetBit(E6)
	Castl[BLACK].KnightsSh.SetBit(F6)
	Castl[BLACK].KnightsSh.SetBit(G6)
	Castl[BLACK].KnightsSh.SetBit(H6)
	// knights stop long castling B
	Castl[BLACK].KnightsL.SetBit(A7)
	Castl[BLACK].KnightsL.SetBit(B7)
	Castl[BLACK].KnightsL.SetBit(C7)
	Castl[BLACK].KnightsL.SetBit(E7)
	Castl[BLACK].KnightsL.SetBit(F7)
	Castl[BLACK].KnightsL.SetBit(G7)
	Castl[BLACK].KnightsL.SetBit(B6)
	Castl[BLACK].KnightsL.SetBit(C6)
	Castl[BLACK].KnightsL.SetBit(D6)
	Castl[BLACK].KnightsL.SetBit(E6)
	Castl[BLACK].KnightsL.SetBit(F6)
}
