package magic

import (
	"fmt"

	"github.com/Tecu23/go-game/pkg/chess/bitboard"
	. "github.com/Tecu23/go-game/pkg/chess/constants"
)

type sMagic struct {
	ToSqBB  []bitboard.BitBoard // all possible atkBoards dep on blockers
	InnerBB bitboard.BitBoard   // atks on empty board
	Magic   uint64
	Shift   uint
}

var (
	MBishopTab [64]sMagic
	MRookTab   [64]sMagic
)

// all attacks from current square
func (m *sMagic) Atks(allBB bitboard.BitBoard) bitboard.BitBoard {
	return m.ToSqBB[int(((allBB&m.InnerBB)*bitboard.BitBoard(m.Magic))>>m.Shift)]
}

func InitMagic() {
	fmt.Println("starting init() for magic.go")

	// bishops
	// fillOptimalMagicsB()
	for sq := A1; sq <= H8; sq++ {
		MBishopTab[sq].Shift = uint(64 - NBBits[sq])
		MBishopTab[sq].InnerBB = bitboard.BitBoard(InnerBAtks(sq))

		MBishopTab[sq].Magic = MagicB[sq]
	}

	// rooks
	FillOptimalMagicsR()
	for sq := A1; sq <= H8; sq++ {
		MRookTab[sq].Shift = uint(64 - NRBits[sq])
		MRookTab[sq].InnerBB = bitboard.BitBoard(InnerRAtks(sq))
		MRookTab[sq].Magic = MagicR[sq]
	}

	PrepareMagicB() // Bishops
	PrepareMagicR() // Rooks
}

// var toSqBB *[]bitboard.BitBoard // pointer to mRookTab[sq].toSqBB or mBishopTab[sq].toSqBB
type dirstr struct {
	rDir int
	fDir int
}

// create move bitboard.BitBoards for bishops on all squares
func PrepareMagicB() {
	dirsB := []dirstr{{+1, +1}, {-1, +1}, {+1, -1}, {-1, -1}}
	for fr := A1; fr <= H8; fr++ {
		maxM := -1
		// all bit combinations for fr and all possible blockers
		cnt := BitCombs(0x0, fr, fr, 0, &maxM, &MBishopTab[fr], dirsB)
		_ = cnt
		//		fmt.Println("bishop on", sq2Fen[fr], "#of combinations", cnt, "maxIx", maxM)
	}
}

// create move bitboard.BitBoards for rooks on all squares
func PrepareMagicR() {
	dirsR := []dirstr{{+1, 0}, {-1, 0}, {0, +1}, {0, -1}}
	for fr := A1; fr <= H8; fr++ {
		maxM := -1
		// all bit combinations for fr and all possible moves (toSqBB)
		cnt := BitCombs(0x0, fr, fr, 0, &maxM, &MRookTab[fr], dirsR)
		_ = cnt
		//		fmt.Println("rook on", sq2Fen[fr], "#of combinations", cnt, "maxIx", maxM)
	}
}

// find all combinations of blockers
func BitCombs(
	wBits bitboard.BitBoard,
	fr, currSq, currIx int,
	maxM *int,
	mTabEntry *sMagic,
	dirs []dirstr,
) int {
	magic := mTabEntry.Magic
	shift := mTabEntry.Shift
	innerBB := uint64(mTabEntry.InnerBB)
	toSqBB := &(*mTabEntry).ToSqBB
	cnt := 0

	currSq = GetNextSq(fr, currSq, &currIx, dirs)
	if currSq == -1 { // vägs ände
		// append new toSqBB
		m := (uint64(wBits) & innerBB) * magic
		m = m >> shift
		if int(m) > *maxM {
			// fill upp with empty entries
			for ; *maxM < int(m); *maxM++ {
				*toSqBB = append(*toSqBB, 0x0)
			}
		}
		toBB := bitboard.BitBoard(ComputeAtks(fr, dirs, uint64(wBits)))

		if (*toSqBB)[int(m)] != 0x0 && (*toSqBB)[int(m)] != toBB { // for bishop
			fmt.Println(
				"we have problem",
				Sq2Fen[fr],
				"with ix",
				int(m),
				"wBits:\n",
				bitboard.BitBoard(wBits).Stringln(),
			)
			fmt.Println((*toSqBB)[int(m)].Stringln())
			fmt.Printf("%X\n", uint64((*toSqBB)[int(m)]))
			fmt.Println(toBB.Stringln())
		}
		(*toSqBB)[int(m)] = toBB
		return 1
	}

	// 1
	//
	//						wBits |= (uint64(1) << uint(currSq))
	wBits.SetBit(currSq)
	cnt += BitCombs(wBits, fr, currSq, currIx, maxM, mTabEntry, dirs)

	// 0
	//
	//						wBits &= ^(uint64(1) << uint(currSq))
	wBits.Clear(currSq)
	cnt += BitCombs(wBits, fr, currSq, currIx, maxM, mTabEntry, dirs)

	return cnt
}

