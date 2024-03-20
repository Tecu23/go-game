package position

import (
	"fmt"
	"math/rand"

	. "github.com/Tecu23/go-game/pkg/chess/constants"
	"github.com/Tecu23/go-game/pkg/chess/moves"
)

// //////////////////////////////////////////////////////
// ////////////////////// HASH //////////////////////////
var (
	RandPcSq  [12 * 64]uint64 // keyvalues for 'pc on sq'
	RandEp    [8]uint64       // keyvalues for 8 ep files
	RandCastl [16]uint64      // keyvalues for castling states
)

// setup random generator with seed
var Rnd = (*rand.Rand)(rand.New(rand.NewSource(1013))) // usage: rnd.Intn(n) NOTE: n > 0

// Rand64 creates one 64 bit random number
func Rand64() uint64 {
	rand := uint64(0)

	for i := 0; i < 4; i++ {
		rand = uint64(int(rand<<16) | Rnd.Intn(1<<16))
	}

	return rand
}

// initKeys computes random hash keyvalues for pc/sq, ep and castlings
func InitKeys() {
	for i := 0; i < 12*64; i++ {
		RandPcSq[i] = Rand64()
	}
	for i := 0; i < 8; i++ {
		RandEp[i] = Rand64()
	}
	for i := 0; i < 16; i++ {
		RandCastl[i] = Rand64()
	}

	// check that all keys are different
	fmt.Println("checking all random keys")
	for pc := 0; pc < 12-1; pc++ {
		for sq := 0; sq < 64; sq++ {
			key1 := PcSqKey(pc, sq)
			for pc2 := pc + 1; pc2 < 12; pc2++ {
				for sq2 := 0; sq2 < 64; sq2++ {
					if key1 == PcSqKey(pc2, sq2) {
						fmt.Printf(
							"pc=%v, sq=%v gives the same key as pc=%v, sq=%v \n",
							pc,
							sq,
							pc2,
							sq2,
						)
					}
				}
			}
			for ep := A3; ep <= H3; ep++ {
				if key1 == EpKey(ep) {
					fmt.Printf("pc=%v, sq=%v gives the same key as ep=%v \n", pc, sq, ep)
				}
			}
			for c := uint(0); c < 16; c++ {
				if key1 == CastlKey(c) {
					fmt.Printf("pc=%v, sq=%v gives the same key as castl=%v \n", pc, sq, c)
				}
			}
		}
	}

	for ep := A3; ep < H3; ep++ {
		key1 := EpKey(ep)
		for ep2 := ep + 1; ep2 <= H3; ep2++ {
			if key1 == EpKey(ep2) {
				fmt.Printf("ep=%vgives the same key as ep=%v \n", ep, ep2)
			}
		}

		for c := uint(0); c < 16; c++ {
			if key1 == CastlKey(c) {
				fmt.Printf("ep=%v is the same key as castl=%v \n", ep, c)
			}
		}
	}

	for c := uint(0); c < 16-1; c++ {
		key1 := CastlKey(c)
		for c2 := c + 1; c2 < 16; c2++ {
			if key1 == CastlKey(c2) {
				fmt.Printf("castl=%v is the same key as castl=%v \n", c, c2)
			}
		}

	}
}

// for color we just flip with XOR ffffffffffffffff
// hash key after change color
func FlipSide(key uint64) uint64 {
	return ^key
}

// pcSqKey returns the keyvalue för piece on square
func PcSqKey(pc, sq int) uint64 {
	return RandPcSq[pc*64+sq]
}

// epKey returns the keyvalue for the current ep state
func EpKey(epSq int) uint64 {
	if epSq == 0 {
		return 0
	}
	return RandEp[epSq%8]
}

// castlKey returns the keyvalue for the current castling state
func CastlKey(castling uint) uint64 {
	return RandCastl[castling]
}

func CheckKey(b *BoardStruct) bool {
	key := uint64(0)
	for sq, pc := range b.Squares {
		if pc == Empty {
			continue
		}
		key ^= PcSqKey(pc, sq)
	}
	if b.Stm == BLACK {
		key = ^key
	}

	if key != b.Key {
		return false
	}

	return true
}

