package engine

import . "github.com/Tecu23/go-game/pkg/chess/constants"

type EbfStruct []uint64

func (e *EbfStruct) New() {
	*e = make(EbfStruct, 0, MaxDepth)
}

func (e *EbfStruct) Add(nodes uint64) {
	*e = append(*e, nodes)
}

func (e *EbfStruct) Clear() {
	*e = (*e)[:0]
}

func (e *EbfStruct) Ebf(depth int) float64 {
	if len(*e) < 4 {
		return 0
	}
	ebf := 0.0
	nodes := float64((*e)[len(*e)-1]) // the last depth
	prevNodes := float64((*e)[len(*e)-2])

	if nodes > 0.0 && prevNodes > 0.0 {
		//		ebf = (prevNodes2/prevNodes3 + prevNodes1/prevNodes2) / 2
		ebf = nodes / prevNodes
	}
	return ebf
}
