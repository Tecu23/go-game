// Package engine contains the main logic behind the chess engine
package engine

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/Tecu23/go-game/pkg/chess/bitboard"
	. "github.com/Tecu23/go-game/pkg/chess/constants"
	"github.com/Tecu23/go-game/pkg/chess/magic"
	"github.com/Tecu23/go-game/pkg/chess/moves"
	"github.com/Tecu23/go-game/pkg/chess/position"
)

var CntNodes uint64

// TODO search limits: count nodes and test for limit.nodes
// TODO search limits: limit.depth

// TODO search limits: time per game w/wo increments
// TODO search limits: time per x moves and after x moves w/wo increments
type SearchLimits struct {
	Depth     int
	Nodes     uint64
	MoveTime  int // in milliseconds
	Infinite  bool
	StartTime time.Time
	LastTime  time.Time

	// Current
	Stop bool
}

// Limits are the engine settings set by the user
var Limits SearchLimits

func (s *SearchLimits) Init() {
	s.Depth = 9999
	s.Nodes = math.MaxUint64
	s.MoveTime = 99999999999
	s.Infinite = false
	s.Stop = false
}

func (s *SearchLimits) SetStop(st bool) {
	s.Stop = st
}

func (s *SearchLimits) SetDepth(d int) {
	s.Depth = d
}

func (s *SearchLimits) SetMoveTime(m int) {
	s.MoveTime = m
}

func (s *SearchLimits) SetInfinite(b bool) {
	s.Infinite = b
}

// Engine should create the 2 channels necessary to communicate to the websocket
func Engine() (chan bool, chan string) {
	frEngine := make(chan string)
	toEngine := make(chan bool)

	go root(toEngine, frEngine)

	return toEngine, frEngine
}

func root(toEngine chan bool, frEngine chan string) {
	var depth, alpha, beta int
	var ebfTab EbfStruct
	var pv PvList
	var childPV PvList
	var ml moves.MoveList
	childPV.New()
	pv.New()
	ml.New(60)
	ebfTab.New()
	b := &position.Board
	for range toEngine {
		Limits.StartTime, Limits.LastTime = time.Now(), time.Now()
		CntNodes = 0
		ebfTab.Clear()
		Killers.Clear()
		ml.Clear()
		pv.Clear()

		position.Trans.InitSearch() // incr age coounters=0

		genAndSort(0, b, &ml)
		depth = 0

		transDepth := 0
		inCheck := b.IsAttacked(b.King[b.Stm], b.Stm.Opposite())
		bm := ml[0]
		bs := NoScore // bm keeps the best from prev iteration in case of immediate stop before first is done in this iteration
		for depth = 1; depth <= Limits.Depth && !Limits.Stop; depth++ {
			ml.Sort()
			bs = NoScore // bm keeps the best from prev iteration in case of immediate stop before first is done in this iterastion
			alpha, beta = MinEval, MaxEval
			for ix, mv := range ml { // root move loop
				childPV.Clear()

				b.Move(mv)
				// Tell(
				// 	fmt.Sprintf(
				// 		"info depth %v currmove %v currmovenumber %v",
				// 		depth,
				// 		mv.String(),
				// 		ix+1,
				// 	),
				// )
				lmrRed := 0
				ext := 0 // TODO: make extension function
				if ext == 0 {
					lmrRed = lmr(mv, inCheck, depth, ix+1, ix, b)
				}
				score := NoScore
				if ix == 0 {
					score = -Search(-beta, -alpha, depth-1+ext, 1, &childPV, b) // full search
				} else {
					score = -Search(-alpha-1, -alpha, depth-1+ext-lmrRed, 1, &childPV, b)
					if score > alpha && !Limits.Stop { // re-search due to PVS and/or lmr
						score = -Search(-beta, -alpha, depth-1+ext, 1, &childPV, b)
					}
				}

				b.Unmove(mv)

				if Limits.Stop {
					break
				}
				ml[ix].PackEval(score)
				if score > bs {
					bs = score
					pv.Catenate(mv, &childPV)

					bm = ml[ix]
					alpha = score
					transDepth = depth
					if depth >= 0 {
						position.Trans.Store(b.FullKey(), mv, transDepth, 0, score, ScoreTypeLower)
					}

					t1 := time.Since(Limits.StartTime)

					// Tell(
					fmt.Sprintf(
						"info score cp %v depth %v nodes %v time %v pv ",
						bm.Eval(),
						depth,
						CntNodes,
						int(t1.Seconds()*1000),
					)
					// 	pv.String(),
					// )
				}
			}
			if !Limits.Stop {
				ebfTab.Add(CntNodes)
			}

		} // end ID
		ml.Sort()

		position.Trans.Store(
			b.FullKey(),
			bm,
			transDepth,
			0,
			bs,
			position.ScoreType(bs, alpha, beta),
		)

		// time, nps, ebf
		t1 := time.Since(Limits.StartTime)
		Nps := float64(0)
		if t1.Seconds() != 0 {
			Nps = float64(CntNodes) / t1.Seconds()
		}
		ebfTab.Ebf(transDepth)
		// Tell(
		fmt.Sprintf(
			"info score cp %v depth %v nodes %v  time %v nps %v pv %v",
			bm.Eval(),
			transDepth,
			CntNodes,
			int(t1.Seconds()*1000),
			uint(Nps),
			pv.String(),
		)
		// )
		frEngine <- fmt.Sprintf("bestmove %v%v", Sq2Fen[bm.Fr()], Sq2Fen[bm.To()])
	}
}