// //////////////////////////////////////////////////////
// ////////////////////// TRANS /////////////////////////
const EntrySize = 128 / 8

type TtEntry struct {
	Lock      uint32 // the lock, extra safety
	Move      uint32 // the best move from the search
	_         uint16 // alignement, not used
	Score     int16  // the score from the search
	Age       uint8  // the age of this entry
	Depth     int8   // the depth that the score is based on
	ScoreType uint8  // the score has this score type
	_         uint8  // alignement, not used
}

// clear one entry
func (e *TtEntry) Clear() {
	// Obs entry skall vara 16 bytes
	e.Lock = 0
	e.Move = uint32(moves.NoMove)
	// entry.utfyllnad = 0  behövs inte
	e.Score = 0
	e.Age = 0
	e.Depth = -1
	e.ScoreType = 0
	// entry.tomt = 0   behövs inte
}

type TranspStruct struct {
	Entries uint // number of entries
	Mask    uint // mask for the index
	CntUsed int  // The transposition table usage
	Age     int  // current age
	Tab     []TtEntry
	// for health tests
	CStores int
	CTried  int
	CFound  int
	CPrune  int
	CBest   int
}

var Trans TranspStruct

// allocate a new transposition table with the size from GUI
func (t *TranspStruct) New(mB int) error {
	if mB > 4000 {
		return fmt.Errorf("max transtable size is 4GB (~4000 MB)")
	}

	byteSize := mB << 20
	bits := SizeToBits(byteSize)

	t.Entries = 1 << uint(bits)
	t.Mask = t.Entries - 1

	t.Age = 0
	t.CntUsed = 0

	t.Tab = make([]TtEntry, t.Entries, t.Entries)
	t.Clear()
	return nil
}

// returns how many bits the mask will need to cover the table size
func SizeToBits(size int) uint {
	bits := uint(0)
	for cntEntries := size / EntrySize; cntEntries > 1; cntEntries /= 2 {
		bits++
	}

	return bits
}

// clear all entries, age and counters
func (t *TranspStruct) Clear() {
	var e TtEntry
	e.Clear()

	for i := uint(0); i < t.Entries; i++ {
		t.Tab[i] = e
	}

	t.Age = 0
	t.CntUsed = 0

	// counts
	t.CFound, t.CStores, t.CTried, t.CPrune, t.CBest = 0, 0, 0, 0, 0
}

// index uses the Key to compute an index into the table
func (t *TranspStruct) Index(fullKey uint64) int64 {
	return int64(fullKey)
}

// Lock extracts the lock value from the hash key
func (t *TranspStruct) Lock(fullKey uint64) uint32 {
	return uint32(fullKey >> 32)
}

func (t *TranspStruct) InitSearch() {
	t.IncAge()
	t.CntUsed = 0

	// Health check counts
	t.CFound, t.CStores, t.CTried, t.CPrune, t.CBest = 0, 0, 0, 0, 0
}

// incAge increments the date for the hahs table.
// We are reborned after the age 255
func (t *TranspStruct) IncAge() {
	t.Age = (t.Age + 1) % 256
}

func (b *BoardStruct) FullKey() uint64 {
	key := b.Key ^ EpKey(b.Ep)
	key ^= CastlKey(uint(b.Castlings))
	return key
}

// store current position in the transp table.
// The key is computed from the position. The lock value is the 32 first bits in the key
// From the key we get an index to the table.
// We will try 4 entries in a sequence if a lock is found
// We always try to replace another age and/or a lower searched depth

