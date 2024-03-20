// Package position will contain all the function to handle a particular position
package position

import (
	"github.com/Tecu23/go-game/pkg/chess/bitboard"
	"github.com/Tecu23/go-game/pkg/chess/castlings"
	. "github.com/Tecu23/go-game/pkg/chess/constants"
	"github.com/Tecu23/go-game/pkg/chess/magic"
	"github.com/Tecu23/go-game/pkg/chess/moves"
	"github.com/Tecu23/go-game/pkg/chess/trans"
)

// BoardStruct defines all the necessary to generate moves and keep track of a certain position
type BoardStruct struct {
	Key                 uint64
	Squares             [64]int
	WbBB                [2]bitboard.BitBoard         // 1 bb for each side (white or black)
	PieceBB             [NoPiecesT]bitboard.BitBoard // 1 bb for each piece type (R, N, B, Q, K, P)
	King                [2]int                       // the position of each king
	Ep                  int                          // en-passant square
	castlings.Castlings                              // whether each castling is allowed
	Stm                 Color                        // Side To Move
	Count               [NoPiecesC]int               // 12 counters that count how many pieces we have
	Rule50              int                          // set to 0 if a pawn or capt move otherwise increment
}

// Board defines the actual board that the game will be set up and played in
var Board = BoardStruct{}

// AllBB should return all bitboards
func (b *BoardStruct) AllBB() bitboard.BitBoard {
	return b.WbBB[0] | b.WbBB[1]
}

// Clear should clear the board, flags, bitboards etc
func (b *BoardStruct) Clear() {
	b.Stm = WHITE
	b.Rule50 = 0
	b.Squares = [64]int{}
	b.King = [2]int{}
	b.Ep = 0
	b.Castlings = 0

	for i := A1; i <= H8; i++ {
		b.Squares[i] = Empty
	}

	for i := WP; i < NoPiecesC; i++ {
		b.Count[i] = 0
	}

	b.WbBB[WHITE], b.WbBB[BLACK] = 0, 0
	for i := 0; i < NoPiecesT; i++ {
		b.PieceBB[i] = 0
	}

	b.Key = 0
}

// NewGame should start a new game from the starting position
func (b *BoardStruct) NewGame() {
	b.Stm = WHITE
	b.Clear()
	ParseFEN(Startpos)
}

// SetSq should set a square sq to a particular piece pc
func (b *BoardStruct) SetSq(pc, sq int) {
	pt := Pc2pt(pc)
	sd := PcColor(pc)

	// If the square is not empty then it is a capture
	if b.Squares[sq] != Empty {
		cp := b.Squares[sq]
		b.Count[cp]--
		b.WbBB[sd^0x1].Clear(sq)
		b.PieceBB[Pc2pt(cp)].Clear(sq)
		b.Key ^= trans.PcSqKey(cp, sq)
	}

	b.Squares[sq] = pc

	if pc == Empty {
		b.WbBB[WHITE].Clear(sq)
		b.WbBB[BLACK].Clear(sq)
		for p := 0; p < NoPiecesT; p++ {
			b.PieceBB[p].Clear(sq)
		}
		return
	}

	b.Key ^= trans.PcSqKey(pc, sq)

	b.Count[pc]++

	if pt == King {
		b.King[sd] = sq
	}

	b.WbBB[sd].SetBit(sq)
	b.PieceBB[pt].SetBit(sq)
}

