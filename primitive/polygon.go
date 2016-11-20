package primitive

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Polygon struct {
	Order  int
	Convex bool
	X, Y   []float64
}

func NewRandomPolygon(worker *Worker, order int, convex bool) *Polygon {
	rnd := worker.Rnd
	m := worker.Mutation
	x := make([]float64, order)
	y := make([]float64, order)
	x[0] = rnd.Float64() * float64(worker.W)
	y[0] = rnd.Float64() * float64(worker.H)
	if worker.X != 0 {
		x[0] = float64(worker.X)
	}
	if worker.Y != 0 {
		y[0] = float64(worker.Y)
	}
	for i := 1; i < order; i++ {
		x[i] = x[0] + rnd.Float64()*m*2 - m
		y[i] = y[0] + rnd.Float64()*m*2 - m
	}
	p := &Polygon{order, convex, x, y}
	p.Mutate(worker)
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
	return ""
}

func (p *Polygon) Command() string {
	s := "polygon"
	for i := 0; i < p.Order; i++ {
		s += fmt.Sprintf(" %f %f", p.X[i], p.Y[i])
	}
	return s
}

func (p *Polygon) Copy() Shape {
	a := *p
	a.X = make([]float64, p.Order)
	a.Y = make([]float64, p.Order)
	copy(a.X, p.X)
	copy(a.Y, p.Y)
	return &a
}

func (p *Polygon) Scale(s float64) Shape {
	a := p.Copy().(*Polygon)
	for i := 0; i < a.Order; i++ {
		a.X[i] *= s
		a.Y[i] *= s
	}
	return a
}

func (p *Polygon) Mutate(worker *Worker) {
	const q = 16
	w := worker.W
	h := worker.H
	m := worker.Mutation
	rnd := worker.Rnd
	for {
		if rnd.Float64() < 0.25 {
			i := rnd.Intn(p.Order)
			j := rnd.Intn(p.Order)
			p.X[i], p.Y[i], p.X[j], p.Y[j] = p.X[j], p.Y[j], p.X[i], p.Y[i]
		} else {
			i := rnd.Intn(p.Order)
			p.X[i] = clamp(p.X[i]+rnd.NormFloat64()*m, -q, float64(w-1+q))
			p.Y[i] = clamp(p.Y[i]+rnd.NormFloat64()*m, -q, float64(h-1+q))
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

func (p *Polygon) Rasterize(worker *Worker) []Scanline {
	var path raster.Path
	for i := 0; i <= p.Order; i++ {
		f := fixp(p.X[i%p.Order], p.Y[i%p.Order])
		if i == 0 {
			path.Start(f)
		} else {
			path.Add1(f)
		}
	}
	return fillPath(worker, path)
}
