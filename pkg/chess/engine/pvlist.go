package engine

import (
	. "github.com/Tecu23/go-game/pkg/chess/constants"
	"github.com/Tecu23/go-game/pkg/chess/moves"
)

type PvList []moves.Move

func (pv *PvList) New() {
	*pv = make(PvList, 0, MaxPly)
}

func (pv *PvList) Add(mv moves.Move) {
	*pv = append(*pv, mv)
}

func (pv *PvList) Clear() {
	*pv = (*pv)[:0]
}

func (pv *PvList) AddPV(pv2 *PvList) {
	*pv = append(*pv, *pv2...)
}

func (pv *PvList) Catenate(mv moves.Move, pv2 *PvList) {
	pv.Clear()
	pv.Add(mv)
	pv.AddPV(pv2)
}

func (pv *PvList) String() string {
	s := ""
	for _, mv := range *pv {
		s += mv.String() + " "
	}
	return s
}