func GetNextSq(fr, currSq int, currIx *int, dirs []dirstr) int {
	for ; *currIx < 4; *currIx++ {
		rk := currSq / 8
		fl := currSq % 8
		r := dirs[*currIx].rDir
		f := dirs[*currIx].fDir
		rk += r
		fl += f
		if (r == 0 || (rk > 0 && rk < 7)) && (f == 0 || (fl > 0 && fl < 7)) {
			return rk*8 + fl
		}
		currSq = fr
	}

	return -1
}

func ComputeAtks(fr int, dirs []dirstr, toBB uint64) uint64 {
	movesBB := uint64(0)

	// 0
	r := dirs[0].rDir
	f := dirs[0].fDir
	rk := fr/8 + r
	fl := fr%8 + f
	for rk >= 0 && fl >= 0 && rk < 8 && fl < 8 {
		sq := uint(rk*8 + fl)
		sqBit := uint64(1) << sq

		movesBB |= sqBit
		if sqBit&toBB != 0 {
			break
		}
		rk += r
		fl += f
	}

	// 1
	r = dirs[1].rDir
	f = dirs[1].fDir
	rk = fr/8 + r
	fl = fr%8 + f
	for rk >= 0 && fl >= 0 && rk < 8 && fl < 8 {
		sq := uint(rk*8 + fl)
		sqBit := uint64(1) << sq

		movesBB |= sqBit
		if sqBit&toBB != 0 {
			break
		}
		rk += r
		fl += f
	}

	// 2
	r = dirs[2].rDir
	f = dirs[2].fDir
	rk = fr/8 + r
	fl = fr%8 + f
	for rk >= 0 && fl >= 0 && rk < 8 && fl < 8 {
		sq := uint(rk*8 + fl)
		sqBit := uint64(1) << sq

		movesBB |= sqBit
		if sqBit&toBB != 0 {
			break
		}
		rk += r
		fl += f
	}

	// 3
	r = dirs[3].rDir
	f = dirs[3].fDir
	rk = fr/8 + r
	fl = fr%8 + f
	for rk >= 0 && fl >= 0 && rk < 8 && fl < 8 {
		sq := uint(rk*8 + fl)
		sqBit := uint64(1) << sq

		movesBB |= sqBit
		if sqBit&toBB != 0 {
			break
		}
		rk += r
		fl += f
	}

	return movesBB
}

// all bishop inner attacks from sq on an empty board
func InnerBAtks(sq int) uint64 {
	atkBB := uint64(0)
	// NE (+9)
	rw := sq / 8
	fl := sq % 8
	r := rw + 1
	f := fl + 1
	for r < 7 && f < 7 {
		atkBB |= uint64(1) << uint(r*8+f)
		r++
		f++
	}

	// NW (+7)
	r = rw + 1
	f = fl - 1
	for r < 7 && f > 0 {
		atkBB |= uint64(1) << uint(r*8+f)
		r++
		f--
	}
	// SW (-7)
	r = rw - 1
	f = fl - 1
	for r > 0 && f > 0 {
		atkBB |= uint64(1) << uint(r*8+f)
		r--
		f--
	}

	// SE (-9)
	r = rw - 1
	f = fl + 1
	for r > 0 && f < 7 {
		atkBB |= uint64(1) << uint(r*8+f)
		r--
		f++
	}
	return atkBB
}

// all rook attacks from sq on an empty board
func InnerRAtks(sq int) uint64 {
	atkBB := uint64(0)
	// N (+8)
	rw := sq / 8
	fl := sq % 8
	r := rw + 1
	f := fl
	for r < 7 {
		atkBB |= uint64(1) << uint(r*8+f)
		r++
	}

	// E (+1)
	r = rw
	f = fl + 1
	for f < 7 {
		atkBB |= uint64(1) << uint(r*8+f)
		f++
	}
	// S (-8)
	r = rw - 1
	f = fl
	for r > 0 {
		atkBB |= uint64(1) << uint(r*8+f)
		r--
	}

	// W (-1)
	r = rw
	f = fl - 1
	for f > 0 {
		atkBB |= uint64(1) << uint(r*8+f)
		f--
	}

	return atkBB
}

