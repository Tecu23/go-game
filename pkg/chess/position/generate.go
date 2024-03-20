package position

import (
	"github.com/Tecu23/go-game/pkg/chess/bitboard"
	"github.com/Tecu23/go-game/pkg/chess/castlings"
	. "github.com/Tecu23/go-game/pkg/chess/constants"
	"github.com/Tecu23/go-game/pkg/chess/magic"
	"github.com/Tecu23/go-game/pkg/chess/moves"
)

var (
	AtksKnights [64]bitboard.BitBoard
	AtksKings   [64]bitboard.BitBoard
)

var (
	isPawnAtkingSq = [2]func(*BoardStruct, int) bool{
		(*BoardStruct).iswPawnAtkingSq,
		(*BoardStruct).isbPawnAtkingSq,
	}
	allPawnAtksBB = [2]func(*BoardStruct) bitboard.BitBoard{
		(*BoardStruct).wPawnAtksBB,
		(*BoardStruct).bPawnAtksBB,
	}
	pawnAtksFr = [2]func(*BoardStruct, int) bitboard.BitBoard{
		(*BoardStruct).WPawnAtksFr,
		(*BoardStruct).BPawnAtksFr,
	}
	pawnAtkers = [2]func(*BoardStruct) bitboard.BitBoard{
		(*BoardStruct).wPawnAtkers,
		(*BoardStruct).bPawnAtkers,
	}
)

// Returns true or false if to-sq is attacked by white pawn
func (b *BoardStruct) iswPawnAtkingSq(to int) bool {
	sqBB := bitboard.BitBoard(1) << uint(to)

	wPawns := b.PieceBB[Pawn] & b.WbBB[WHITE]

	// Attacks left and right
	toCap := ((wPawns & ^FileA) << NW) & b.WbBB[BLACK]
	toCap |= ((wPawns & ^FileH) << NE) & b.WbBB[BLACK]
	return (toCap & sqBB) != 0
}

// Returns true or false if to-sq is attacked by white pawn
func (b *BoardStruct) isbPawnAtkingSq(to int) bool {
	sqBB := bitboard.BitBoard(1) << uint(to)

	bPawns := b.PieceBB[Pawn] & b.WbBB[BLACK]

	// Attacks left and right
	toCap := ((bPawns & ^FileA) >> (-SW)) & b.WbBB[WHITE]
	toCap |= ((bPawns & ^FileH) >> (-SE)) & b.WbBB[WHITE]

	return (toCap & sqBB) != 0
}

// returns all w pawns that attacka black pieces
func (b *BoardStruct) wPawnAtkers() bitboard.BitBoard {
	BB := b.WbBB[BLACK] // all their pieces
	// pretend that all their pieces are pawns
	// Get pawn Attacks left and right from their pieces into our pawns that now are all our pwan attackers
	ourPawnAttackers := ((BB & ^FileA) >> (-SW)) & b.WbBB[WHITE] & b.PieceBB[Pawn]
	ourPawnAttackers |= ((BB & ^FileH) >> (-SE)) & b.WbBB[WHITE] & b.PieceBB[Pawn]

	return ourPawnAttackers
}

// returns all bl pawns that attacks white pieces
func (b *BoardStruct) bPawnAtkers() bitboard.BitBoard {
	BB := b.WbBB[WHITE] // all their pieces
	// pretend that all their pieces are pawns
	// Get pawn Attacks left and right from their pieces into our pawns that now are all our pwan attackers
	ourPawnAttackers := ((BB & ^FileA) << NW) & b.WbBB[BLACK] & b.PieceBB[Pawn]
	ourPawnAttackers |= ((BB & ^FileH) << NE) & b.WbBB[BLACK] & b.PieceBB[Pawn]

	return ourPawnAttackers
}

// returns captures from fr-sq
func (b *BoardStruct) WPawnAtksFr(fr int) bitboard.BitBoard {
	frBB := bitboard.BitBoard(1) << uint(fr)

	// Attacks left and right
	toCap := ((frBB & ^FileA) << NW) & b.WbBB[BLACK]
	toCap |= ((frBB & ^FileH) << NE) & b.WbBB[BLACK]
	return toCap
}