// TODO search: Late Move Reduction

// TODO search: Internal Iterative Depening
// TODO search: Futility/Delta Pruning
// TODO search: more complicated time handling schemes
// TODO search: other reductions and extensions
func Search(alpha, beta, depth, ply int, pv *PvList, b *position.BoardStruct) int {
	CntNodes++
	if depth <= 0 {
		// return signEval(b.stm, evaluate(b))
		return Qs(beta, b)
	}

	// Are we in mate search?
	if mateSc := position.AddMatePly(MateEval-1, ply); mateSc < beta {
		beta = mateSc

		if mateSc <= alpha {
			return mateSc
		}
	}

	pv.Clear()

	pvNode := depth > 0 && beta != alpha+1

	transMove := moves.NoMove
	useTT := depth >= 0
	transDepth := depth
	inCheck := b.IsAttacked(b.King[b.Stm], b.Stm.Opposite())

	if depth < 0 && inCheck {
		useTT = true
		transDepth = 0
	}

	if useTT {
		var transSc, scType int
		ok := false

		if transMove, transSc, scType, ok = position.Trans.Retrieve(b.FullKey(), transDepth, ply); ok &&
			!pvNode {
			switch {
			case scType == ScoreTypeLower && transSc >= beta:
				position.Trans.CPrune++
				return transSc
			case scType == ScoreTypeUpper && transSc <= alpha:
				position.Trans.CPrune++
				return transSc
			case scType == ScoreTypeBetween:
				position.Trans.CPrune++
				return transSc
			}
		}
	}

	var childPV PvList
	childPV.New() // TODO? make it smaller for each depth maxDepth-ply
	/////////////////////////////////////// NULL MOVE /////////////////////////////////////////
	ev := SignEval(b.Stm, position.Evaluate(b))
	// null-move pruning
	if !pvNode && depth > 0 && !position.IsMateScore(beta) && !inCheck && !b.IsAntiNullMove() &&
		ev >= beta {
		nullMv := b.MoveNull()
		sc := MinEval
		if depth <= 3 { // static
			// if you don't beat me with 100 points,
			// then I think your position sucks
			sc = -Qs(-beta+1, b) // TODO: maybe 75-100 points bonus for opponent?
		} else { // dynamic
			sc = -Search(-beta, -beta+1, depth-3-1, ply, &childPV, b)
		}

		b.UndoNull(nullMv)

		if sc >= beta {
			if useTT {
				position.Trans.Store(b.FullKey(), moves.NoMove, transDepth, ply, sc, ScoreTypeLower)
			}
			return sc
		}
	}
	/////////////////////// NULL MOVE END //////////////////////

	bs, score := NoScore, NoScore
	bm := moves.NoMove

	genInfo := GenInfoStruct{Sv: 0, Ply: ply, TransMove: transMove}
	cntMoves := 0
	Next = NextNormal
	for mv, msg := Next(&genInfo, b); mv != moves.NoMove; mv, msg = Next(&genInfo, b) {
		_ = msg

		if !b.Move(mv) {
			continue
		}

		childPV.Clear()
		lmrRed := 0
		ext := 0 // TODO: make extension function
		if ext == 0 {
			lmrRed = lmr(mv, inCheck, depth, genInfo.Sv, cntMoves, b)
		}
		if pvNode && cntMoves == 0 {
			score = -Search(-beta, -alpha, depth-1+ext, ply+1, &childPV, b)
		} else {
			score = -Search(-alpha-1, -alpha, depth-1+ext-lmrRed, ply+1, &childPV, b)
			if score > alpha {
				score = -Search(-beta, -alpha, depth-1+ext, ply+1, &childPV, b)
			}
		}

		b.Unmove(mv)
		cntMoves++

		if score > bs {
			bs = score
			bm = mv
			pv.Catenate(mv, &childPV)
			if score > alpha {
				alpha = score
				if useTT {
					position.Trans.Store(
						b.FullKey(),
						mv,
						transDepth,
						ply,
						score,
						position.ScoreType(score, alpha, beta),
					)
				}
			}

			if score >= beta { // beta cutoff
				// add killer and update history
				if mv.Cp() == Empty && mv.Pr() == Empty {
					Killers.Add(mv, ply)
					History.Inc(mv.Fr(), mv.To(), b.Stm, depth)
				}
				if mv.Cmp(transMove) {
					position.Trans.CPrune++
				}
				return score
			}
		}

		tStep := time.Since(Limits.LastTime) - time.Duration(time.Millisecond*200)
		if tStep >= 0 {
			Limits.LastTime = time.Now().Add(-time.Duration(tStep))
			t1 := time.Since(Limits.StartTime)
			ms := uint64(t1.Nanoseconds() / 1000000)
			if t1.Seconds() > 1 {
				if ms%1000 <= 5 {
					// Tell(
					fmt.Sprintf(
						"info time %v nodes %v nps %v",
						ms,
						CntNodes,
						CntNodes/uint64(t1.Seconds()),
					)
					// )
				}
			}

			if ms >= uint64(Limits.MoveTime)-200 {
				//				fmt.Println("t1", uint64(t1.Nanoseconds()/1000000)-100, "limit", uint64(limits.moveTime))
				Limits.Stop = true
			}
		}

		if Limits.Stop {
			return alpha
		}
	}

	if cntMoves == 0 { // we didn't find any legal moves - either mate or stalemate
		sc := 0      // we could have a contempt value here instead
		if inCheck { // must be a mate
			sc = -MateEval + ply + 1
		}

		if useTT {
			position.Trans.Store(b.FullKey(), moves.NoMove, transDepth, ply, sc, ScoreTypeBetween)
		}
		return sc
	}

	if bm.Cmp(transMove) {
		position.Trans.CBest++
	}
	return bs
}

