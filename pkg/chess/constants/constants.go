// Package constants contains the varibles used throughout the project
package constants

import (
	"github.com/Tecu23/go-game/pkg/chess/bitboard"
)

// Evaluate constants
const (
	MaxEval  = +10000
	MinEval  = -MaxEval
	MateEval = MaxEval + 1
	NoScore  = MinEval - 1
)

var PieceVal = [16]int{
	100,
	-100,
	325,
	-325,
	350,
	-350,
	500,
	-500,
	950,
	-950,
	10000,
	-10000,
	0,
	0,
	0,
	0,
}

var (
	KnightFile = [8]int{-4, -3, -2, +2, +2, 0, -2, -4}
	KnightRank = [8]int{-15, 0, +5, +6, +7, +8, +2, -4}
	CenterFile = [8]int{-8, -1, 0, +1, +1, 0, -1, -3}
	KingFile   = [8]int{+1, +2, 0, -2, -2, 0, +2, +1}
	KingRank   = [8]int{+1, 0, -2, -4, -6, -8, -10, -12}
	PawnRank   = [8]int{0, 0, 0, 0, +2, +6, +25, 0}
	PawnFile   = [8]int{0, 0, +1, +10, +10, +8, +10, +8}
)

const LongDiag = 10

// Piece Square Table
var PSqTab [12][64]int

// Engine Constants
const (
	MaxDepth = 100
	MaxPly   = 100
)

const (
	NoPiecesC = 12       // NoPiecesC is the number of total pieces
	NoPiecesT = 6        // NoPiecesT is the number of pieces each side has
	WHITE     = Color(0) // WHITE is the color white
	BLACK     = Color(1) // BLACK is the color black

	Startpos = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - " // Startpos is Game Starting position

	Row1  = bitboard.BitBoard(0x00000000000000FF)
	Row2  = bitboard.BitBoard(0x000000000000FF00)
	Row3  = bitboard.BitBoard(0x0000000000FF0000)
	Row4  = bitboard.BitBoard(0x00000000FF000000)
	Row5  = bitboard.BitBoard(0x000000FF00000000)
	Row6  = bitboard.BitBoard(0x0000FF0000000000)
	Row7  = bitboard.BitBoard(0x00FF000000000000)
	Row8  = bitboard.BitBoard(0xFF00000000000000)
	FileA = bitboard.BitBoard(0x0101010101010101)
	FileB = bitboard.BitBoard(0x0202020202020202)
	FileG = bitboard.BitBoard(0x4040404040404040)
	FileH = bitboard.BitBoard(0x8080808080808080)
)

// piece char definitions
const (
	PcFen = "PpNnBbRrQqKk     "
	PtFen = "PNBRQK?"
)

// Fen2Sq maps fen-sq to int
var Fen2Sq = make(map[string]int)

// Sq2Fen maps int-sq to fen
var Sq2Fen = make(map[int]string)

