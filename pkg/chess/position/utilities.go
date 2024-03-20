package position

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Tecu23/go-game/pkg/chess/castlings"
	. "github.com/Tecu23/go-game/pkg/chess/constants"
	"github.com/Tecu23/go-game/pkg/chess/moves"
)

// ParseFEN should parse a FEN string and retrieve the board
// rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq -
func ParseFEN(FEN string) {
	Board.Clear()

	fenIdx := 0
	sq := 0

	// parsing the FEN from the start and setting the from top to bottom
	for row := 7; row >= 0; row-- {
		for sq = row * 8; sq < row*8+8; {

			char := string(FEN[fenIdx])
			fenIdx++

			if char == "/" {
				continue
			}

			// if we find a number we should skip that many squares from our current board
			if i, err := strconv.Atoi(char); err == nil {
				for j := 0; j < i; j++ {
					Board.SetSq(Empty, sq)
					sq++
				}
				continue
			}

			// if we find an invalid piece we skip
			if strings.IndexAny(PcFen, char) == -1 {
				log.Errorf("error string invalid piece %s try next one", char)
				continue
			}

			Board.SetSq(Fen2pc(char), sq)

			sq++
		}
	}

	remaining := strings.Split(strings.TrimSpace(FEN[fenIdx:]), " ")

	// Setting the Side to Move
	if len(remaining) > 0 {
		if remaining[0] == "w" {
			Board.Stm = WHITE
		} else if remaining[0] == "b" {
			Board.Stm = BLACK
		} else {
			log.Errorf("info string remaining=%v; sq=%v;  fenIx=%v;", strings.Join(remaining, " "), sq, fenIdx)
			log.Errorf("info string %s invalid stm color", remaining[0])
			Board.Stm = WHITE
		}
	}

	if Board.Stm == BLACK {
		Board.Key = ^Board.Key
	}

	// Checking for castling
	Board.Castlings = 0
	if len(remaining) > 1 {
		Board.Castlings = castlings.ParseCastlings(remaining[1])
	}

	// En Passant
	Board.Ep = 0
	if len(remaining) > 2 {
		if remaining[2] != "-" {
			Board.Ep = Fen2Sq[remaining[2]]
		}
	}

	// Cheking for 50 move rule
	Board.Rule50 = 0
	if len(remaining) > 3 {
		Board.Rule50 = parse50(remaining[3])
	}
}

// parseMvs should parse and make the moves retrieved from the position command
func ParseMvs(mvstr string) error {
	mvs := strings.Fields(strings.ToLower(mvstr))

	for _, mv := range mvs {
		mv = strings.TrimSpace(mv)

		if len(mv) < 4 || len(mv) > 5 {
			e := fmt.Sprintf("error string %s in the position command is not a correct move", mv)
			log.Error(e)
			return errors.New(e)

		}

		// does the from square exists
		fr, ok := Fen2Sq[mv[:2]]
		if !ok {
			e := fmt.Sprintf(
				"error string %s in the position command is not a correct from square",
				mv,
			)
			log.Error(e)
			return errors.New(e)
		}

		pc := Board.Squares[fr]
		if pc == Empty {
			e := fmt.Sprintf(
				"error string %s in the position command.The from square is an empty square",
				mv,
			)
			log.Error(e)
			return errors.New(e)
		}

		pcCol := PcColor(pc)
		if pcCol != Board.Stm {
			e := fmt.Sprintf(
				"error string %s in the position command.The from piece has the wrong color",
				mv,
			)
			log.Error(e)
			return errors.New(e)

		}

		// does the to square exists
		to, ok := Fen2Sq[mv[2:4]]
		if !ok {
			e := fmt.Sprintf(
				"error string %s in the position command has an incorrect from square",
				mv,
			)
			log.Error(e)
			return errors.New(e)
		}

		// does the promotion piece exists?
		pr := Empty
		if len(mv) == 5 {
			if !strings.ContainsAny(mv[4:5], "QRNBqrnb") {
				e := fmt.Sprintf(
					"error string promotion piece in %s in the position command is not correct",
					mv,
				)
				log.Error(e)
				return errors.New(e)
			}

			pr = Fen2pc(mv[4:5])
			pt := Pc2pt(pr)
			pr = Pt2pc(pt, Board.Stm)
		}

		cp := Board.Squares[to]

		var intMv moves.Move // external move format
		intMv.PackMove(fr, to, pc, cp, pr, Board.Ep, Board.Castlings)

		if !Board.Move(intMv) {
			e := fmt.Sprintf(
				"error string %v-%v is an illegal move",
				Sq2Fen[fr], Sq2Fen[to],
			)
			log.Error(e)
			return errors.New(e)
		}
	}
	return nil
}

// parse the 50 move rule in remaining portion of the fenstring
func parse50(fen50 string) int {
	r50, err := strconv.Atoi(fen50)
	if err != nil || r50 < 0 {
		log.Errorf("error string 50 move rule in fenstring %s is not valid", fen50)
		return 0
	}

	return r50
}

// Fen2pc convert pieceString to pc int
func Fen2pc(c string) int {
	for p, x := range PcFen {
		if string(x) == c {
			return p
		}
	}
	return Empty
}

// Pc2Fen convert pc to fenString
func Pc2Fen(pc int) string {
	if pc == Empty {
		return " "
	}
	return string(PcFen[pc])
}

// Pc2pt returns the pt from pc
func Pc2pt(pc int) int {
	return pc >> 1
}

// PcColor returns the color of a pc form
func PcColor(pc int) Color {
	return Color(pc & 0x1)
}

// Pt2pc returns pc from pt and sd
func Pt2pc(pt int, sd Color) int {
	return (pt << 1) | int(sd)
}