// compute late move reduction
func lmr(mv moves.Move, inCheck bool, depth, sv, CntMoves int, b *position.BoardStruct) int {
	interesting := inCheck || mv.Cp() != Empty || mv.Pr() != Empty ||
		b.IsAttacked(b.King[b.Stm], b.Stm.Opposite()) ||
		(b.Stm == WHITE && mv.Pc() == WP && mv.To() >= A6) ||
		(b.Stm == BLACK && mv.Pc() == BP && mv.To() <= H3) // even big threats? castling?
	red := 0
	if !interesting && depth >= 3 && sv >= NextFirstNonCp {
		red = 1
		if depth >= 5 && sv >= NextFirstNonCp {
			red = depth / 3
		}
	}
	return red
}

func initQS(ml *moves.MoveList, b *position.BoardStruct) {
	ml.Clear()
	b.GenAllCaptures(ml)
}

func Qs(beta int, b *position.BoardStruct) int {
	ev := SignEval(b.Stm, position.Evaluate(b))
	if ev >= beta {
		// we are good. No need to try captures
		return ev
	}
	bs := ev

	qsList := make(moves.MoveList, 0, 60)
	initQS(&qsList, b) // create attacks
	done := bitboard.BitBoard(0)

	// move loop
	for _, mv := range qsList {
		fr := mv.Fr()
		to := mv.To()

		// This works because we pick lower value pieces first
		if done.IsBitSet(to) { // Don't do the same to-sw again
			continue
		}
		done.SetBit(to)

		see := See(fr, to, b)

		if see == 0 && mv.Cp() == Empty {
			// must be a promotion that didn't captureand was not captured
			see = PieceVal[WQ] - PieceVal[WP]
		}

		if see <= 0 {
			continue // equal captures not interesting
		}

		sc := ev + see
		if sc > bs {
			bs = sc
			if sc >= beta {
				return sc
			}
		}
	}

	return bs
}

