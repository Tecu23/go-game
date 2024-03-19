// Package engine contains the main logic behind the chess engine
package engine

import (
	"math"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	maxDepth = 100
	maxPly   = 100
)

var countNodes uint64

// TODO search limits: count nodes and test for limit.nodes
// TODO search limits: limit.depth

// TODO search limits: time per game w/wo increments
// TODO search limits: time per x moves and after x moves w/wo increments
type searchLimits struct {
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
var Limits searchLimits

func (s *searchLimits) Init() {
	s.Depth = 9999
	s.Nodes = math.MaxUint64
	s.MoveTime = 99999999999
	s.Infinite = false
	s.Stop = false
}

func (s *searchLimits) SetStop(st bool) {
	s.Stop = st
}

func (s *searchLimits) SetDepth(d int) {
	s.Depth = d
}

func (s *searchLimits) SetMoveTime(m int) {
	s.MoveTime = m
}

func (s *searchLimits) SetInfinite(b bool) {
	s.Infinite = b
}

// Engine should create the 2 channels necessary to communicate to the websocket
func Engine() (chan string, chan string) {
	frEngine := make(chan string)
	toEngine := make(chan string)

	go func() {
		for cmd := range toEngine {
			log.Infof("engine got %s", cmd)
			switch cmd {
			case "stop":
			case "quit":
			}
		}
	}()

	return toEngine, frEngine
}
