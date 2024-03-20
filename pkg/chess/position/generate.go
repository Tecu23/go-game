package position

import (
	"github.com/Tecu23/go-game/pkg/chess/bitboard"
	. "github.com/Tecu23/go-game/pkg/chess/constants"
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
		(*BoardStruct).wPawnAtksFr,
		(*BoardStruct).bPawnAtksFr,
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
func (b *BoardStruct) wPawnAtksFr(fr int) bitboard.BitBoard {
	frBB := bitboard.BitBoard(1) << uint(fr)

	// Attacks left and right
	toCap := ((frBB & ^FileA) << NW) & b.WbBB[BLACK]
	toCap |= ((frBB & ^FileH) << NE) & b.WbBB[BLACK]
	return toCap
}

// returns captures from fr-sq
func (b *BoardStruct) bPawnAtksFr(fr int) bitboard.BitBoard {
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
