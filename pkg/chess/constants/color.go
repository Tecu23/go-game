package constants

type Color int

func (c Color) Opposite() Color {
	return c ^ 0x1
}

func (c Color) String() string {
	if c == WHITE {
		return "W"
	}
	return "B"
}