// see (Static Echange Evaluation)
// Start with the capture fr-to and find out all the other captures to to-sq
func See(fr, to int, b *position.BoardStruct) int {
	pVal := [16]int{
		100,
		-100,
		325,
		-325,
		325,
		-325,
		500,
		-500,
		950,
		-950,
		10000,
		-10000,
		0,
		0,
		0,
		0,
	}
	pc := b.Squares[fr]
	cp := b.Squares[to]
	cnt := 1
	us := position.PcColor(pc)
	them := us.Opposite()

	// All the attackers to the to-sq, but first remove the moving piece and use X-ray to the to-sq
	occ := b.AllBB()
	occ.Clear(fr)
	attackingBB := magic.MRookTab[to].Atks(occ)&(b.PieceBB[Rook]|b.PieceBB[Queen]) |
		magic.MBishopTab[to].Atks(occ)&(b.PieceBB[Bishop]|b.PieceBB[Queen]) |
		(position.AtksKnights[to] & b.PieceBB[Knight]) |
		(position.AtksKings[to] & b.PieceBB[King]) |
		(b.WPawnAtksFr(to) & b.PieceBB[Pawn] & b.WbBB[BLACK]) |
		(b.BPawnAtksFr(to) & b.PieceBB[Pawn] & b.WbBB[WHITE])
	attackingBB &= occ

	if (attackingBB & b.WbBB[them]) == 0 { // 'they' have no attackers - good bye
		return Abs(pVal[cp]) // always return score from 'our' point of view
	}

	// Now we continue to keep track of the material gain/loss for each capture
	// Always remove the last attacker and use x-ray to find possible new attackers

	lastAtkVal := Abs(pVal[pc]) // save attacker piece value for later use
	var captureList [32]int
	captureList[0] = Abs(pVal[cp])
	n := 1

	stm := them // change side to move

	for {
		cnt++

		var pt int
		switch { // select the least valuable attacker
		case (attackingBB & b.PieceBB[Pawn] & b.WbBB[stm]) != 0:
			pt = Pawn
		case (attackingBB & b.PieceBB[Knight] & b.WbBB[stm]) != 0:
			pt = Knight
		case (attackingBB & b.PieceBB[Bishop] & b.WbBB[stm]) != 0:
			pt = Bishop
		case (attackingBB & b.PieceBB[Rook] & b.WbBB[stm]) != 0:
			pt = Rook
		case (attackingBB & b.PieceBB[Queen] & b.WbBB[stm]) != 0:
			pt = Queen
		case (attackingBB & b.PieceBB[King] & b.WbBB[stm]) != 0:
			pt = King
		default:
			panic("Don't come here in see! ")
		}

		// now remove the pt above from the attackingBB and scan for new attackers by possible x-ray
		BB := attackingBB & (attackingBB & b.PieceBB[pt] & b.WbBB[stm])
		occ ^= (BB & -BB) // turn off the rightmost bit from BB in occ

		//  pick sliding attacks again (do it from to-sq)
		attackingBB |= magic.MRookTab[to].Atks(occ)&(b.PieceBB[Rook]|b.PieceBB[Queen]) |
			magic.MBishopTab[to].Atks(occ)&(b.PieceBB[Bishop]|b.PieceBB[Queen])
		attackingBB &= occ // but only attacking pieces

		captureList[n] = -captureList[n-1] + lastAtkVal
		n++

		// save the value of tha capturing piece to be used later
		lastAtkVal = pVal[position.Pt2pc(pt, WHITE)] // using WHITE always gives positive integer
		stm = stm.Opposite()                         // next side to move

		if pt == King && (attackingBB&b.WbBB[stm]) != 0 { // NOTE: just changed stm-color above
			// if king capture and 'they' are atting we have to stop
			captureList[n] = pVal[WK]
			n++
			break
		}

		if attackingBB&b.WbBB[stm] == 0 { // if no more attackers
			break
		}

	}

	// find the optimal capture sequence and 'our' material value will be on top
	for n--; n != 0; n-- {
		captureList[n-1] = min(-captureList[n], captureList[n-1])
	}

	return captureList[0]
}