// Move should make a move on the board
func (b *BoardStruct) Move(mv moves.Move) bool {
	newEp := 0

	// Assume that the move is legally correct (except for inCheck())
	fr := mv.Fr()
	to := mv.To()
	pr := mv.Pr()
	pc := b.Squares[fr]

	switch {
	case pc == WK:
		b.Castlings.Off(castlings.ShortW | castlings.LongW)

		if Abs(int(to)-int(fr)) == 2 {
			if to == G1 {
				b.SetSq(WR, F1)
				b.SetSq(Empty, H1)
			} else {
				b.SetSq(WR, D1)
				b.SetSq(Empty, A1)
			}
		}
	case pc == BK:
		b.Castlings.Off(castlings.ShortB | castlings.LongB)

		if Abs(int(to)-int(fr)) == 2 {
			if to == G8 {
				b.SetSq(BR, F8)
				b.SetSq(Empty, H8)
			} else {
				b.SetSq(BR, D8)
				b.SetSq(Empty, A8)
			}
		}
	case pc == WR:
		if fr == A1 {
			b.Off(castlings.LongW)
		} else if fr == H1 {
			b.Off(castlings.ShortW)
		}
	case pc == BR:
		if fr == A8 {
			b.Off(castlings.LongB)
		} else if fr == H8 {
			b.Off(castlings.ShortB)
		}
	case pc == WP && b.Squares[to] == Empty: // en passant move or set en passant
		if to-fr == 16 {
			newEp = fr + 8
		} else if to-fr == 7 { // must be ep
			b.SetSq(Empty, to-8)
		} else if to-fr == 9 { // must be ep
			b.SetSq(Empty, to-8)
		}
	case pc == BP && b.Squares[to] == Empty: // en passant move or set en passant
		if fr-to == 16 {
			newEp = to + 8
		} else if fr-to == 7 { // must be ep
			b.SetSq(Empty, to+8)
		} else if fr-to == 9 { // must be ep
			b.SetSq(Empty, to+8)
		}
	}

	b.Ep = newEp
	b.SetSq(Empty, fr)

	if pr != Empty {
		b.SetSq(pr, to)
	} else {
		b.SetSq(pc, to)
	}

	b.Key = ^b.Key
	b.Stm = b.Stm ^ 0x1

	// Undo the move if the king is in check
	if b.isAttacked(b.King[b.Stm^0x1], b.Stm) {
		b.Unmove(mv)
		return false
	}

	return true
}

// Unmove undoes a certain move
func (b *BoardStruct) Unmove(mv moves.Move) {
	b.Ep = mv.Ep(b.Stm.Opposite())
	b.Castlings = mv.Castl()

	pc := int(mv.Pc())
	fr := int(mv.Fr())
	to := int(mv.To())

	b.SetSq(mv.Cp(), to)
	b.SetSq(pc, fr)

	if Pc2pt(pc) == Pawn {
		if to == b.Ep && b.Ep != 0 {
			b.SetSq(Empty, to)
			switch to - fr {
			case NW, NE:
				b.SetSq(BP, to-N)
			case SW, SE:
				b.SetSq(WP, to-S)

			}
		}
	} else if Pc2pt(pc) == King {
		sd := PcColor(pc)
		if fr-to == 2 {
			b.SetSq(castlings.Castl[sd].Rook, int(castlings.Castl[sd].RookL))
			b.SetSq(Empty, fr-1)
		} else if fr-to == -2 {
			b.SetSq(castlings.Castl[sd].Rook, int(castlings.Castl[sd].RookSh))
			b.SetSq(Empty, fr+1)
		}
	}
	b.Key = ^b.Key
	b.Stm = b.Stm ^ 0x1
}

// isAttacked should return whether the sq is attacked by the sd color side
func (b *BoardStruct) isAttacked(to int, sd Color) bool {
	if isPawnAtkingSq[sd](b, to) {
		return true
	}

	if AtksKnights[to]&b.PieceBB[Knight]&b.WbBB[sd] != 0 {
		return true
	}
	if AtksKings[to]&b.PieceBB[King]&b.WbBB[sd] != 0 {
		return true
	}
	if (magic.MBishopTab[to].Atks(b.AllBB()) & (b.PieceBB[Bishop] | b.PieceBB[Queen]) & b.WbBB[sd]) != 0 {
		return true
	}
	if (magic.MRookTab[to].Atks(b.AllBB()) & (b.PieceBB[Rook] | b.PieceBB[Queen]) & b.WbBB[sd]) != 0 {
		return true
	}

	return false
}
