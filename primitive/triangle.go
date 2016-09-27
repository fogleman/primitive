package primitive

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Triangle struct {
	W, H   int
	X1, Y1 int
	X2, Y2 int
	X3, Y3 int
}

func NewRandomTriangle(w, h int, rnd *rand.Rand) *Triangle {
	x1 := rnd.Intn(w)
	y1 := rnd.Intn(h)
	x2 := x1 + rnd.Intn(31) - 15
	y2 := y1 + rnd.Intn(31) - 15
	x3 := x1 + rnd.Intn(31) - 15
	y3 := y1 + rnd.Intn(31) - 15
	t := &Triangle{w, h, x1, y1, x2, y2, x3, y3}
	t.Mutate(rnd)
	return t
}

func (t *Triangle) Draw(dc *gg.Context, scale float64) {
	dc.LineTo(float64(t.X1), float64(t.Y1))
	dc.LineTo(float64(t.X2), float64(t.Y2))
	dc.LineTo(float64(t.X3), float64(t.Y3))
	dc.ClosePath()
	dc.Fill()
}

func (t *Triangle) SVG(attrs string) string {
	return fmt.Sprintf(
		"<polygon %s points=\"%d,%d %d,%d %d,%d\" />",
		attrs, t.X1, t.Y1, t.X2, t.Y2, t.X3, t.Y3)
}

func (t *Triangle) Copy() Shape {
	a := *t
	return &a
}

func (t *Triangle) Mutate(rnd *rand.Rand) {
	const m = 16
	for {
		switch rnd.Intn(3) {
		case 0:
			t.X1 = clampInt(t.X1+rnd.Intn(21)-10, -m, t.W-1+m)
			t.Y1 = clampInt(t.Y1+rnd.Intn(21)-10, -m, t.H-1+m)
		case 1:
			t.X2 = clampInt(t.X2+rnd.Intn(21)-10, -m, t.W-1+m)
			t.Y2 = clampInt(t.Y2+rnd.Intn(21)-10, -m, t.H-1+m)
		case 2:
			t.X3 = clampInt(t.X3+rnd.Intn(21)-10, -m, t.W-1+m)
			t.Y3 = clampInt(t.Y3+rnd.Intn(21)-10, -m, t.H-1+m)
		}
		if t.Valid() {
			break
		}
	}
}

func (t *Triangle) Valid() bool {
	const minDegrees = 15
	var a1, a2, a3 float64
	{
		x1 := float64(t.X2 - t.X1)
		y1 := float64(t.Y2 - t.Y1)
		x2 := float64(t.X3 - t.X1)
		y2 := float64(t.Y3 - t.Y1)
		d1 := math.Sqrt(x1*x1 + y1*y1)
		d2 := math.Sqrt(x2*x2 + y2*y2)
		x1 /= d1
		y1 /= d1
		x2 /= d2
		y2 /= d2
		a1 = degrees(math.Acos(x1*x2 + y1*y2))
	}
	{
		x1 := float64(t.X1 - t.X2)
		y1 := float64(t.Y1 - t.Y2)
		x2 := float64(t.X3 - t.X2)
		y2 := float64(t.Y3 - t.Y2)
		d1 := math.Sqrt(x1*x1 + y1*y1)
		d2 := math.Sqrt(x2*x2 + y2*y2)
		x1 /= d1
		y1 /= d1
		x2 /= d2
		y2 /= d2
		a2 = degrees(math.Acos(x1*x2 + y1*y2))
	}
	a3 = 180 - a1 - a2
	return a1 > minDegrees && a2 > minDegrees && a3 > minDegrees
}

func (t *Triangle) Rasterize(r *raster.Rasterizer) []Scanline {
	var path raster.Path
	p1 := fixp(float64(t.X1), float64(t.Y1))
	p2 := fixp(float64(t.X2), float64(t.Y2))
	p3 := fixp(float64(t.X3), float64(t.Y3))
	path.Start(p1)
	path.Add1(p2)
	path.Add1(p3)
	path.Add1(p1)
	lines := fillPath(r, path)
	return lines
}