/*
	 func genAndSort(b *boardStruct, ml *moveList) {
		ml.clear()
		b.genAllLegals(ml)
		for ix, mv := range *ml {
			b.move(mv)
			v := evaluate(b)
			b.unmove(mv)

			v = signEval(b.stm, v)
			(*ml)[ix].packEval(v)
		}

		ml.sort()
	}
*/
func genAndSort(ply int, b *position.BoardStruct, ml *moves.MoveList) {
	if ply > MaxPly {
		panic("wtf maxply")
	}

	ml.Clear()
	b.GenAllLegals(ml)

	for ix, mv := range *ml {
		b.Move(mv)
		v := position.Evaluate(b)
		b.Unmove(mv)
		if Killers[ply].K1.Cmp(mv) {
			v += 1000
		} else if Killers[ply].K2.Cmp(mv) {
			v += 900
		}

		v = SignEval(b.Stm, v)

		(*ml)[ix].PackEval(v)
	}

	ml.Sort()
}

// generate capture moves first, then killers, then non captures
func genInOrder(b *position.BoardStruct, ml *moves.MoveList, ply int, transMove moves.Move) {
	ml.Clear()
	b.GenAllCaptures(ml)
	noCaptIx := len(*ml)
	b.GenAllNonCaptures(ml)
	if transMove != moves.NoMove {
		for ix := 0; ix < len(*ml); ix++ {
			mv := (*ml)[ix]
			if transMove.Cmp(mv) {
				(*ml)[ix], (*ml)[0] = (*ml)[0], (*ml)[ix]
				break
			}
		}
	}
	pos1, pos2 := noCaptIx, noCaptIx+1
	if (*ml)[pos1].Cmp(transMove) {
		noCaptIx++
		pos1++
		pos2++
	}

	if len(*ml)-noCaptIx > 2 {
		// place killers first among non captures
		cnt := 0
		for ix := noCaptIx; ix < len(*ml); ix++ {
			mv := (*ml)[ix]
			switch {
			case Killers[ply].K1.CmpFrTo(mv) && !mv.CmpFrTo(transMove) && b.Squares[mv.To()] == Empty:
				mv.PackMove(
					mv.Fr(),
					mv.To(),
					b.Squares[mv.Fr()],
					b.Squares[mv.To()],
					mv.Pr(),
					b.Ep,
					b.Castlings,
				)
				(*ml)[ix] = mv
				(*ml)[ix], (*ml)[pos1] = (*ml)[pos1], (*ml)[ix]
				cnt++
			case Killers[ply].K2.CmpFrTo(mv) && !mv.CmpFrTo(transMove) && b.Squares[mv.To()] == Empty:
				mv.PackMove(
					mv.Fr(),
					mv.To(),
					b.Squares[mv.Fr()],
					b.Squares[mv.To()],
					mv.Pr(),
					b.Ep,
					b.Castlings,
				)
				(*ml)[ix] = mv
				(*ml)[ix], (*ml)[pos2] = (*ml)[pos2], (*ml)[ix]
				cnt++
			}
			if cnt >= 2 {
				break
			}
		}
	}
}

