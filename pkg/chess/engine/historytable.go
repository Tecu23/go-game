package engine

// /////////////////////////// history table //////////////////////////////////
import (
	"fmt"

	. "github.com/Tecu23/go-game/pkg/chess/constants"
)

type HistoryStruct [2][64][64]uint

func (h *HistoryStruct) Inc(fr, to int, stm Color, depth int) {
	h[stm][fr][to] += uint(depth * depth)
}

func (h *HistoryStruct) Get(fr, to int, stm Color) uint {
	return h[stm][fr][to]
}

func (h *HistoryStruct) Clear() {
	for fr := 0; fr < 64; fr++ {
		for to := 0; to < 64; to++ {
			h[0][fr][to] = 0
			h[1][fr][to] = 0
		}
	}
}

func (h HistoryStruct) Print(n int) {
	fmt.Println("history top", n)
	type top50 struct{ fr, to, sd, sc uint }
	hTab := make([]top50, n, n)
	for ix := range hTab {
		hTab[ix].fr, hTab[ix].to, hTab[ix].sd, hTab[ix].sc = 0, 0, 0, 0
	}

	W, B := uint(WHITE), uint(BLACK)
	for fr := uint(0); fr < 64; fr++ {
		for to := uint(0); to < 64; to++ {
			sc := h.Get(int(fr), int(to), WHITE)
			for ix := range hTab {
				if sc > hTab[ix].sc {
					for ix2 := n - 2; ix2 >= ix; ix2-- {
						hTab[ix2+1] = hTab[ix2]
					}
					hTab[ix].fr, hTab[ix].to, hTab[ix].sd, hTab[ix].sc = fr, to, W, sc
					break
				}
			}

			sc = h.Get(int(fr), int(to), BLACK)
			for ix := range hTab {
				if sc > hTab[ix].sc {
					for ix2 := n - 2; ix2 >= ix; ix2-- {
						hTab[ix2+1] = hTab[ix2]
					}
					hTab[ix].fr, hTab[ix].to, hTab[ix].sd, hTab[ix].sc = fr, to, B, sc
					break
				}
			}
		}
	}
	for ix, ht := range hTab {
		if ht.fr == 0 && ht.to == 0 {
			continue
		}
		fmt.Printf(
			"%2v: %v %v-%v   %v  \n",
			ix+1,
			Color(ht.sd).String(),
			Sq2Fen[int(ht.fr)],
			Sq2Fen[int(ht.to)],
			ht.sc,
		)
	}
}

var History HistoryStruct
