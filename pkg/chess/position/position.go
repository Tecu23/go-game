// Package position will contain all the function to handle a particular position
package position

import (
	"github.com/Tecu23/go-game/pkg/chess/bitboard"
	"github.com/Tecu23/go-game/pkg/chess/castlings"
	. "github.com/Tecu23/go-game/pkg/chess/constants"
	"github.com/Tecu23/go-game/pkg/chess/magic"
	"github.com/Tecu23/go-game/pkg/chess/moves"
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
		b.Key ^= PcSqKey(cp, sq)
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

	b.Key ^= PcSqKey(pc, sq)

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
	if b.IsAttacked(b.King[b.Stm^0x1], b.Stm) {
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
func (b *BoardStruct) IsAttacked(to int, sd Color) bool {
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

func (b *BoardStruct) MoveNull() moves.Move {
	mv := moves.NoMove
	mv.PackMove(0, 0, Empty, Empty, Empty, b.Ep, b.Castlings)

	b.Ep = 0
	b.Key = ^b.Key
	b.Stm = b.Stm ^ 0x1
	return mv
}

func (b *BoardStruct) UndoNull(mv moves.Move) {
	b.Key = ^b.Key
	b.Stm = b.Stm ^ 0x1

	b.Ep = mv.Ep(b.Stm)
	// b.castlings = mv.castl()     // no need!
}

// is the move legal (except from inCheck)
func (b *BoardStruct) IsLegal(mv moves.Move) bool {
	fr := mv.Fr()
	pc := mv.Pc()
	if b.Squares[fr] != pc || pc == Empty {
		return false
	}
	if b.Stm != PcColor(pc) {
		return false
	}

	to := mv.To()
	cp := mv.Cp()
	if !((pc == WP || pc == BP) && to == b.Ep && b.Ep != 0) {
		if b.Squares[to] != cp {
			return false
		}
		if cp != Empty && PcColor(cp) == PcColor(pc) {
			return false
		}
	}

	switch {
	case pc == WP:
		if to-fr == 8 { // wP one step
			if b.Squares[to] == Empty {
				return true
			}
		} else if to-fr == 16 {
			if b.Squares[fr+8] == Empty && b.Squares[fr+16] == Empty { // wP two step
				return true
			}
		} else if b.Ep == mv.Ep(b.Stm) && b.Squares[to-8] == BP { // wP ep
			return true
		} else if to-fr == 7 && cp != Empty { // wP capture left
			return true
		} else if to-fr == 9 && cp != Empty { // wp capture right
			return true
		}

		return false
	case pc == BP:
		if fr-to == 8 { // bP one step
			if b.Squares[to] == Empty {
				return true
			}
		} else if fr-to == 16 {
			if b.Squares[fr-8] == Empty && b.Squares[fr-16] == Empty { // bP two step
				return true
			}
		} else if b.Ep == mv.Ep(b.Stm) && b.Squares[to+8] == WP { // bP ep
			return true
		} else if fr-to == 7 && cp != Empty { // bP capture right
			return true
		} else if fr-to == 9 && cp != Empty { // bp capture left
			return true
		}

		return false
	case pc == WB, pc == BB:
		toBB := bitboard.BitBoard(1) << uint(to)
		if magic.MBishopTab[fr].Atks(b.AllBB())&toBB != 0 {
			return true
		}
		return false
	case pc == WR, pc == BR:
		toBB := bitboard.BitBoard(1) << uint(to)
		if magic.MRookTab[fr].Atks(b.AllBB())&toBB != 0 {
			return true
		}
		return false
	case pc == WQ, pc == BQ:
		toBB := bitboard.BitBoard(1) << uint(to)
		if magic.MBishopTab[fr].Atks(b.AllBB())&toBB != 0 {
			return true
		}
		if magic.MRookTab[fr].Atks(b.AllBB())&toBB != 0 {
			return true
		}
		return false
	case pc == WK:
		if Abs(int(to)-int(fr)) == 2 { // castlings
			if to == G1 {
				if b.Squares[H1] != WR || b.Squares[E1] != WK {
					return false
				}

				if b.Squares[F1] != Empty || b.Squares[G1] != Empty {
					return false
				}

				if !b.IsShortOk(b.Stm) {
					return false
				}
			} else {
				if b.Squares[A1] != WR || b.Squares[E1] != WK {
					return false
				}
				if to != C1 {
					return false
				}
				if b.Squares[B1] != Empty || b.Squares[C1] != Empty || b.Squares[D1] != Empty {
					return false
				}
				if !b.IsLongOk(b.Stm) {
					return false
				}
			}
		}
		return true
	case pc == BK:
		if Abs(int(to)-int(fr)) == 2 { // castlings
			if to == G8 {
				if b.Squares[H8] != BR || b.Squares[E8] != BK {
					return false
				}
				if b.Squares[F8] != Empty || b.Squares[G8] != Empty {
					return false
				}
				if !b.IsShortOk(b.Stm) {
					return false
				}
			} else {
				if b.Squares[A8] != BR || b.Squares[E8] != BK {
					return false
				}
				if to != C8 {
					return false
				}
				if b.Squares[B8] != Empty || b.Squares[C8] != Empty || b.Squares[D8] != Empty {
					return false
				}
				if !b.IsLongOk(b.Stm) {
					return false
				}
			}
		}
		return true
	}

	return true
}

// check if short castlings is legal
func (b *BoardStruct) IsShortOk(sd Color) bool {
	if !b.ShortFlag(sd) {
		return false
	}

	opp := sd ^ 0x1
	if castlings.Castl[sd].PawnsSh&b.PieceBB[Pawn]&b.WbBB[opp] != 0 { // stopped by pawns?
		return false
	}
	if castlings.Castl[sd].PawnsSh&b.PieceBB[King]&b.WbBB[opp] != 0 { // stopped by king?
		return false
	}
	if castlings.Castl[sd].KnightsSh&b.PieceBB[Knight]&b.WbBB[opp] != 0 { // stopped by Knights?
		return false
	}

	// sliding to e1/e8	//NOTE: Maybe not needed during search because we know if we are in check
	sq := b.King[sd]
	if (magic.MBishopTab[sq].Atks(b.AllBB()) & (b.PieceBB[Bishop] | b.PieceBB[Queen]) & b.WbBB[opp]) != 0 {
		return false
	}
	if (magic.MRookTab[sq].Atks(b.AllBB()) & (b.PieceBB[Rook] | b.PieceBB[Queen]) & b.WbBB[opp]) != 0 {
		return false
	}

	// slidings to f1/f8
	if (magic.MBishopTab[sq+1].Atks(b.AllBB()) & (b.PieceBB[Bishop] | b.PieceBB[Queen]) & b.WbBB[opp]) != 0 {
		return false
	}
	if (magic.MRookTab[sq+1].Atks(b.AllBB()) & (b.PieceBB[Rook] | b.PieceBB[Queen]) & b.WbBB[opp]) != 0 {
		return false
	}

	// slidings to g1/g8		//NOTE: Maybe not needed because we always make isAttacked() after a move
	if (magic.MBishopTab[sq+2].Atks(b.AllBB()) & (b.PieceBB[Bishop] | b.PieceBB[Queen]) & b.WbBB[opp]) != 0 {
		return false
	}
	if (magic.MRookTab[sq+2].Atks(b.AllBB()) & (b.PieceBB[Rook] | b.PieceBB[Queen]) & b.WbBB[opp]) != 0 {
		return false
	}
	return true
}

// check if long castlings is legal
func (b *BoardStruct) IsLongOk(sd Color) bool {
	if !b.LongFlag(sd) {
		return false
	}

	opp := sd ^ 0x1
	if castlings.Castl[sd].PawnsL&b.PieceBB[Pawn]&b.WbBB[opp] != 0 {
		return false
	}
	if castlings.Castl[sd].PawnsL&b.PieceBB[King]&b.WbBB[opp] != 0 {
		return false
	}
	if castlings.Castl[sd].KnightsL&b.PieceBB[Knight]&b.WbBB[opp] != 0 {
		return false
	}

	// sliding e1/e8
	sq := b.King[sd]
	if (magic.MBishopTab[sq].Atks(b.AllBB()) & (b.PieceBB[Bishop] | b.PieceBB[Queen]) & b.WbBB[opp]) != 0 {
		return false
	}
	if (magic.MRookTab[sq].Atks(b.AllBB()) & (b.PieceBB[Rook] | b.PieceBB[Queen]) & b.WbBB[opp]) != 0 {
		return false
	}

	// sliding d1/d8
	if (magic.MBishopTab[sq-1].Atks(b.AllBB()) & (b.PieceBB[Bishop] | b.PieceBB[Queen]) & b.WbBB[opp]) != 0 {
		return false
	}
	if (magic.MRookTab[sq-1].Atks(b.AllBB()) & (b.PieceBB[Rook] | b.PieceBB[Queen]) & b.WbBB[opp]) != 0 {
		return false
	}

	// sliding c1/c8	//NOTE: Maybe not needed because we always make inCheck() before a move
	if (magic.MBishopTab[sq-2].Atks(b.AllBB()) & (b.PieceBB[Bishop] | b.PieceBB[Queen]) & b.WbBB[opp]) != 0 {
		return false
	}
	if (magic.MRookTab[sq-2].Atks(b.AllBB()) & (b.PieceBB[Rook] | b.PieceBB[Queen]) & b.WbBB[opp]) != 0 {
		return false
	}
	return true
}

// is this a position to avoid null move?
func (b *BoardStruct) IsAntiNullMove() bool {
	if b.WbBB[b.Stm] == b.PieceBB[King]&b.WbBB[b.Stm] {
		return true
	}
	return false
}

// ////////////////////////////// TODO: remove this after benchmarking ////////////////////////////////////////
func (b *BoardStruct) GenSimpleRookMoves(ml *moves.MoveList, sd Color) {
	allRBB := b.PieceBB[Rook] & b.WbBB[sd]
	pc := Pt2pc(Rook, Color(sd))
	ep := b.Ep
	castlings := b.Castlings
	var mv moves.Move
	for fr := allRBB.FirstOne(); fr != 64; fr = allRBB.FirstOne() {
		rk := fr / 8
		fl := fr % 8
		// N
		for r := rk + 1; r < 8; r++ {
			to := r*8 + fl
			cp := b.Squares[to]
			if cp != Empty && PcColor(int(cp)) == sd {
				break
			}
			mv.PackMove(fr, to, pc, cp, Empty, ep, castlings)
			ml.Add(mv)
			if cp != Empty {
				break
			}
		}
		// S
		for r := rk - 1; r >= 0; r-- {
			to := r*8 + fl
			cp := b.Squares[to]
			if cp != Empty && PcColor(int(cp)) == sd {
				break
			}
			mv.PackMove(fr, to, pc, cp, Empty, ep, castlings)
			ml.Add(mv)
			if cp != Empty {
				break
			}
		}
		// E
		for f := fl + 1; f < 8; f++ {
			to := rk*8 + f
			cp := b.Squares[to]
			if cp != Empty && PcColor(int(cp)) == sd {
				break
			}
			mv.PackMove(fr, to, pc, cp, Empty, ep, castlings)
			ml.Add(mv)
			if cp != Empty {
				break
			}
		}
		// W
		for f := fl - 1; f >= 0; f-- {
			to := rk*8 + f
			cp := b.Squares[to]
			if cp != Empty && PcColor(int(cp)) == sd {
				break
			}
			mv.PackMove(fr, to, pc, cp, Empty, ep, castlings)
			ml.Add(mv)
			if cp != Empty {
				break
			}
		}
	}
}