func SignEval(stm Color, ev int) int {
	if stm == BLACK {
		return -ev
	}
	return ev
}

// ///////////////////////// Next move /////////////////////////////////////
var Next func(*GenInfoStruct, *position.BoardStruct) (moves.Move, string) // or nextKEvasion or nextQS

const (
	InitNext = iota
	NextTr
	NextFirstGoodCp
	NextGoodCp
	NextK1
	NextK2
	NextCounterMv
	NextFirstNonCp
	NextNonCp
	NextBadCp
	NextEnd
)

type GenInfoStruct struct {
	// to be filled in, before first call to the next-function
	Sv, Ply   int
	TransMove moves.Move

	// handle by the next-function
	Captures, NonCapt moves.MoveList
	CounterMv         moves.Move
}

func NextNormal(genInfo *GenInfoStruct, b *position.BoardStruct) (moves.Move, string) {
	switch genInfo.Sv {
	case InitNext:
		genInfo.Sv = NextTr
		fallthrough
	case NextTr:
		genInfo.Sv = NextFirstGoodCp
		if genInfo.TransMove != moves.NoMove {
			if b.IsLegal(genInfo.TransMove) {
				return genInfo.TransMove, "transMove"
			}
			genInfo.TransMove = moves.NoMove
		}
		fallthrough
	case NextFirstGoodCp:
		genInfo.Captures.New(20)
		b.GenAllCaptures(&genInfo.Captures)
		// pick a good capt - use see - not transMove
		bs := -1
		bIx := 0
		ml := &genInfo.Captures
		for ix := 0; ix < len(*ml); ix++ {
			if (*ml)[ix].Cmp(genInfo.TransMove) {
				continue
			}
			sc := See((*ml)[ix].Fr(), (*ml)[ix].To(), b)
			(*ml)[ix].PackEval(sc)
			if sc > bs {
				bs = sc
				bIx = ix
			}
		}
		if bs >= 0 {
			mv := (*ml)[bIx]
			(*ml)[bIx], (*ml)[len(*ml)-1] = (*ml)[len(*ml)-1], (*ml)[bIx]
			*ml = (*ml)[:len(*ml)-1]
			genInfo.Sv = NextGoodCp
			return mv, "first good capt"
		}

		genInfo.Sv = NextK1
		fallthrough
	case NextGoodCp:
		// pick a good capt - use see - not transMove
		bs := -1
		bIx := 0
		ml := &genInfo.Captures
		for ix := 0; ix < len(*ml); ix++ {
			if (*ml)[ix].Cmp(genInfo.TransMove) {
				continue
			}
			sc := (*ml)[ix].Eval()
			if sc > bs {
				bs = sc
				bIx = ix
			}
		}
		if bs >= 0 {
			mv := (*ml)[bIx]
			(*ml)[bIx], (*ml)[len(*ml)-1] = (*ml)[len(*ml)-1], (*ml)[bIx]
			*ml = (*ml)[:len(*ml)-1]
			bs, bIx = MinEval, -1
			return mv, "good capt"
		}
		genInfo.Sv = NextK1
		fallthrough
	case NextK1: // not transMove
		genInfo.Sv = NextK2
		if Killers[genInfo.Ply].K1 != moves.NoMove &&
			!genInfo.TransMove.CmpFrToP(Killers[genInfo.Ply].K1) {
			if b.IsLegal(Killers[genInfo.Ply].K1) {
				var mv moves.Move
				mv.PackMove(
					Killers[genInfo.Ply].K1.Fr(),
					Killers[genInfo.Ply].K1.To(),
					b.Squares[Killers[genInfo.Ply].K1.Fr()],
					b.Squares[Killers[genInfo.Ply].K1.To()],
					Killers[genInfo.Ply].K1.Pr(),
					b.Ep,
					b.Castlings,
				)
				return mv, "K1"
			}
		}

		fallthrough
	case NextK2: // not transMove
		genInfo.Sv = NextCounterMv
		if Killers[genInfo.Ply].K2 != moves.NoMove &&
			!genInfo.TransMove.CmpFrToP(Killers[genInfo.Ply].K2) {
			if b.IsLegal(Killers[genInfo.Ply].K2) {
				var mv moves.Move
				mv.PackMove(
					Killers[genInfo.Ply].K2.Fr(),
					Killers[genInfo.Ply].K2.To(),
					b.Squares[Killers[genInfo.Ply].K2.Fr()],
					b.Squares[Killers[genInfo.Ply].K2.To()],
					Killers[genInfo.Ply].K2.Pr(),
					b.Ep,
					b.Castlings,
				)
				return mv, "K2"
			}
		}

		fallthrough
	case NextCounterMv: // not transMove, not killer1, not killer2 genInfo.CounterMv = moves.NoMove
		genInfo.Sv = NextFirstNonCp
		//	if genInfo.counterMv != noMove && !genInfo.counterMv.cmpFrTo(genInfo.transMove) &&
		//     genInfo.counterMv.cmpFrTo(killers[genInfo.ply].k1) &&  genInfo.counterMv.cmpFrTo(killers[genInfo.ply].k2) {
		//		var mv move
		//		mv.packMove(counterMv.fr(), counterMv.to(),b.sq[counterMv.fr()],b.sq[counterMv.to()],counterMv.pr(), b.ep,b.castlings)
		// check if it is a valid move
		//		return sv, counterMovex[genInfo.ply][mv.to()]
		//	}

		fallthrough
	case NextFirstNonCp: // not transMove, not counterMove, not killer1, not killer2
		genInfo.NonCapt.New(50)
		ml := &genInfo.NonCapt
		b.GenAllNonCaptures(ml)
		// pick by HistoryTab (see will probably not give anything) - I don't want to sort it. hist may change between moves
		bs := MinEval
		bIx := -1
		for ix := 0; ix < len(*ml); ix++ {
			if (*ml)[ix].CmpFrToP(genInfo.TransMove) || (*ml)[ix].CmpFrToP(genInfo.CounterMv) ||
				(*ml)[ix].CmpFrToP(
					Killers[genInfo.Ply].K1,
				) || (*ml)[ix].CmpFrToP(Killers[genInfo.Ply].K2) {
				continue
			}
			sc := int(History.Get((*ml)[ix].Fr(), (*ml)[ix].To(), b.Stm))
			if sc > bs {
				bs = sc
				bIx = ix
			}
		}
		if bIx >= 0 {
			mv := (*ml)[bIx]

			(*ml)[bIx], (*ml)[len(*ml)-1] = (*ml)[len(*ml)-1], (*ml)[bIx]
			*ml = (*ml)[:len(*ml)-1]

			genInfo.Sv = NextNonCp

			return mv, "first non capt"
		}

		genInfo.Sv = NextBadCp
		fallthrough
	case NextNonCp: // not transMove, not counterMove, not killer1, not killer2
		// pick by HistoryTab (see will probably not give anything)
		bs := MinEval
		bIx := -1
		ml := &genInfo.NonCapt
		for ix := 0; ix < len(*ml); ix++ {
			if (*ml)[ix].CmpFrToP(genInfo.TransMove) || (*ml)[ix].CmpFrToP(genInfo.CounterMv) ||
				(*ml)[ix].CmpFrToP(
					Killers[genInfo.Ply].K1,
				) || (*ml)[ix].CmpFrToP(Killers[genInfo.Ply].K2) {
				continue
			}
			sc := int(History.Get((*ml)[ix].Fr(), (*ml)[ix].To(), b.Stm))
			if sc > bs {
				bs = sc
				bIx = ix
			}
		}
		if bIx >= 0 {
			mv := (*ml)[bIx]
			(*ml)[bIx], (*ml)[len(*ml)-1] = (*ml)[len(*ml)-1], (*ml)[bIx]
			*ml = (*ml)[:len(*ml)-1]

			return mv, "non Capt"
		}

		genInfo.Sv = NextBadCp
		fallthrough
	case NextBadCp: // not transMove
		// pick a bad capt  - use see?
		mv := moves.NoMove
		ml := &genInfo.Captures
		for ix := len(*ml) - 1; ix >= 0; ix-- {
			if (*ml)[ix].Cmp(genInfo.TransMove) {
				*ml = (*ml)[:len(*ml)-1]
				continue
			}

			mv = (*ml)[ix]
			//		(*ml)[ix], (*ml)[len(*ml)-1] = (*ml)[len(*ml)-1], (*ml)[ix]
			*ml = (*ml)[:len(*ml)-1]
			break
		}

		return mv, "bad capt"
	default: // shouldn't happen
		panic("never come here! nextNormal sv=" + strconv.Itoa(genInfo.Sv))
	}
}

