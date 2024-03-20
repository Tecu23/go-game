package position

import (
	. "github.com/Tecu23/go-game/pkg/chess/constants"
)

// TODO: eval hash
// TODO: pawn hash
// TODO: pawn structures. isolated, backward, duo, passed (guarded and not), double and more...
// TODO: bishop pair
// TODO: King safety. pawn shelter, guarding pieces
// TODO: King attack. Attacking area surrounding the enemy king, closeness to the enemy king
// TODO: space, center control, knight outposts, connected rooks, 7th row and more
// TODO: combine middle game and end game values

// evaluate returns score from white pov
func Evaluate(b *BoardStruct) int {
	ev := 0
	for sq := A1; sq <= H8; sq++ {
		pc := b.Squares[sq]
		if pc == Empty {
			continue
		}
		ev += PieceVal[pc]
		ev += PcSqScore(pc, sq)
	}
	return ev
}

// Score returns the piece square table value for a given piece on a given square. Stage = MG/EG
func PcSqScore(pc, sq int) int {
	return PSqTab[pc][sq]
}

// PstInit intits the pieces-square-tables when the program starts
func pcSqInit() {
	for pc := 0; pc < 12; pc++ {
		for sq := 0; sq < 64; sq++ {
			PSqTab[pc][sq] = 0
		}
	}

	for sq := 0; sq < 64; sq++ {

		fl := sq % 8
		rk := sq / 8

		PSqTab[WP][sq] = PawnFile[fl] + PawnRank[rk]

		PSqTab[WN][sq] = KnightFile[fl] + KnightRank[rk]
		PSqTab[WB][sq] = CenterFile[fl] + CenterFile[rk]*2

		PSqTab[WR][sq] = CenterFile[fl] * 5

		PSqTab[WQ][sq] = CenterFile[fl] + CenterFile[rk]

		PSqTab[WK][sq] = (KingFile[fl] + KingRank[rk]) * 8
	}

	// bonus for e4 d5 and c4
	PSqTab[WP][E2], PSqTab[WP][D2], PSqTab[WP][E3], PSqTab[WP][D3], PSqTab[WP][E4], PSqTab[WP][D4], PSqTab[WP][C4] = 0, 0, 6, 6, 24, 20, 12

	// long diagonal
	for sq := A1; sq <= H8; sq += NE {
		PSqTab[WB][sq] += LongDiag - 2
	}
	for sq := H1; sq <= A8; sq += NW {
		PSqTab[WB][sq] += LongDiag
	}

	// for Black
	for pt := Pawn; pt <= King; pt++ {

		wPiece := Pt2pc(pt, WHITE)
		bPiece := Pt2pc(pt, BLACK)

		for bSq := 0; bSq < 64; bSq++ {
			wSq := oppRank(bSq)
			PSqTab[bPiece][bSq] = -PSqTab[wPiece][wSq]
		}
	}
}

// mirror the rank_sq
func oppRank(sq int) int {
	fl := sq % 8
	rk := sq / 8
	rk = 7 - rk
	return rk*8 + fl
}
