package moves

import (
	"fmt"
	"strings"

	"github.com/Tecu23/go-game/pkg/chess/castlings"
	. "github.com/Tecu23/go-game/pkg/chess/constants"
)

type Move uint64

const NoMove = Move(0)

func (m *Move) PackMove(fr, to, pc, cp, pr, epSq int, castl castlings.Castlings) {
	// 6 bits fr, 6 bits to, 4 bits pc, 4 bits cp, 4 bits prom, 4 bits ep, 4 bits castl = 32 bits

	if epSq == Empty {
		// Handle this somehow
	}

	epFile := 0

	if epSq != 0 {
		epFile = epSq%8 + 1
	}

	*m = Move(fr | (to << ToShift) | (pc << PcShift) |
		(cp << CpShift) | (pr << PrShift) |
		(epFile << EpShift) | int(castl<<CastlShift))
}

func (m *Move) PackEval(score int) {
	(*m) &= Move(^EvalMask) // clear eval
	(*m) |= Move(score+30000) << EvalShift
}

// compare two moves - only frSq and toSq
func (m Move) CmpFrTo(m2 Move) bool {
	// return (m & move(^evalMask)) == (m2 & move(^evalMask))
	return m.Fr() == m2.Fr() && m.To() == m2.To()
}

// compare two moves - only frSq, toSq and pc
func (m Move) CmpFrToP(m2 Move) bool {
	// return (m & move(^evalMask)) == (m2 & move(^evalMask))
	return m.Fr() == m2.Fr() && m.To() == m2.To() && m.Pc() == m2.Pc()
}

func (m Move) Cmp(m2 Move) bool {
	return (m & Move(^EvalMask)) == (m2 & Move(^EvalMask))
}

func (m Move) Eval() int {
	return int((uint(m)&uint(EvalMask))>>EvalShift) - 30000
}

func (m Move) Fr() int {
	return int(m & FrMask)
}

func (m Move) To() int {
	return int(m&ToMask) >> ToShift
}

func (m Move) Pc() int {
	return int(m&PcMask) >> PcShift
}

func (m Move) Cp() int {
	return int(m&CpMask) >> CpShift
}

func (m Move) Pr() int {
	return int(m&PrMask) >> PrShift
}

func (m Move) Castl() castlings.Castlings {
	return castlings.Castlings(m&CastlMask) >> CastlShift
}

func (m Move) Ep(sd Color) int {
	// sd is the side that can capture
	file := int(m&EpMask) >> EpShift
	if file == 0 {
		return 0 // no ep
	}

	// there is an ep sq
	rank := 5
	if sd == BLACK {
		rank = 2
	}

	return rank*8 + file - 1
}

// move without eval
func (m Move) OnlyMv() Move {
	return m & Move(^EvalMask)
}

func (m Move) String() string {
	s := m.StringFull()
	s = s[1:3] + s[5:]
	return s
}

func (m Move) StringFull() string {
	fr := Sq2Fen[int(m.Fr())]
	to := Sq2Fen[int(m.To())]
	pc := Pc2Fen(int(m.Pc()))
	cp := Pc2Fen(int(m.Cp())) + " "
	pr := Pc2Fen(int(m.Pr()))
	return strings.TrimSpace(fmt.Sprintf("%v%v-%v%v%v", pc, fr, cp[:1], to, pr))
}

var PieceRules [NoPiecesT][]int // not pawns
