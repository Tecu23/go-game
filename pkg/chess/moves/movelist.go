package moves

type MoveList []Move

func (ml *MoveList) New(size int) {
	*ml = make(MoveList, 0, size)
}

func (ml *MoveList) Clear() {
	*ml = (*ml)[:0]
}

func (ml *MoveList) Add(mv Move) {
	*ml = append(*ml, mv)
}

func (ml *MoveList) Remove(ix int) {
	if len(*ml) > ix && ix >= 0 {
		*ml = append((*ml)[:ix], (*ml)[ix+1:]...)
	}
}

// Sort is sorting the moves in the Score/Move list according to the score per move
func (ml *MoveList) Sort() {
	bSwap := true
	for bSwap {
		bSwap = false
		for i := 0; i < len(*ml)-1; i++ {
			if (*ml)[i+1].Eval() > (*ml)[i].Eval() {
				(*ml)[i], (*ml)[i+1] = (*ml)[i+1], (*ml)[i]
				bSwap = true
			}
		}
	}
}

func (ml MoveList) String() string {
	theString := ""
	for _, mv := range ml {
		theString += mv.String() + " "
	}
	return theString
}
