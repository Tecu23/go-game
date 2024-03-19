// Package bitboard implements utility functions
// for using bitboards
package bitboard

import (
	"fmt"
	"math/bits"
	"strings"
)

// BitBoard reprents a 64 bits unsigned int
type BitBoard uint64

// Count should count the number of 1s in a bitboard
func (b BitBoard) Count() int {
	return bits.OnesCount64(uint64(b))
}

// SetBit should change a bit on the bitboard to 1 at a certain position
func (b *BitBoard) SetBit(pos int) {
	*b |= BitBoard(uint64(1) << uint(pos))
}

// IsBitSet should check if the bit at position pos is set
func (b BitBoard) IsBitSet(pos int) bool {
	return (b & BitBoard(uint64(1)<<uint(pos))) != 0
}

// Clear should set the bit at position pos to 0
func (b *BitBoard) Clear(pos int) {
	*b &= BitBoard(^(uint64(1) << uint(pos)))
}

// FirstOne should remove the trailing zero and returns the position
func (b *BitBoard) FirstOne() int {
	bit := bits.TrailingZeros64(uint64(*b))
	if bit == 64 { // the method returns 64 when the bb is 0
		return 64
	}
	*b = (*b >> uint(bit+1)) << uint(bit+1)
	return bit
}

// LastOne removes the leading zero and returns the position
func (b *BitBoard) LastOne() int {
	bit := bits.LeadingZeros64(uint64(*b))
	if bit == 64 { // the method returns 64 when the bb is 0
		return 64
	}
	*b = (*b << uint(bit+1)) >> uint(bit+1)
	return 63 - bit
}

// ToString returns the full bitstring (with leading zeroes) of the BitBoard
func (b BitBoard) ToString() string {
	zeroes := ""
	for ix := 0; ix < 64; ix++ {
		zeroes = zeroes + "0"
	}

	bits := zeroes + fmt.Sprintf("%b", b)
	return bits[len(bits)-64:]
}

// Stringln returns the bitboard string 8x8
func (b BitBoard) Stringln() string {
	s := b.ToString()
	row := [8]string{}
	row[0] = s[0:8]
	row[1] = s[8:16]
	row[2] = s[16:24]
	row[3] = s[24:32]
	row[4] = s[32:40]
	row[5] = s[40:48]
	row[6] = s[48:56]
	row[7] = s[56:]
	for ix, r := range row {
		row[ix] = fmt.Sprintf(
			"%v%v%v%v%v%v%v%v\n",
			r[7:8],
			r[6:7],
			r[5:6],
			r[4:5],
			r[3:4],
			r[2:3],
			r[1:2],
			r[0:1],
		)
	}

	s = strings.Join(row[:], "")
	s = strings.Replace(s, "1", "1 ", -1)
	s = strings.Replace(s, "0", "0 ", -1)
	return s
}