// returns captures from fr-sq
func (b *BoardStruct) BPawnAtksFr(fr int) bitboard.BitBoard {
	frBB := bitboard.BitBoard(1) << uint(fr)

	// Attacks left and right
	toCap := ((frBB & ^FileA) >> (-SW)) & b.WbBB[WHITE]
	toCap |= ((frBB & ^FileH) >> (-SE)) & b.WbBB[WHITE]

	return toCap
}

// returns bitBoard with all attacks, empty or not, from all white Pawns
func (b *BoardStruct) wPawnAtksBB() bitboard.BitBoard {
	frBB := b.PieceBB[Pawn] & b.WbBB[WHITE]

	// Attacks left and right
	toCap := ((frBB & ^FileA) << NW)
	toCap |= ((frBB & ^FileH) << NE)
	return toCap
}

// returns bitBoard with all attacks, empty or not, from all white Pawns
func (b *BoardStruct) bPawnAtksBB() bitboard.BitBoard {
	frBB := b.PieceBB[Pawn] & b.WbBB[BLACK]

	// Attacks left and right
	toCap := ((frBB & ^FileA) << NW)
	toCap |= ((frBB & ^FileH) << NE)
	return toCap
}

func (b *BoardStruct) GenRookMoves(ml *moves.MoveList, targetBB bitboard.BitBoard) {
	sd := b.Stm
	allRBB := b.PieceBB[Rook] & b.WbBB[sd]
	pc := Pt2pc(Rook, Color(sd))
	var mv moves.Move
	for fr := allRBB.FirstOne(); fr != 64; fr = allRBB.FirstOne() {
		toBB := magic.MRookTab[fr].Atks(b.AllBB()) & targetBB
		for to := toBB.FirstOne(); to != 64; to = toBB.FirstOne() {
			mv.PackMove(fr, to, pc, b.Squares[to], Empty, b.Ep, b.Castlings)
			ml.Add(mv)
		}
	}
}

func (b *BoardStruct) GenBishopMoves(ml *moves.MoveList, targetBB bitboard.BitBoard) {
	sd := b.Stm
	allBBB := b.PieceBB[Bishop] & b.WbBB[sd]
	pc := Pt2pc(Bishop, Color(sd))
	ep := b.Ep
	castlings := b.Castlings
	var mv moves.Move

	for fr := allBBB.FirstOne(); fr != 64; fr = allBBB.FirstOne() {
		toBB := magic.MBishopTab[fr].Atks(b.AllBB()) & targetBB
		for to := toBB.LastOne(); to != 64; to = toBB.LastOne() {
			mv.PackMove(fr, to, pc, b.Squares[to], Empty, ep, castlings)
			ml.Add(mv)
		}
	}
}

func (b *BoardStruct) GenQueenMoves(mlq *moves.MoveList, targetBB bitboard.BitBoard) {
	sd := b.Stm
	allQBB := b.PieceBB[Queen] & b.WbBB[sd]
	pc := Pt2pc(Queen, Color(sd))
	ep := b.Ep
	castlings := b.Castlings
	var mv moves.Move

	for fr := allQBB.FirstOne(); fr != 64; fr = allQBB.FirstOne() {
		toBB := magic.MBishopTab[fr].Atks(b.AllBB()) & targetBB
		toBB |= magic.MRookTab[fr].Atks(b.AllBB()) & targetBB
		for to := toBB.FirstOne(); to != 64; to = toBB.FirstOne() {
			mv.PackMove(fr, to, pc, b.Squares[to], Empty, ep, castlings)
			mlq.Add(mv)
		}
	}
}

func (b *BoardStruct) GenKnightMoves(ml *moves.MoveList, targetBB bitboard.BitBoard) {
	sd := b.Stm
	allNBB := b.PieceBB[Knight] & b.WbBB[sd]
	pc := Pt2pc(Knight, Color(sd))
	ep := b.Ep
	castlings := b.Castlings
	var mv moves.Move
	for fr := allNBB.FirstOne(); fr != 64; fr = allNBB.FirstOne() {
		toBB := AtksKnights[fr] & targetBB
		for to := toBB.FirstOne(); to != 64; to = toBB.FirstOne() {
			mv.PackMove(fr, to, pc, b.Squares[to], Empty, ep, castlings)
			ml.Add(mv)
		}
	}
}

