package constants

type Color int

func (c Color) Opposite() Color {
	return c ^ 0x1
}

func (c Color) ToString() string {
	if c == WHITE {
		return "W"
	}
	return "B"
}
