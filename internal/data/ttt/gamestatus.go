package entity

type GameStatus string

const (
	InProgress GameStatus = "In-Progress"
	Win        GameStatus = "Win"
	Draw       GameStatus = "Draw"
	Lose       GameStatus = "Lose"
)