// Each piece has a particular nummber for easy accessing
const (
	Pawn int = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

// 12 pieces with color plus empty
const (
	WP = iota
	BP
	WN
	BN
	WB
	BB
	WR
	BR
	WQ
	BQ
	WK
	BK
	Empty = 15
)

const (
	E  = +1
	W  = -1
	N  = 8
	S  = -8
	NW = +7
	NE = +9
	SW = -NE
	SE = -NW

	FrMask     = 0x0000003f                 // 0000 0000  0000 0000  0000 0000  0011 1111
	ToMask     = 0x00000fd0                 // 0000 0000  0000 0000  0000 1111  1100 0000
	PcMask     = 0x0000f000                 // 0000 0000  0000 0000  1111 0000  0000 0000
	CpMask     = 0x000f0000                 // 0000 0000  0000 1111  0000 0000  0000 0000
	PrMask     = 0x00f00000                 // 0000 0000  1111 0000  0000 0000  0000 0000
	EpMask     = 0x0f000000                 // 0000 1111  0000 0000  0000 0000  0000 0000
	CastlMask  = 0xf0000000                 // 1111 0000  0000 0000  0000 0000  0000 0000
	EvalMask   = uint64(0xffff000000000000) // The 16 first bits in uint64
	ToShift    = 6
	PcShift    = 12 // 6+6
	CpShift    = 16 // 6+6+4
	PrShift    = 20 // 6+6+4+4
	EpShift    = 24 // 6+6+4+4+4
	CastlShift = 28 // 6+6+4+4+4+4
	EvalShift  = 64 - 16
)

// Each square is assignes a particular number
const (
	A1 = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1

	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2

	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3

	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4

	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5

	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6

	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7

	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)

// init the square map from string to int and int to string
func InitFen2Sq() {
	Fen2Sq["a1"] = A1
	Fen2Sq["a2"] = A2
	Fen2Sq["a3"] = A3
	Fen2Sq["a4"] = A4
	Fen2Sq["a5"] = A5
	Fen2Sq["a6"] = A6
	Fen2Sq["a7"] = A7
	Fen2Sq["a8"] = A8

	Fen2Sq["b1"] = B1
	Fen2Sq["b2"] = B2
	Fen2Sq["b3"] = B3
	Fen2Sq["b4"] = B4
	Fen2Sq["b5"] = B5
	Fen2Sq["b6"] = B6
	Fen2Sq["b7"] = B7
	Fen2Sq["b8"] = B8

	Fen2Sq["c1"] = C1
	Fen2Sq["c2"] = C2
	Fen2Sq["c3"] = C3
	Fen2Sq["c4"] = C4
	Fen2Sq["c5"] = C5
	Fen2Sq["c6"] = C6
	Fen2Sq["c7"] = C7
	Fen2Sq["c8"] = C8

	Fen2Sq["d1"] = D1
	Fen2Sq["d2"] = D2
	Fen2Sq["d3"] = D3
	Fen2Sq["d4"] = D4
	Fen2Sq["d5"] = D5
	Fen2Sq["d6"] = D6
	Fen2Sq["d7"] = D7
	Fen2Sq["d8"] = D8

	Fen2Sq["e1"] = E1
	Fen2Sq["e2"] = E2
	Fen2Sq["e3"] = E3
	Fen2Sq["e4"] = E4
	Fen2Sq["e5"] = E5
	Fen2Sq["e6"] = E6
	Fen2Sq["e7"] = E7
	Fen2Sq["e8"] = E8

	Fen2Sq["f1"] = F1
	Fen2Sq["f2"] = F2
	Fen2Sq["f3"] = F3
	Fen2Sq["f4"] = F4
	Fen2Sq["f5"] = F5
	Fen2Sq["f6"] = F6
	Fen2Sq["f7"] = F7
	Fen2Sq["f8"] = F8

	Fen2Sq["g1"] = G1
	Fen2Sq["g2"] = G2
	Fen2Sq["g3"] = G3
	Fen2Sq["g4"] = G4
	Fen2Sq["g5"] = G5
	Fen2Sq["g6"] = G6
	Fen2Sq["g7"] = G7
	Fen2Sq["g8"] = G8

	Fen2Sq["h1"] = H1
	Fen2Sq["h2"] = H2
	Fen2Sq["h3"] = H3
	Fen2Sq["h4"] = H4
	Fen2Sq["h5"] = H5
	Fen2Sq["h6"] = H6
	Fen2Sq["h7"] = H7
	Fen2Sq["h8"] = H8

	// -------------- Sq2Fen
	Sq2Fen[A1] = "a1"
	Sq2Fen[A2] = "a2"
	Sq2Fen[A3] = "a3"
	Sq2Fen[A4] = "a4"
	Sq2Fen[A5] = "a5"
	Sq2Fen[A6] = "a6"
	Sq2Fen[A7] = "a7"
	Sq2Fen[A8] = "a8"

	Sq2Fen[B1] = "b1"
	Sq2Fen[B2] = "b2"
	Sq2Fen[B3] = "b3"
	Sq2Fen[B4] = "b4"
	Sq2Fen[B5] = "b5"
	Sq2Fen[B6] = "b6"
	Sq2Fen[B7] = "b7"
	Sq2Fen[B8] = "b8"

	Sq2Fen[C1] = "c1"
	Sq2Fen[C2] = "c2"
	Sq2Fen[C3] = "c3"
	Sq2Fen[C4] = "c4"
	Sq2Fen[C5] = "c5"
	Sq2Fen[C6] = "c6"
	Sq2Fen[C7] = "c7"
	Sq2Fen[C8] = "c8"

	Sq2Fen[D1] = "d1"
	Sq2Fen[D2] = "d2"
	Sq2Fen[D3] = "d3"
	Sq2Fen[D4] = "d4"
	Sq2Fen[D5] = "d5"
	Sq2Fen[D6] = "d6"
	Sq2Fen[D7] = "d7"
	Sq2Fen[D8] = "d8"

	Sq2Fen[E1] = "e1"
	Sq2Fen[E2] = "e2"
	Sq2Fen[E3] = "e3"
	Sq2Fen[E4] = "e4"
	Sq2Fen[E5] = "e5"
	Sq2Fen[E6] = "e6"
	Sq2Fen[E7] = "e7"
	Sq2Fen[E8] = "e8"

	Sq2Fen[F1] = "f1"
	Sq2Fen[F2] = "f2"
	Sq2Fen[F3] = "f3"
	Sq2Fen[F4] = "f4"
	Sq2Fen[F5] = "f5"
	Sq2Fen[F6] = "f6"
	Sq2Fen[F7] = "f7"
	Sq2Fen[F8] = "f8"

	Sq2Fen[G1] = "g1"
	Sq2Fen[G2] = "g2"
	Sq2Fen[G3] = "g3"
	Sq2Fen[G4] = "g4"
	Sq2Fen[G5] = "g5"
	Sq2Fen[G6] = "g6"
	Sq2Fen[G7] = "g7"
	Sq2Fen[G8] = "g8"

	Sq2Fen[H1] = "h1"
	Sq2Fen[H2] = "h2"
	Sq2Fen[H3] = "h3"
	Sq2Fen[H4] = "h4"
	Sq2Fen[H5] = "h5"
	Sq2Fen[H6] = "h6"
	Sq2Fen[H7] = "h7"
	Sq2Fen[H8] = "h8"
}

const (
	// no scoretype = 0
	ScoreTypeLower   = 0x1                             // sc > alpha
	ScoreTypeUpper   = 0x2                             // sc < beta
	ScoreTypeBetween = ScoreTypeLower | ScoreTypeUpper // alpha < sc < beta
)
