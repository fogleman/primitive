package primitive

import (
	"math/rand"

	"github.com/fogleman/gg"
)

type Rectangle struct {
	W, H   int
	X1, Y1 int
	X2, Y2 int
}

func NewRandomRectangle(w, h int) *Rectangle {
	x1 := rand.Intn(w)
	y1 := rand.Intn(h)
	x2 := rand.Intn(w)
	y2 := rand.Intn(h)
	return &Rectangle{w, h, x1, y1, x2, y2}
}

func (r *Rectangle) Draw(dc *gg.Context) {
	x1, y1 := r.X1, r.Y1
	x2, y2 := r.X2, r.Y2
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	dc.DrawRectangle(pt(x1), pt(y1), pt(x2-x1+1), pt(y2-y1+1))
	dc.Fill()
}

func (r *Rectangle) Copy() Shape {
	a := *r
	return &a
}

func (r *Rectangle) Mutate() {
	switch rand.Intn(2) {
	case 0:
		r.X1 = clampInt(r.X1+rand.Intn(21)-10, 0, r.W-1)
		r.Y1 = clampInt(r.Y1+rand.Intn(21)-10, 0, r.H-1)
	case 1:
		r.X2 = clampInt(r.X2+rand.Intn(21)-10, 0, r.W-1)
		r.Y2 = clampInt(r.Y2+rand.Intn(21)-10, 0, r.H-1)
	}
}

func (r *Rectangle) Rasterize() []Scanline {
	x1, y1 := r.X1, r.Y1
	x2, y2 := r.X2, r.Y2
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	lines := make([]Scanline, y2-y1+1)
	i := 0
	for y := y1; y <= y2; y++ {
		lines[i] = Scanline{y, x1, x2}
		i++
	}
	return lines
}
