package bitboard

import (
	"testing"

	assert "github.com/Tecu23/go-engine/internal/testing"
)

func TestBitBoardCount(t *testing.T) {
	tests := []struct {
		id   string
		b    bitBoard
		want int
	}{
		{
			id:   "1",
			b:    0x0,
			want: 0,
		}, {
			id:   "2",
			b:    0x10,
			want: 1,
		}, {
			id:   "3",
			b:    0x2cd0ab4173295da4,
			want: 29,
		}, {
			id:   "4",
			b:    0xe128a79595994b75,
			want: 32,
		}, {
			id:   "5",
			b:    0xffffffffffffffff,
			want: 64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			res := tt.b.count()

			assert.Equal(t, res, tt.want)
		})
	}
}

func Test_bitBoard_some(t *testing.T) {
	tests := []struct {
		name string
		b    bitBoard
		pos  int
	}{
		{"", 0xF, 63},
		{"", 0x0, 0},
		{"", 0x1, 0},
		{"", 0xFFFF, 63},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.set(tt.pos)
			if !tt.b.test(tt.pos) {
				t.Fatalf("set(%v) gives %v in %b. Want %v", tt.pos, false, tt.b, true)
			}

			tt.b.clr(tt.pos)
			if tt.b.test(tt.pos) {
				t.Errorf("clr(%v) gives %v in %b. Want %v", tt.pos, true, tt.b, false)
			}
		})
	}
}

func Test_bitBoard_firstOne(t *testing.T) {
	tests := []struct {
		name string
		b    bitBoard
		want int
	}{
		{"", 0x0, 64},
		{"", 0x1, 0},
		{"", 0xFFFFFFFFFFFFFFFF, 0},
		{"", 0xFFFFFFFFFFFFFF00, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := tt.b
			if got := tt.b.firstOne(); got != tt.want {
				t.Errorf("bitBoard.firstOne(%x) = %v, want %v", x, got, tt.want)
			}
		})
	}
}

func Test_bitBoard_lasttOne(t *testing.T) {
	tests := []struct {
		name string
		b    bitBoard
		want int
	}{
		{"", 0x0, 64},
		{"", 0x1, 0},
		{"", 0xFFFFFFFFFFFFFFFF, 63},
		{"", 0x00FFFFFFFFFFFFFF, 55},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := tt.b
			if got := tt.b.lastOne(); got != tt.want {
				t.Errorf("bitBoard.lastOne(%x) = %v, want %v", x, got, tt.want)
			}
		})
	}
}
