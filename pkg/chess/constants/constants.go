// Package constants contains the varibles used throughout the project
package constants

import "github.com/Tecu23/go-game/pkg/chess/bitboard"

const (
	TotalPieces     = 12 // TotalPieces is the number of total pieces
	TotalSidePieces = 6  // TotalSidePieces is the number of pieces each side has
	WHITE           = Color(0)
	BLACK           = Color(1)

	Startpos = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - " // Game Starting position

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