func FillOptimalMagicsB() {
	NBBits[A1] = 5
	MagicB[A1] = 0xffedf9fd7cfcffff
	NBBits[B1] = 4
	MagicB[B1] = 0xfc0962854a77f576
	NBBits[C1] = 5
	MagicB[C1] = 0xE433BF9FF9BD3C0D
	NBBits[D1] = 5
	MagicB[D1] = 0x8F0BBE9CF98C0405
	NBBits[E1] = 5
	MagicB[E1] = 0x7E11DFD9DDFBDBF0
	NBBits[G1] = 4
	MagicB[G1] = 0xfc0a66c64a7ef576
	NBBits[H1] = 5
	MagicB[H1] = 0x7ffdfdfcbd79ffff
	NBBits[A2] = 4
	MagicB[A2] = 0xfc0846a64a34fff6
	NBBits[B2] = 4
	MagicB[B2] = 0xfc087a874a3cf7f6
	NBBits[C2] = 5
	MagicB[C2] = 0x0040020042188680
	NBBits[D2] = 5
	MagicB[D2] = 0x0080000108D80200
	NBBits[E2] = 5
	MagicB[E2] = 0xF2048D48B0240820
	NBBits[F2] = 5
	MagicB[F2] = 0x810040B921030010
	NBBits[G2] = 4
	MagicB[G2] = 0xfc0864ae59b4ff76
	NBBits[H2] = 4
	MagicB[H2] = 0x3c0860af4b35ff76
	NBBits[A3] = 4
	MagicB[A3] = 0x73C01AF56CF4CFFB
	NBBits[B3] = 4
	MagicB[B3] = 0x41A01CFAD64AAFFC
	NBBits[G3] = 4
	MagicB[G3] = 0x7c0c028f5b34ff76
	NBBits[H3] = 4
	MagicB[H3] = 0xfc0a028e5ab4df76
	NBBits[A6] = 4
	MagicB[A6] = 0xDCEFD9B54BFCC09F
	NBBits[B6] = 4
	MagicB[B6] = 0xF95FFA765AFD602B
	NBBits[G6] = 4
	MagicB[G6] = 0x43ff9a5cf4ca0c01
	NBBits[H6] = 4
	MagicB[H6] = 0x4BFFCD8E7C587601
	NBBits[A7] = 4
	MagicB[A7] = 0xfc0ff2865334f576
	NBBits[B7] = 4
	MagicB[B7] = 0xfc0bf6ce5924f576
	NBBits[G7] = 4
	MagicB[G7] = 0xc3ffb7dc36ca8c89
	NBBits[H7] = 4
	MagicB[H7] = 0xc3ff8a54f4ca2c89
	NBBits[A8] = 5
	MagicB[A8] = 0xfffffcfcfd79edff
	NBBits[B8] = 4
	MagicB[B8] = 0xfc0863fccb147576
	NBBits[G8] = 4
	MagicB[G8] = 0xfc087e8e4bb2f736
	NBBits[H8] = 5
	MagicB[H8] = 0x43ff9e4ef4ca2c89
}

func FillOptimalMagicsR() {
	NRBits[A7] = 10
	MagicR[A7] = 0x48FFFE99FECFAA00
	NRBits[B7] = 9
	MagicR[B7] = 0x48FFFE99FECFAA00
	NRBits[C7] = 9
	MagicR[C7] = 0x497FFFADFF9C2E00
	NRBits[D7] = 9
	MagicR[D7] = 0x613FFFDDFFCE9200
	NRBits[E7] = 9
	MagicR[E7] = 0xffffffe9ffe7ce00
	NRBits[F7] = 9
	MagicR[F7] = 0xfffffff5fff3e600
	NRBits[G7] = 9
	MagicR[G7] = 0x3ff95e5e6a4c0
	NRBits[H7] = 10
	MagicR[H7] = 0x510FFFF5F63C96A0
	NRBits[A8] = 11
	MagicR[A8] = 0xEBFFFFB9FF9FC526
	NRBits[B8] = 10
	MagicR[B8] = 0x61FFFEDDFEEDAEAE
	NRBits[C8] = 10
	MagicR[C8] = 0x53BFFFEDFFDEB1A2
	NRBits[D8] = 10
	MagicR[D8] = 0x127FFFB9FFDFB5F6
	NRBits[E8] = 10
	MagicR[E8] = 0x411FFFDDFFDBF4D6
	NRBits[G8] = 10
	MagicR[G8] = 0x0003ffef27eebe74
	NRBits[H8] = 11
	MagicR[H8] = 0x7645FFFECBFEA79E
}

var NRBits = [64]int{
	12, 11, 11, 11, 11, 11, 11, 12,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	11, 10, 10, 10, 10, 10, 10, 11,
	12, 11, 11, 11, 11, 11, 11, 12,
}

