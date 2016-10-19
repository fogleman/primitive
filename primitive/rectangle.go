package primitive

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
)

type Rectangle struct {
	Worker *Worker
	X1, Y1 int
	X2, Y2 int
}

// NewRandomRectangle uses the passed Worker's random seed to create a new Rectangle.
// Returns a reference to the new Rectangle.
func NewRandomRectangle(worker *Worker) *Rectangle {
	rnd := worker.Rnd
	x1 := rnd.Intn(worker.W)
	y1 := rnd.Intn(worker.H)
	x2 := clampInt(x1+rnd.Intn(32)+1, 0, worker.W-1)
	y2 := clampInt(y1+rnd.Intn(32)+1, 0, worker.H-1)
	return &Rectangle{worker, x1, y1, x2, y2}
}

// bounds regularizes the Rectangle shape by making the first x,y vertex
// the minimum (left) vertex.
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

// Draw adds/paints the Rectangle shape to the Go Graphics Context.
func (r *Rectangle) Draw(dc *gg.Context, scale float64) {
	x1, y1, x2, y2 := r.bounds()
	dc.DrawRectangle(float64(x1), float64(y1), float64(x2-x1+1), float64(y2-y1+1))
	dc.Fill()
}

// SVG returns the string markup that represents the Rectangle in SVG format.
func (r *Rectangle) SVG(attrs string) string {
	x1, y1, x2, y2 := r.bounds()
	w := x2 - x1 + 1
	h := y2 - y1 + 1
	return fmt.Sprintf(
		"<rect %s x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" />",
		attrs, x1, y1, w, h)
}

// Copy returns a reference to a copy of the rectangle.
func (r *Rectangle) Copy() Shape {
	a := *r
	return &a
}

// Mutate randomly changes the shape of the Rectangle randomly
// adjusting its X,Y vertex values.
func (r *Rectangle) Mutate() {
	w := r.Worker.W
	h := r.Worker.H
	rnd := r.Worker.Rnd
	switch rnd.Intn(2) {
	case 0:
		r.X1 = clampInt(r.X1+int(rnd.NormFloat64()*16), 0, w-1)
		r.Y1 = clampInt(r.Y1+int(rnd.NormFloat64()*16), 0, h-1)
	case 1:
		r.X2 = clampInt(r.X2+int(rnd.NormFloat64()*16), 0, w-1)
		r.Y2 = clampInt(r.Y2+int(rnd.NormFloat64()*16), 0, h-1)
	}
}

// Rasterize converts the Rectangle shape object into scanlines of pixels.
func (r *Rectangle) Rasterize() []Scanline {
	x1, y1, x2, y2 := r.bounds()
	lines := r.Worker.Lines[:0]
	for y := y1; y <= y2; y++ {
		lines = append(lines, Scanline{y, x1, x2, 0xffff})
	}
	return lines
}

type RotatedRectangle struct {
	Worker *Worker
	X, Y   int
	Sx, Sy int
	Angle  int
}

// NewRandomRotatedRectangle uses the passed Worker's random seed to create a new rotated Rectangle.
// Returns a reference to the new rotated Rectangle.
func NewRandomRotatedRectangle(worker *Worker) *RotatedRectangle {
	rnd := worker.Rnd
	x := rnd.Intn(worker.W)
	y := rnd.Intn(worker.H)
	sx := rnd.Intn(32) + 1
	sy := rnd.Intn(32) + 1
	a := rnd.Intn(360)
	r := &RotatedRectangle{worker, x, y, sx, sy, a}
	r.Mutate()
	return r
}

// Draw adds the Rotated Rectangle object into the Go Graphics Context.
func (r *RotatedRectangle) Draw(dc *gg.Context, scale float64) {
	sx, sy := float64(r.Sx), float64(r.Sy)
	dc.Push()
	dc.Translate(float64(r.X), float64(r.Y))
	dc.Rotate(radians(float64(r.Angle)))
	dc.DrawRectangle(-sx/2, -sy/2, sx, sy)
	dc.Pop()
	dc.Fill()
}

// SVG returns the string markup that represents the RotatedRectangle in SVG format.
func (r *RotatedRectangle) SVG(attrs string) string {
	return fmt.Sprintf(
		"<g transform=\"translate(%d %d) rotate(%d) scale(%d %d)\"><rect %s x=\"-0.5\" y=\"-0.5\" width=\"1\" height=\"1\" /></g>",
		r.X, r.Y, r.Angle, r.Sx, r.Sy, attrs)
}

// Copy returns a reference to a copy of the rotated rectangle.
func (r *RotatedRectangle) Copy() Shape {
	a := *r
	return &a
}

// Mutate randomly changes the shape of the RotatedRectangle randomly
// adjusting its X,Y vertex values and Angle.
func (r *RotatedRectangle) Mutate() {
	w := r.Worker.W
	h := r.Worker.H
	rnd := r.Worker.Rnd
	switch rnd.Intn(3) {
	case 0:
		r.X = clampInt(r.X+int(rnd.NormFloat64()*16), 0, w-1)
		r.Y = clampInt(r.Y+int(rnd.NormFloat64()*16), 0, h-1)
	case 1:
		r.Sx = clampInt(r.Sx+int(rnd.NormFloat64()*16), 1, w-1)
		r.Sy = clampInt(r.Sy+int(rnd.NormFloat64()*16), 1, h-1)
	case 2:
		r.Angle = r.Angle + int(rnd.NormFloat64()*32)
	}
	// for !r.Valid() {
	// 	r.Sx = clampInt(r.Sx+int(rnd.NormFloat64()*16), 0, w-1)
	// 	r.Sy = clampInt(r.Sy+int(rnd.NormFloat64()*16), 0, h-1)
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

// Rasterize converts the RotatedRectangle into Scanlines of pixels.
// Returns a slice of Scanlines which contain the pixel representation of the
// RotatedRectangle.
func (r *RotatedRectangle) Rasterize() []Scanline {
	w := r.Worker.W
	h := r.Worker.H
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
		min[i] = w
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
	lines := r.Worker.Lines[:0]
	for i := 0; i < n; i++ {
		y := miny + i
		if y < 0 || y >= h {
			continue
		}
		a := maxInt(min[i], 0)
		b := minInt(max[i], w-1)
		if b >= a {
			lines = append(lines, Scanline{y, a, b, 0xffff})
		}
	}
	return lines
}