// StartPerft starts the Perft command that generates all moves until the given depth.
// It counts the leafs only taht is printed out for each possible move from current pos
func StartPerft(depth int, bd *position.BoardStruct) uint64 {
	if depth <= 0 {
		fmt.Printf("Total:\t%v\n", 1)
		return 0
	}

	transMove := moves.NoMove
	transMove, _, _, _ = position.Trans.Retrieve(bd.FullKey(), depth, 0)

	totCount := uint64(0)
	genInfo := GenInfoStruct{Sv: 0, Ply: 0, TransMove: transMove}
	Next = NextNormal
	ix := 0
	for mv, msg := Next(&genInfo, bd); mv != moves.NoMove; mv, msg = Next(&genInfo, bd) {
		if !bd.Move(mv) {
			continue
		}
		dbg := false
		/*
			/////////////////////////////////////////////////////////////
			if mv.fr() == D4 && mv.to() == F4 {
				dbg = true
			}
			/////////////////////////////////////////////////////////////
		*/
		count := perft(dbg, depth-1, 1, bd)
		totCount += count
		fmt.Printf("%2d: %v \t%v \t%v\n", ix+1, mv.String(), count, msg)

		bd.Unmove(mv)
		ix++
	}
	fmt.Println("------------------")
	fmt.Println()
	fmt.Printf("Total:\t%v\n", totCount)
	return totCount
}

