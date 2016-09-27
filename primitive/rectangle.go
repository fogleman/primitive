package primitive

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

type Rectangle struct {
	W, H   int
	X1, Y1 int
	X2, Y2 int
}

func NewRandomRectangle(w, h int, rnd *rand.Rand) *Rectangle {
	x1 := rnd.Intn(w)
	y1 := rnd.Intn(h)
	x2 := rnd.Intn(w)
	y2 := rnd.Intn(h)
	return &Rectangle{w, h, x1, y1, x2, y2}
}

func (r *Rectangle) bounds() (x1, y1, x2, y2 int) {
	x1, y1 = r.X1, r.Y1
	x2, y2 = r.X2, r.Y2
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	return
}

func (r *Rectangle) Draw(dc *gg.Context) {
	x1, y1, x2, y2 := r.bounds()
	dc.DrawRectangle(float64(x1), float64(y1), float64(x2-x1+1), float64(y2-y1+1))
}

func (r *Rectangle) SVG(attrs string) string {
	x1, y1, x2, y2 := r.bounds()
	w := x2 - x1 + 1
	h := y2 - y1 + 1
	return fmt.Sprintf(
		"<rect %s x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" />",
		attrs, x1, y1, w, h)
}

func (r *Rectangle) Copy() Shape {
	a := *r
	return &a
}

func (r *Rectangle) Mutate(rnd *rand.Rand) {
	switch rnd.Intn(2) {
	case 0:
		r.X1 = clampInt(r.X1+rnd.Intn(21)-10, 0, r.W-1)
		r.Y1 = clampInt(r.Y1+rnd.Intn(21)-10, 0, r.H-1)
	case 1:
		r.X2 = clampInt(r.X2+rnd.Intn(21)-10, 0, r.W-1)
		r.Y2 = clampInt(r.Y2+rnd.Intn(21)-10, 0, r.H-1)
	}
}

func (r *Rectangle) Rasterize() []Scanline {
	x1, y1, x2, y2 := r.bounds()
	lines := make([]Scanline, y2-y1+1)
	i := 0
	for y := y1; y <= y2; y++ {
		lines[i] = Scanline{y, x1, x2}
		i++
	}
	return lines
}

type RotatedRectangle struct {
	W, H   int
	X, Y   int
	Sx, Sy int
	Angle  int
}

func NewRandomRotatedRectangle(w, h int, rnd *rand.Rand) *RotatedRectangle {
	x := rnd.Intn(w)
	y := rnd.Intn(h)
	sx := rnd.Intn(32) + 1
	sy := rnd.Intn(32) + 1
	a := rnd.Intn(360)
	r := &RotatedRectangle{w, h, x, y, sx, sy, a}
	r.Mutate(rnd)
	return r
}

func (r *RotatedRectangle) Draw(dc *gg.Context) {
	sx, sy := float64(r.Sx), float64(r.Sy)
	dc.Push()
	dc.Translate(float64(r.X), float64(r.Y))
	dc.Rotate(radians(float64(r.Angle)))
	dc.DrawRectangle(-sx/2, -sy/2, sx, sy)
	dc.Pop()
}

func (r *RotatedRectangle) SVG(attrs string) string {
	return fmt.Sprintf(
		"<g transform=\"translate(%d %d) rotate(%d) scale(%d %d)\"><rect %s x=\"-0.5\" y=\"-0.5\" width=\"1\" height=\"1\" /></g>",
		r.X, r.Y, r.Angle, r.Sx, r.Sy, attrs)
}

func (r *RotatedRectangle) Copy() Shape {
	a := *r
	return &a
}

func (r *RotatedRectangle) Mutate(rnd *rand.Rand) {
	switch rnd.Intn(3) {
	case 0:
		r.X = clampInt(r.X+rnd.Intn(21)-10, 0, r.W-1)
		r.Y = clampInt(r.Y+rnd.Intn(21)-10, 0, r.H-1)
	case 1:
		r.Sx = clampInt(r.Sx+rnd.Intn(21)-10, 1, r.W-1)
		r.Sy = clampInt(r.Sy+rnd.Intn(21)-10, 1, r.H-1)
	case 2:
		r.Angle = r.Angle + rnd.Intn(41) - 20
	}
	// for !r.Valid() {
	// 	r.Sx = clampInt(r.Sx+rnd.Intn(21)-10, 0, r.W-1)
	// 	r.Sy = clampInt(r.Sy+rnd.Intn(21)-10, 0, r.H-1)
	// }
}

func (r *RotatedRectangle) Valid() bool {
	a, b := r.Sx, r.Sy
	if a < b {
		a, b = b, a
	}
	aspect := float64(a) / float64(b)
	return aspect <= 5
}

func (r *RotatedRectangle) Rasterize() []Scanline {
	sx, sy := float64(r.Sx), float64(r.Sy)
	angle := radians(float64(r.Angle))
	rx1, ry1 := rotate(-sx/2, -sy/2, angle)
	rx2, ry2 := rotate(sx/2, -sy/2, angle)
	rx3, ry3 := rotate(sx/2, sy/2, angle)
	rx4, ry4 := rotate(-sx/2, sy/2, angle)
	x1, y1 := int(rx1)+r.X, int(ry1)+r.Y
	x2, y2 := int(rx2)+r.X, int(ry2)+r.Y
	x3, y3 := int(rx3)+r.X, int(ry3)+r.Y
	x4, y4 := int(rx4)+r.X, int(ry4)+r.Y
	miny := minInt(y1, minInt(y2, minInt(y3, y4)))
	maxy := maxInt(y1, maxInt(y2, maxInt(y3, y4)))
	n := maxy - miny + 1
	min := make([]int, n)
	max := make([]int, n)
	for i := range min {
		min[i] = r.W
	}
	xs := []int{x1, x2, x3, x4, x1}
	ys := []int{y1, y2, y3, y4, y1}
	// TODO: this could be better probably
	for i := 0; i < 4; i++ {
		x, y := float64(xs[i]), float64(ys[i])
		dx, dy := float64(xs[i+1]-xs[i]), float64(ys[i+1]-ys[i])
		count := int(math.Sqrt(dx*dx+dy*dy)) * 2
		for j := 0; j < count; j++ {
			t := float64(j) / float64(count-1)
			xi := int(x + dx*t)
			yi := int(y+dy*t) - miny
			min[yi] = minInt(min[yi], xi)
			max[yi] = maxInt(max[yi], xi)
		}
	}
	lines := make([]Scanline, 0, n)
	for i := 0; i < n; i++ {
		y := miny + i
		if y < 0 || y >= r.H {
			continue
		}
		a := maxInt(min[i], 0)
		b := minInt(max[i], r.W-1)
		if b >= a {
			lines = append(lines, Scanline{y, a, b})
		}
	}
	return lines
}
