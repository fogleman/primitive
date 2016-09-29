package primitive

import (
	"fmt"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Quadratic struct {
	Worker *Worker
	X1, Y1 float64
	X2, Y2 float64
	X3, Y3 float64
	Width  float64
}

func NewRandomQuadratic(worker *Worker) *Quadratic {
	rnd := worker.Rnd
	x1 := rnd.Float64() * float64(worker.W)
	y1 := rnd.Float64() * float64(worker.H)
	x2 := x1 + rnd.Float64()*40 - 20
	y2 := y1 + rnd.Float64()*40 - 20
	x3 := x2 + rnd.Float64()*40 - 20
	y3 := y2 + rnd.Float64()*40 - 20
	width := 1.0
	q := &Quadratic{worker, x1, y1, x2, y2, x3, y3, width}
	q.Mutate()
	return q
}

func (q *Quadratic) Draw(dc *gg.Context, scale float64) {
	dc.MoveTo(float64(q.X1), float64(q.Y1))
	dc.QuadraticTo(float64(q.X2), float64(q.Y2), float64(q.X3), float64(q.Y3))
	dc.SetLineWidth(float64(q.Width) * scale)
	dc.Stroke()
}

func (q *Quadratic) SVG(attrs string) string {
	// TODO: this is a little silly
	attrs = strings.Replace(attrs, "fill", "stroke", -1)
	return fmt.Sprintf(
		"<path %s fill=\"none\" d=\"M %f %f Q %f %f, %f %f\" stroke-width=\"%f\" />",
		attrs, q.X1, q.Y1, q.X2, q.Y2, q.X3, q.Y3, q.Width)
}

func (q *Quadratic) Copy() Shape {
	a := *q
	return &a
}

func (q *Quadratic) Mutate() {
	const m = 16
	w := q.Worker.W
	h := q.Worker.H
	rnd := q.Worker.Rnd
	for {
		switch rnd.Intn(3) {
		case 0:
			q.X1 = clamp(q.X1+rnd.NormFloat64()*16, -m, float64(w-1+m))
			q.Y1 = clamp(q.Y1+rnd.NormFloat64()*16, -m, float64(h-1+m))
		case 1:
			q.X2 = clamp(q.X2+rnd.NormFloat64()*16, -m, float64(w-1+m))
			q.Y2 = clamp(q.Y2+rnd.NormFloat64()*16, -m, float64(h-1+m))
		case 2:
			q.X3 = clamp(q.X3+rnd.NormFloat64()*16, -m, float64(w-1+m))
			q.Y3 = clamp(q.Y3+rnd.NormFloat64()*16, -m, float64(h-1+m))
		case 3:
			q.Width = clamp(q.Width+rnd.NormFloat64(), 1, 16)
		}
		if q.Valid() {
			break
		}
	}
}

func (q *Quadratic) Valid() bool {
	dx12 := int(q.X1 - q.X2)
	dy12 := int(q.Y1 - q.Y2)
	dx23 := int(q.X2 - q.X3)
	dy23 := int(q.Y2 - q.Y3)
	dx13 := int(q.X1 - q.X3)
	dy13 := int(q.Y1 - q.Y3)
	d12 := dx12*dx12 + dy12*dy12
	d23 := dx23*dx23 + dy23*dy23
	d13 := dx13*dx13 + dy13*dy13
	return d13 > d12 && d13 > d23
}

func (q *Quadratic) Rasterize() []Scanline {
	var path raster.Path
	p1 := fixp(q.X1, q.Y1)
	p2 := fixp(q.X2, q.Y2)
	p3 := fixp(q.X3, q.Y3)
	path.Start(p1)
	path.Add2(p2, p3)
	width := fix(q.Width)
	return strokePath(q.Worker, path, width, raster.RoundCapper, raster.RoundJoiner)
}