func perft(dbg bool, depth, ply int, bd *position.BoardStruct) uint64 {
	if depth == 0 {
		return 1
	}

	transMove := moves.NoMove
	transMove, _, _, _ = position.Trans.Retrieve(bd.FullKey(), depth, ply)
	ix := 0
	count := uint64(0)
	genInfo := GenInfoStruct{Sv: 0, Ply: ply, TransMove: transMove}
	Next = NextNormal
	for mv, msg := Next(&genInfo, bd); mv != moves.NoMove; mv, msg = Next(&genInfo, bd) {
		if !bd.Move(mv) {
			continue
		}
		_ = msg
		deb := false
		/*
			////////////////////////////////////////////////////////////////
			if dbg && mv.fr() == F5 && mv.to() == F4 {
				deb = true
			}
			if dbg && mv.fr() == E2 && mv.to() == E4 {
				deb = true
			}
			////////////////////////////////////////////////////////////////
		*/
		cnt := perft(deb, depth-1, ply+1, bd)
		count += cnt
		/*
			/////////////////////////////////////////////
			if dbg && !deb  {
				fmt.Println(ix+1, ":(e4)     ", mv.String(), msg, "\t", cnt)
				if ix==1{
					fmt.Println("K1",killers[ply].k1.StringFull())
					fmt.Println("K2",killers[ply].k2.StringFull())
				}
			}
			////////////////////////////////////////////
		*/
		bd.Unmove(mv)
		ix++
	}

	return count
}
