package entity

type Player struct {
	name         string
	playingPiece PlayingPiece
	id           int
}

func NewPlayer(name string, piece PlayingPiece, id int) Player {
	return Player{
		name:         name,
		playingPiece: piece,
		id:           id,
	}
}