func (b *BoardStruct) GenKingMoves(ml *moves.MoveList, targetBB bitboard.BitBoard) {
	sd := b.Stm
	// 'normal' moves
	pc := Pt2pc(King, Color(sd))
	ep := b.Ep
	Castlings := b.Castlings
	var mv moves.Move

	toBB := AtksKings[b.King[sd]] & targetBB
	for to := toBB.FirstOne(); to != 64; to = toBB.FirstOne() {
		mv.PackMove(b.King[sd], to, pc, b.Squares[to], Empty, ep, Castlings)
		ml.Add(mv)
	}

	// castlings
	if b.King[sd] == castlings.Castl[sd].KingPos { // NOTE: Maybe not needed. We should know that the king is there if the flags are ok
		if targetBB.IsBitSet(b.King[sd] + 2) {
			// short castling
			if b.Squares[castlings.Castl[sd].RookSh] == castlings.Castl[sd].Rook && // NOTE: Maybe not needed. We should know that the rook is there if the flags are ok
				(castlings.Castl[sd].BetweenSh&b.AllBB()) == 0 {
				if b.IsShortOk(sd) {
					mv.PackMove(
						b.King[sd],
						b.King[sd]+2,
						b.Squares[b.King[sd]],
						Empty,
						Empty,
						b.Ep,
						b.Castlings,
					)
					ml.Add(mv)
				}
			}
		}

		if targetBB.IsBitSet(b.King[sd] - 2) {
			// long castling
			if b.Squares[castlings.Castl[sd].RookL] == castlings.Castl[sd].Rook && // NOTE: Maybe not needed. We should know that the rook is there if the flags are ok
				(castlings.Castl[sd].BetweenL&b.AllBB()) == 0 {
				if b.IsLongOk(sd) {
					mv.PackMove(
						b.King[sd],
						b.King[sd]-2,
						b.Squares[b.King[sd]],
						Empty,
						Empty,
						b.Ep,
						b.Castlings,
					)
					ml.Add(mv)
				}
			}
		}
	}
}

var (
	GenPawns = [2]func(*BoardStruct, *moves.MoveList){
		(*BoardStruct).GenWPawnMoves,
		(*BoardStruct).GenBPawnMoves,
	}
	GenPawnCapt = [2]func(*BoardStruct, *moves.MoveList){
		(*BoardStruct).GenWPawnCapt,
		(*BoardStruct).GenBPawnCapt,
	}
	GenPawnNonCapt = [2]func(*BoardStruct, *moves.MoveList){
		(*BoardStruct).GenWPawnNonCapt,
		(*BoardStruct).GenBPawnNonCapt,
	}
)

func (b *BoardStruct) GenPawnMoves(ml *moves.MoveList) {
	GenPawns[b.Stm](b, ml)
}

func (b *BoardStruct) GenPawnCapt(ml *moves.MoveList) {
	GenPawnCapt[b.Stm](b, ml)
}

func (b *BoardStruct) GenPawnNonCapt(ml *moves.MoveList) {
	GenPawnNonCapt[b.Stm](b, ml)
}