var NBBits = [64]int{
	6, 5, 5, 5, 5, 5, 5, 6,
	5, 5, 5, 5, 5, 5, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 9, 9, 7, 5, 5,
	5, 5, 7, 7, 7, 7, 5, 5,
	5, 5, 5, 5, 5, 5, 5, 5,
	6, 5, 5, 5, 5, 5, 5, 6,
}

// Thanks to Tord Romstad code from:
// https://chessprogramming.wikispaces.com/Looking+for+Magics
// This site will soon be closed down and moved somewhere else
// Search for "Feeding in randoms Tord Romstad"

var MagicB = [64]uint64{
	0xc085080200420200,
	0x60014902028010,
	0x401240100c201,
	0x580ca104020080,
	0x8434052000230010,
	0x102080208820420,
	0x2188410410403024,
	0x40120805282800,
	0x4420410888208083,
	0x1049494040560,
	0x6090100400842200,
	0x1000090405002001,
	0x48044030808c409,
	0x20802080384,
	0x2012008401084008,
	0x9741088200826030,
	0x822000400204c100,
	0x14806004248220,
	0x30200101020090,
	0x148150082004004,
	0x6020402112104,
	0x4001000290080d22,
	0x2029100900400,
	0x804203145080880,
	0x60a10048020440,
	0xc08080b20028081,
	0x1009001420c0410,
	0x101004004040002,
	0x1004405014000,
	0x10029a0021005200,
	0x4002308000480800,
	0x301025015004800,
	0x2402304004108200,
	0x480110c802220800,
	0x2004482801300741,
	0x400400820a60200,
	0x410040040040,
	0x2828080020011000,
	0x4008020050040110,
	0x8202022026220089,
	0x204092050200808,
	0x404010802400812,
	0x422002088009040,
	0x180604202002020,
	0x400109008200,
	0x2420042000104,
	0x40902089c008208,
	0x4001021400420100,
	0x484410082009,
	0x2002051108125200,
	0x22e4044108050,
	0x800020880042,
	0xb2020010021204a4,
	0x2442100200802d,
	0x10100401c4040000,
	0x2004a48200c828,
	0x9090082014000,
	0x800008088011040,
	0x4000000a0900b808,
	0x900420000420208,
	0x4040104104,
	0x120208c190820080,
	0x4000102042040840,
	0x8002421001010100,
}

var MagicR = [64]uint64{
	0x11800040001481a0,
	0x2040400010002000,
	0xa280200308801000,
	0x100082005021000,
	0x280280080040006,
	0x200080104100200,
	0xc00040221100088,
	0xe00072200408c01,
	0x2002045008600,
	0xa410804000200089,
	0x4081002000401102,
	0x2000c20420010,
	0x800800400080080,
	0x40060010041a0009,
	0x441004442000100,
	0x462800080004900,
	0x80004020004001,
	0x1840420021021081,
	0x8020004010004800,
	0x940220008420010,
	0x2210808008000400,
	0x24808002000400,
	0x803604001019a802,
	0x520000440081,
	0x802080004000,
	0x1200810500400024,
	0x8000100080802000,
	0x2008080080100480,
	0x8000404002040,
	0xc012040801104020,
	0xc015000900240200,
	0x20040200208041,
	0x1080004000802080,
	0x400081002110,
	0x30002000808010,
	0x2000100080800800,
	0x2c0800400800800,
	0x1004800400800200,
	0x818804000210,
	0x340082000a45,
	0x8520400020818000,
	0x2008900460020,
	0x100020008080,
	0x601001000a30009,
	0xc001000408010010,
	0x2040002008080,
	0x11008218018c0030,
	0x20c0080620011,
	0x400080002080,
	0x8810040002500,
	0x400801000200080,
	0x2402801000080480,
	0x204040280080080,
	0x31044090200801,
	0x40c10830020400,
	0x442800100004080,
	0x10080002d005041,
	0x134302820010a2c2,
	0x6202001080200842,
	0x1820041000210009,
	0x1002001008210402,
	0x2000108100402,
	0x10310090a00b824,
	0x800040100944822,
}

/*
func findDuplicates(toBBTab []bitboard.BitBoard, fr int, magic uint64) int {
	cntDup := 0
	if len(toBBTab) < 2 {
		return 0
	}
	for ix, toBB := range toBBTab[:len(toBBTab)-1] {
		dup := false
		for _, toBB2 := range toBBTab[ix+1:] {
			if toBB == toBB2 && toBB != 0x0 {
				dup = true
			}
		}
		if dup {
			cntDup++
		}
	}
	return cntDup
}
*/
