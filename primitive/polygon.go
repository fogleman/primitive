package primitive

import (
	"fmt"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Polygon struct {
	Worker *Worker
	Order  int
	Convex bool
	X, Y   []float64
}

//NewRandomPolygon creates random polygon
func NewRandomPolygon(worker *Worker, order int, convex bool) *Polygon {
	rnd := worker.Rnd
	x := make([]float64, order)
	y := make([]float64, order)
	x[0] = rnd.Float64() * float64(worker.W)
	y[0] = rnd.Float64() * float64(worker.H)
	for i := 1; i < order; i++ {
		x[i] = x[0] + rnd.Float64()*40 - 20
		y[i] = y[0] + rnd.Float64()*40 - 20
	}
	p := &Polygon{worker, order, convex, x, y}
	p.Mutate()
	return p
}

func (p *Polygon) Draw(dc *gg.Context, scale float64) {
	dc.NewSubPath()
	for i := 0; i < p.Order; i++ {
		dc.LineTo(p.X[i], p.Y[i])
	}
	dc.ClosePath()
	dc.Fill()
}

func (p *Polygon) SVG(attrs string) string {
	ret := fmt.Sprintf(
		"<polygon %s points=\"",
		attrs)
	points := make([]string, len(p.X))
	for i := 0; i < len(p.X); i++ {
		points[i] = fmt.Sprintf("%f,%f", p.X[i], p.Y[i])
	}

	return ret + strings.Join(points, ",") + "\" />"
}

func (p *Polygon) Copy() Shape {
	a := *p
	a.X = make([]float64, p.Order)
	a.Y = make([]float64, p.Order)
	copy(a.X, p.X)
	copy(a.Y, p.Y)
	return &a
}

func (p *Polygon) Mutate() {
	const m = 16
	w := p.Worker.W
	h := p.Worker.H
	rnd := p.Worker.Rnd
	for {
		if rnd.Float64() < 0.25 {
			i := rnd.Intn(p.Order)
			j := rnd.Intn(p.Order)
			p.X[i], p.Y[i], p.X[j], p.Y[j] = p.X[j], p.Y[j], p.X[i], p.Y[i]
		} else {
			i := rnd.Intn(p.Order)
			p.X[i] = clamp(p.X[i]+rnd.NormFloat64()*16, -m, float64(w-1+m))
			p.Y[i] = clamp(p.Y[i]+rnd.NormFloat64()*16, -m, float64(h-1+m))
		}
		if p.Valid() {
			break
		}
	}
}

func (p *Polygon) Valid() bool {
	if !p.Convex {
		return true
	}
	var sign bool
	for a := 0; a < p.Order; a++ {
		i := (a + 0) % p.Order
		j := (a + 1) % p.Order
		k := (a + 2) % p.Order
		c := cross3(p.X[i], p.Y[i], p.X[j], p.Y[j], p.X[k], p.Y[k])
		if a == 0 {
			sign = c > 0
		} else if c > 0 != sign {
			return false
		}
	}
	return true
}

func cross3(x1, y1, x2, y2, x3, y3 float64) float64 {
	dx1 := x2 - x1
	dy1 := y2 - y1
	dx2 := x3 - x2
	dy2 := y3 - y2
	return dx1*dy2 - dy1*dx2
}

func (p *Polygon) Rasterize() []Scanline {
	var path raster.Path
	for i := 0; i <= p.Order; i++ {
		f := fixp(p.X[i%p.Order], p.Y[i%p.Order])
		if i == 0 {
			path.Start(f)
		} else {
			path.Add1(f)
		}
	}
	return fillPath(p.Worker, path)
}