func (t *TranspStruct) Store(fullKey uint64, mv moves.Move, depth, ply, sc, scoreType int) {
	t.CStores++
	sc = RemoveMatePly(sc, ply)

	index := fullKey & uint64(t.Mask)
	lock := t.Lock(fullKey)

	var newEntry *TtEntry
	bestDep := -2000

	for i := uint64(0); i < 4; i++ {
		idx := (uint64(index) + i) & uint64(t.Mask)

		entry := &t.Tab[idx]

		if entry.Lock == lock {
			if int(entry.Age) != t.Age {
				entry.Age = uint8(t.Age)
				t.CntUsed++
			}

			if depth >= int(entry.Depth) {
				if mv != moves.NoMove {
					entry.Move = uint32(mv.OnlyMv())
				}
				entry.Depth = int8(depth)
				entry.Score = int16(sc)
				entry.ScoreType = uint8(scoreType)
				return
			}

			if entry.Move == uint32(moves.NoMove) {
				entry.Move = uint32(mv.OnlyMv())
			}

			return
		}

		adjDepth := -int(entry.Depth)
		if entry.Age != uint8(t.Age) {
			adjDepth += 1000 // make age more important than depth
		}

		if adjDepth > bestDep {
			newEntry = entry
			bestDep = adjDepth
		}
	}

	if newEntry.Age != uint8(t.Age) {
		t.CntUsed++
	}

	newEntry.Lock = lock
	newEntry.Age = uint8(t.Age)
	newEntry.Depth = int8(depth)
	newEntry.Move = uint32(mv)
	newEntry.Score = int16(sc)
	newEntry.ScoreType = uint8(scoreType)
}

// retrieve get move and score to the current position from the transp Table if the key and lock is correct
// if no entry is matching return false else return true, depth not ok return false but with move filled in
// We will try the 4 entries in sequence until lock match otherwise return false
func (t *TranspStruct) Retrieve(
	fullKey uint64,
	depth, ply int,
) (mv moves.Move, sc, scoreType int, ok bool) {
	t.CTried++
	mv = moves.NoMove
	ok = false
	sc = NoScore
	scoreType = 0

	index := fullKey & uint64(t.Mask)
	lock := t.Lock(fullKey)

	for i := uint64(0); i < 4; i++ {

		idx := uint(index+i) & t.Mask
		entry := &t.Tab[idx]

		if entry.Lock == lock { // there is a matching position already here
			t.CFound++

			if int(entry.Age) != t.Age { // from another generation?
				entry.Age = uint8(t.Age)
				t.CntUsed++
			}
			mv = moves.Move(entry.Move)
			sc = AddMatePly(int(entry.Score), ply)
			scoreType = int(entry.ScoreType)
			ok = true
			if int(entry.Depth) >= depth {
				//				fmt.Println("retrieve (key,depth,ply)",fullKey,depth,ply,"(mv,sc,scTyp):", mv, sc, scoreType)
				return
			}

			if IsMateScore(sc) {
				scoreType &= ^ScoreTypeUpper
				if sc < 0 {
					scoreType &= ^ScoreTypeLower
				}
				//				fmt.Println("retrieve (key,depth,ply)",fullKey,depth,ply,"(mv,sc,scTyp):", mv, sc, scoreType)
				return
			}
			//			fmt.Println("Nothing to retrieve-depth? (key,depth,ply):", fullKey, depth,ply)
			ok = false
			return
		}
	}
	//	fmt.Println("Nothing to retrieve - no hit (key,depth,ply):", fullKey, depth,ply)
	ok = false
	return
}

// isMateScore returns true if the score is a mate score
func IsMateScore(sc int) bool {
	return sc < MinEval+MaxPly || sc > MaxEval-MaxPly
}

// removeMatePly removes ply from the score value (score - ply) if mate
// in order to mix up different depths
func RemoveMatePly(sc, Ply int) int {
	if sc < MinEval+MaxPly {
		return -MateEval
	} else if sc > MaxEval-MaxPly {
		return MateEval
	}
	return sc
}

// addMatePly adjusts mate value with ply if mate score
func AddMatePly(sc, ply int) int {
	if sc < MinEval+MaxPly {
		return -MateEval + ply
	} else if sc > MaxEval-MaxPly {
		return MateEval - ply
	}
	return sc
}

// scoreType sets if it is an upper or lower score
func ScoreType(sc, alpha, beta int) int {
	scoreType := 0
	if sc > alpha {
		scoreType |= ScoreTypeLower
	}
	if sc < beta {
		scoreType |= ScoreTypeUpper
	}

	return scoreType
}