func (b *BoardStruct) GenWPawnMoves(ml *moves.MoveList) {
	wPawns := b.PieceBB[Pawn] & b.WbBB[WHITE]

	// one step
	to1Step := (wPawns << N) & ^b.AllBB()
	// two steps,
	to2Step := ((to1Step & Row3) << N) & ^b.AllBB()
	// captures
	toCapL := ((wPawns & ^FileA) << NW) & b.WbBB[BLACK]
	toCapR := ((wPawns & ^FileH) << NE) & b.WbBB[BLACK]

	mv := moves.NoMove

	// prom
	prom := (to1Step | toCapL | toCapR) & Row8
	if prom != 0 {
		for to := prom.FirstOne(); to != 64; to = prom.FirstOne() {
			cp := b.Squares[to]
			frTab := make([]int, 0, 3)
			if b.Squares[to] == Empty {
				frTab = append(frTab, to-N) // not capture
			} else {
				if toCapL.IsBitSet(to) { // capture left
					frTab = append(frTab, to-NW)
				}
				if toCapR.IsBitSet(to) { // capture right
					frTab = append(frTab, to-NE)
				}
			}

			for _, fr := range frTab {
				mv.PackMove(fr, to, WP, cp, WQ, b.Ep, b.Castlings)
				ml.Add(mv)
				mv.PackMove(fr, to, WP, cp, WR, b.Ep, b.Castlings)
				ml.Add(mv)
				mv.PackMove(fr, to, WP, cp, WN, b.Ep, b.Castlings)
				ml.Add(mv)
				mv.PackMove(fr, to, WP, cp, WB, b.Ep, b.Castlings)
				ml.Add(mv)
			}
		}
		to1Step &= ^Row8
		toCapL &= ^Row8
		toCapR &= ^Row8
	}

	// ep move
	if b.Ep != 0 {
		epBB := bitboard.BitBoard(1) << uint(b.Ep)
		// ep left
		epToL := ((wPawns & ^FileA) << NW) & epBB
		if epToL != 0 {
			mv.PackMove(b.Ep-NW, b.Ep, WP, BP, Empty, b.Ep, b.Castlings)
			ml.Add(mv)
		}
		epToR := ((wPawns & ^FileH) << NE) & epBB
		if epToR != 0 {
			mv.PackMove(b.Ep-NE, b.Ep, WP, BP, Empty, b.Ep, b.Castlings)
			ml.Add(mv)
		}
	}
	// Add one step forward
	for to := to1Step.FirstOne(); to != 64; to = to1Step.FirstOne() {
		mv.PackMove(to-N, to, WP, Empty, Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}
	// Add two steps forward
	for to := to2Step.FirstOne(); to != 64; to = to2Step.FirstOne() {
		mv.PackMove(to-2*N, to, WP, Empty, Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}

	// add Captures left
	for to := toCapL.FirstOne(); to != 64; to = toCapL.FirstOne() {
		mv.PackMove(to-NW, to, WP, b.Squares[to], Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}

	// add Captures right
	for to := toCapR.FirstOne(); to != 64; to = toCapR.FirstOne() {
		mv.PackMove(to-NE, to, WP, b.Squares[to], Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}
}

func (b *BoardStruct) GenBPawnMoves(ml *moves.MoveList) {
	bPawns := b.PieceBB[Pawn] & b.WbBB[BLACK]

	// one step
	to1Step := (bPawns >> (-S)) & ^b.AllBB()
	// two steps,
	to2Step := ((to1Step & Row6) >> (-S)) & ^b.AllBB()
	// captures
	toCapL := ((bPawns & ^FileA) >> (-SW)) & b.WbBB[WHITE]
	toCapR := ((bPawns & ^FileH) >> (-SE)) & b.WbBB[WHITE]

	var mv moves.Move

	// prom
	prom := (to1Step | toCapL | toCapR) & Row1
	if prom != 0 {
		for to := prom.FirstOne(); to != 64; to = prom.FirstOne() {
			cp := b.Squares[to]
			frTab := make([]int, 0, 3)
			if b.Squares[to] == Empty {
				frTab = append(frTab, to-S) // not capture
			} else {
				if toCapL.IsBitSet(to) { // capture left
					frTab = append(frTab, to-SW)
				}
				if toCapR.IsBitSet(to) { // capture right
					frTab = append(frTab, to-SE)
				}
			}

			for _, fr := range frTab {
				mv.PackMove(fr, to, BP, cp, BQ, b.Ep, b.Castlings)
				ml.Add(mv)
				mv.PackMove(fr, to, BP, cp, BR, b.Ep, b.Castlings)
				ml.Add(mv)
				mv.PackMove(fr, to, BP, cp, BN, b.Ep, b.Castlings)
				ml.Add(mv)
				mv.PackMove(fr, to, BP, cp, BB, b.Ep, b.Castlings)
				ml.Add(mv)
			}
		}
		to1Step &= ^Row1
		toCapL &= ^Row1
		toCapR &= ^Row1
	}
	// ep move
	if b.Ep != 0 {
		epBB := bitboard.BitBoard(1) << uint(b.Ep)
		// ep left
		epToL := ((bPawns & ^FileA) >> (-SW)) & epBB
		if epToL != 0 {
			mv.PackMove(b.Ep-SW, b.Ep, BP, WP, Empty, b.Ep, b.Castlings)
			ml.Add(mv)
		}
		epToR := ((bPawns & ^FileH) >> (-SE)) & epBB
		if epToR != 0 {
			mv.PackMove(b.Ep-SE, b.Ep, BP, WP, Empty, b.Ep, b.Castlings)
			ml.Add(mv)
		}
	}
	// Add one step forward
	for to := to1Step.FirstOne(); to != 64; to = to1Step.FirstOne() {
		mv.PackMove(to-S, to, BP, Empty, Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}
	// Add two steps forward
	for to := to2Step.FirstOne(); to != 64; to = to2Step.FirstOne() {
		mv.PackMove(to-2*S, to, BP, Empty, Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}

	// add Captures left
	for to := toCapL.FirstOne(); to != 64; to = toCapL.FirstOne() {
		mv.PackMove(to-SW, to, BP, b.Squares[to], Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}

	// add Captures right
	for to := toCapR.FirstOne(); to != 64; to = toCapR.FirstOne() {
		mv.PackMove(to-SE, to, BP, b.Squares[to], Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}
}

// W pawns  captures or promotions alt 2
func (b *BoardStruct) GenWPawnCapt(ml *moves.MoveList) {
	wPawns := b.PieceBB[Pawn] & b.WbBB[WHITE]

	// captures
	toCapL := ((wPawns & ^FileA) << NW) & b.WbBB[BLACK]
	toCapR := ((wPawns & ^FileH) << NE) & b.WbBB[BLACK]
	// prom
	prom := Row8 & ((toCapL | toCapR) | ((wPawns << N) & ^b.AllBB()))

	var mv moves.Move
	if prom != 0 {
		for to := prom.FirstOne(); to != 64; to = prom.FirstOne() {
			cp := b.Squares[to]
			frTab := make([]int, 0, 3)
			if b.Squares[to] == Empty {
				frTab = append(frTab, to-N) // not capture
			} else {
				if toCapL.IsBitSet(to) { // capture left
					frTab = append(frTab, to-NW)
				}
				if toCapR.IsBitSet(to) { // capture right
					frTab = append(frTab, to-NE)
				}
			}
			for _, fr := range frTab {
				mv.PackMove(fr, to, WP, cp, WQ, b.Ep, b.Castlings)
				ml.Add(mv)
				mv.PackMove(fr, to, WP, cp, WR, b.Ep, b.Castlings)
				ml.Add(mv)
				mv.PackMove(fr, to, WP, cp, WN, b.Ep, b.Castlings)
				ml.Add(mv)
				mv.PackMove(fr, to, WP, cp, WB, b.Ep, b.Castlings)
				ml.Add(mv)
			}
		}
		toCapL &= ^Row8
		toCapR &= ^Row8
	}
	// ep move
	if b.Ep != 0 {
		epBB := bitboard.BitBoard(1) << uint(b.Ep)
		// ep left
		epToL := ((wPawns & ^FileA) << NW) & epBB
		if epToL != 0 {
			mv.PackMove(b.Ep-NW, b.Ep, WP, BP, Empty, b.Ep, b.Castlings)
			ml.Add(mv)
		}
		epToR := ((wPawns & ^FileH) << NE) & epBB
		if epToR != 0 {
			mv.PackMove(b.Ep-NE, b.Ep, WP, BP, Empty, b.Ep, b.Castlings)
			ml.Add(mv)
		}
	}

	// add Captures left
	for to := toCapL.FirstOne(); to != 64; to = toCapL.FirstOne() {
		mv.PackMove(to-NW, to, WP, b.Squares[to], Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}

	// add Captures right
	for to := toCapR.FirstOne(); to != 64; to = toCapR.FirstOne() {
		mv.PackMove(to-NE, to, WP, b.Squares[to], Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}
}

// B pawn captures or promotions alternativ 2
func (b *BoardStruct) GenBPawnCapt(ml *moves.MoveList) {
	bPawns := b.PieceBB[Pawn] & b.WbBB[BLACK]

	// captures
	toCapL := ((bPawns & ^FileA) >> (-SW)) & b.WbBB[WHITE]
	toCapR := ((bPawns & ^FileH) >> (-SE)) & b.WbBB[WHITE]

	var mv moves.Move

	// prom
	prom := Row1 & ((toCapL | toCapR) | ((bPawns >> (-S)) & ^b.AllBB()))
	if prom != 0 {
		for to := prom.FirstOne(); to != 64; to = prom.FirstOne() {
			cp := b.Squares[to]
			frTab := make([]int, 0, 3)
			if b.Squares[to] == Empty {
				frTab = append(frTab, to-S) // not capture
			} else {
				if toCapL.IsBitSet(to) { // capture left
					frTab = append(frTab, to-SW)
				}
				if toCapR.IsBitSet(to) { // capture right
					frTab = append(frTab, to-SE)
				}
			}

			for _, fr := range frTab {
				mv.PackMove(fr, to, BP, cp, BQ, b.Ep, b.Castlings)
				ml.Add(mv)
				mv.PackMove(fr, to, BP, cp, BR, b.Ep, b.Castlings)
				ml.Add(mv)
				mv.PackMove(fr, to, BP, cp, BN, b.Ep, b.Castlings)
				ml.Add(mv)
				mv.PackMove(fr, to, BP, cp, BB, b.Ep, b.Castlings)
				ml.Add(mv)
			}
		}
		toCapL &= ^Row1
		toCapR &= ^Row1
	}
	// ep move
	if b.Ep != 0 {
		epBB := bitboard.BitBoard(1) << uint(b.Ep)
		// ep left
		epToL := ((bPawns & ^FileA) >> (-SW)) & epBB
		if epToL != 0 {
			mv.PackMove(b.Ep-SW, b.Ep, BP, WP, Empty, b.Ep, b.Castlings)
			ml.Add(mv)
		}
		epToR := ((bPawns & ^FileH) >> (-SE)) & epBB
		if epToR != 0 {
			mv.PackMove(b.Ep-SE, b.Ep, BP, WP, Empty, b.Ep, b.Castlings)
			ml.Add(mv)
		}
	}

	// add Captures left
	for to := toCapL.FirstOne(); to != 64; to = toCapL.FirstOne() {
		mv.PackMove(to-SW, to, BP, b.Squares[to], Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}

	// add Captures right
	for to := toCapR.FirstOne(); to != 64; to = toCapR.FirstOne() {
		mv.PackMove(to-SE, to, BP, b.Squares[to], Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}
}

// W pawns moves that doesn't capture aand not promotions
func (b *BoardStruct) GenWPawnNonCapt(ml *moves.MoveList) {
	var mv moves.Move
	wPawns := b.PieceBB[Pawn] & b.WbBB[WHITE]

	// one step
	to1Step := (wPawns << N) & ^b.AllBB()
	// two steps,
	to2Step := ((to1Step & Row3) << N) & ^b.AllBB()
	to1Step &= ^Row8

	// Add one step forward
	for to := to1Step.FirstOne(); to != 64; to = to1Step.FirstOne() {
		mv.PackMove(to-N, to, WP, Empty, Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}
	// Add two steps forward
	for to := to2Step.FirstOne(); to != 64; to = to2Step.FirstOne() {
		mv.PackMove(to-2*N, to, WP, Empty, Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}
}

// B pawns moves that doesn't capture aand not promotions
func (b *BoardStruct) GenBPawnNonCapt(ml *moves.MoveList) {
	var mv moves.Move
	bPawns := b.PieceBB[Pawn] & b.WbBB[BLACK]

	// one step
	to1Step := (bPawns >> (-S)) & ^b.AllBB()
	// two steps,
	to2Step := ((to1Step & Row6) >> (-S)) & ^b.AllBB()
	to1Step &= ^Row1

	// Add one step forward
	for to := to1Step.FirstOne(); to != 64; to = to1Step.FirstOne() {
		mv.PackMove(to-S, to, BP, Empty, Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}
	// Add two steps forward
	for to := to2Step.FirstOne(); to != 64; to = to2Step.FirstOne() {
		mv.PackMove(to-2*S, to, BP, Empty, Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}
}

// generates all pseudomoves
func (b *BoardStruct) GenAllMoves(ml *moves.MoveList) {
	b.GenPawnMoves(ml)
	b.GenKnightMoves(ml, ^b.WbBB[b.Stm])
	b.GenBishopMoves(ml, ^b.WbBB[b.Stm])
	b.GenRookMoves(ml, ^b.WbBB[b.Stm])
	b.GenQueenMoves(ml, ^b.WbBB[b.Stm])
	b.GenKingMoves(ml, ^b.WbBB[b.Stm])
}

func (b *BoardStruct) GenAllCaptures(ml *moves.MoveList) {
	oppBB := b.WbBB[b.Stm.Opposite()]
	b.GenPawnCapt(ml)
	b.GenKnightMoves(ml, oppBB)
	b.GenBishopMoves(ml, oppBB)
	b.GenRookMoves(ml, oppBB)
	b.GenQueenMoves(ml, oppBB)
	b.GenKingMoves(ml, oppBB)
}

// Create a list of captures from pawns to Kings (including promotions) - alternative
func (b *BoardStruct) GenAllCapturesy(ml *moves.MoveList) {
	us := b.Stm
	them := us.Opposite()
	usBB := b.WbBB[us]
	themBB := b.WbBB[them]
	allBB := b.AllBB()
	var atkBB, frBB bitboard.BitBoard
	var mv moves.Move

	// Pawns (including ep and promotions)
	b.GenPawnCapt(ml)

	// Knights
	pc := Pt2pc(Knight, us)
	frBB = b.PieceBB[Knight] & usBB
	for fr := frBB.FirstOne(); fr != 64; fr = frBB.FirstOne() {
		atkBB = AtksKnights[fr] & themBB
		for to := atkBB.FirstOne(); to != 64; to = atkBB.FirstOne() {
			cp := b.Squares[to]
			mv.PackMove(fr, to, pc, cp, Empty, b.Ep, b.Castlings)
			ml.Add(mv)
		}
	}

	// Bishops
	pc = Pt2pc(Bishop, us)
	frBB = b.PieceBB[Bishop] & usBB
	for fr := frBB.FirstOne(); fr != 64; fr = frBB.FirstOne() {
		atkBB = magic.MBishopTab[fr].Atks(allBB) & themBB
		for to := atkBB.FirstOne(); to != 64; to = atkBB.FirstOne() {
			cp := b.Squares[to]
			mv.PackMove(fr, to, pc, cp, Empty, b.Ep, b.Castlings)
			ml.Add(mv)
		}
	}

	// Rooks
	pc = Pt2pc(Rook, us)
	frBB = b.PieceBB[Rook] & usBB
	for fr := frBB.FirstOne(); fr != 64; fr = frBB.FirstOne() {
		atkBB = magic.MRookTab[fr].Atks(allBB) & themBB
		for to := atkBB.FirstOne(); to != 64; to = atkBB.FirstOne() {
			cp := b.Squares[to]
			mv.PackMove(fr, to, pc, cp, Empty, b.Ep, b.Castlings)
			ml.Add(mv)
		}
	}

	// Queens
	pc = Pt2pc(Queen, us)
	frBB = b.PieceBB[Queen] & usBB
	for fr := frBB.FirstOne(); fr != 64; fr = frBB.FirstOne() {
		atkBB = magic.MBishopTab[fr].Atks(allBB) & themBB
		atkBB |= magic.MRookTab[fr].Atks(allBB) & themBB
		for to := atkBB.FirstOne(); to != 64; to = atkBB.FirstOne() {
			cp := b.Squares[to]
			mv.PackMove(fr, to, pc, cp, Empty, b.Ep, b.Castlings)
			ml.Add(mv)
		}
	}

	// King
	pc = Pt2pc(King, us)
	fr := b.King[us]
	atkBB = AtksKings[fr] & themBB
	for to := atkBB.FirstOne(); to != 64; to = atkBB.FirstOne() {
		cp := b.Squares[to]
		mv.PackMove(fr, to, pc, cp, Empty, b.Ep, b.Castlings)
		ml.Add(mv)
	}
}

func (b *BoardStruct) GenAllNonCaptures(ml *moves.MoveList) {
	emptyBB := ^b.AllBB()
	b.GenPawnNonCapt(ml)
	b.GenKnightMoves(ml, emptyBB)
	b.GenBishopMoves(ml, emptyBB)
	b.GenRookMoves(ml, emptyBB)
	b.GenQueenMoves(ml, emptyBB)
	b.GenKingMoves(ml, emptyBB)
}

// generates all legal moves
func (b *BoardStruct) GenAllLegals(ml *moves.MoveList) {
	b.GenAllMoves(ml)
	b.FilterLegals(ml)
}

// generate all legal moves
func (b *BoardStruct) FilterLegals(ml *moves.MoveList) {
	for ix := len(*ml) - 1; ix >= 0; ix-- {
		mov := (*ml)[ix]
		if b.Move(mov) {
			b.Unmove(mov)
		} else {
			ml.Remove(ix)
		}
	}
}

func (b *BoardStruct) GenFrMoves(pc int, toBB bitboard.BitBoard, ml *moves.MoveList) {
}
