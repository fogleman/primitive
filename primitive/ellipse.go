package primitive

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Ellipse struct {
	Worker *Worker
	X, Y   int
	Rx, Ry int
	Circle bool
}

func NewRandomEllipse(worker *Worker) *Ellipse {
	rnd := worker.Rnd
	x := rnd.Intn(worker.W)
	y := rnd.Intn(worker.H)
	rx := rnd.Intn(32) + 1
	ry := rnd.Intn(32) + 1
	return &Ellipse{worker, x, y, rx, ry, false}
}

func NewRandomCircle(worker *Worker) *Ellipse {
	rnd := worker.Rnd
	x := rnd.Intn(worker.W)
	y := rnd.Intn(worker.H)
	r := rnd.Intn(32) + 1
	return &Ellipse{worker, x, y, r, r, true}
}

func (c *Ellipse) Draw(dc *gg.Context, scale float64) {
	dc.DrawEllipse(float64(c.X), float64(c.Y), float64(c.Rx), float64(c.Ry))
	dc.Fill()
}

func (c *Ellipse) SVG(attrs string) string {
	return fmt.Sprintf(
		"<ellipse %s cx=\"%d\" cy=\"%d\" rx=\"%d\" ry=\"%d\" />",
		attrs, c.X, c.Y, c.Rx, c.Ry)
}

func (c *Ellipse) Copy() Shape {
	a := *c
	return &a
}

func (c *Ellipse) Mutate() {
	w := c.Worker.W
	h := c.Worker.H
	rnd := c.Worker.Rnd
	switch rnd.Intn(3) {
	case 0:
		c.X = clampInt(c.X+int(rnd.NormFloat64()*16), 0, w-1)
		c.Y = clampInt(c.Y+int(rnd.NormFloat64()*16), 0, h-1)
	case 1:
		c.Rx = clampInt(c.Rx+int(rnd.NormFloat64()*16), 1, w-1)
		if c.Circle {
			c.Ry = c.Rx
		}
	case 2:
		c.Ry = clampInt(c.Ry+int(rnd.NormFloat64()*16), 1, w-1)
		if c.Circle {
			c.Rx = c.Ry
		}
	}
}

func (c *Ellipse) Rasterize() []Scanline {
	w := c.Worker.W
	h := c.Worker.H
	lines := c.Worker.Lines[:0]
	aspect := float64(c.Rx) / float64(c.Ry)
	for dy := 0; dy < c.Ry; dy++ {
		y1 := c.Y - dy
		y2 := c.Y + dy
		if (y1 < 0 || y1 >= h) && (y2 < 0 || y2 >= h) {
			continue
		}
		s := int(math.Sqrt(float64(c.Ry*c.Ry-dy*dy)) * aspect)
		x1 := c.X - s
		x2 := c.X + s
		if x1 < 0 {
			x1 = 0
		}
		if x2 >= w {
			x2 = w - 1
		}
		if y1 >= 0 && y1 < h {
			lines = append(lines, Scanline{y1, x1, x2, 0xffff})
		}
		if y2 >= 0 && y2 < h && dy > 0 {
			lines = append(lines, Scanline{y2, x1, x2, 0xffff})
		}
	}
	return lines
}

type RotatedEllipse struct {
	Worker *Worker
	X, Y   float64
	Rx, Ry float64
	Angle  float64
}

func NewRandomRotatedEllipse(worker *Worker) *RotatedEllipse {
	rnd := worker.Rnd
	x := rnd.Float64() * float64(worker.W)
	y := rnd.Float64() * float64(worker.H)
	rx := rnd.Float64()*32 + 1
	ry := rnd.Float64()*32 + 1
	a := rnd.Float64() * 360
	return &RotatedEllipse{worker, x, y, rx, ry, a}
}

func (c *RotatedEllipse) Draw(dc *gg.Context, scale float64) {
	dc.Push()
	dc.RotateAbout(radians(c.Angle), c.X, c.Y)
	dc.DrawEllipse(c.X, c.Y, c.Rx, c.Ry)
	dc.Fill()
	dc.Pop()
}

func (c *RotatedEllipse) SVG(attrs string) string {
	return fmt.Sprintf(
		"<g transform=\"translate(%f %f) rotate(%f) scale(%f %f)\"><ellipse %s cx=\"0\" cy=\"0\" rx=\"1\" ry=\"1\" /></g>",
		c.X, c.Y, c.Angle, c.Rx, c.Ry, attrs)
}

func (c *RotatedEllipse) Copy() Shape {
	a := *c
	return &a
}

func (c *RotatedEllipse) Mutate() {
	w := c.Worker.W
	h := c.Worker.H
	rnd := c.Worker.Rnd
	switch rnd.Intn(3) {
	case 0:
		c.X = clamp(c.X+rnd.NormFloat64()*16, 0, float64(w-1))
		c.Y = clamp(c.Y+rnd.NormFloat64()*16, 0, float64(h-1))
	case 1:
		c.Rx = clamp(c.Rx+rnd.NormFloat64()*16, 1, float64(w-1))
		c.Ry = clamp(c.Ry+rnd.NormFloat64()*16, 1, float64(w-1))
	case 2:
		c.Angle = c.Angle + rnd.NormFloat64()*32
	}
}

func (c *RotatedEllipse) Rasterize() []Scanline {
	var path raster.Path
	const n = 16
	for i := 0; i < n; i++ {
		p1 := float64(i+0) / n
		p2 := float64(i+1) / n
		a1 := p1 * 2 * math.Pi
		a2 := p2 * 2 * math.Pi
		x0 := c.Rx * math.Cos(a1)
		y0 := c.Ry * math.Sin(a1)
		x1 := c.Rx * math.Cos(a1+(a2-a1)/2)
		y1 := c.Ry * math.Sin(a1+(a2-a1)/2)
		x2 := c.Rx * math.Cos(a2)
		y2 := c.Ry * math.Sin(a2)
		cx := 2*x1 - x0/2 - x2/2
		cy := 2*y1 - y0/2 - y2/2
		x0, y0 = rotate(x0, y0, radians(c.Angle))
		cx, cy = rotate(cx, cy, radians(c.Angle))
		x2, y2 = rotate(x2, y2, radians(c.Angle))
		if i == 0 {
			path.Start(fixp(x0+c.X, y0+c.Y))
		}
		path.Add2(fixp(cx+c.X, cy+c.Y), fixp(x2+c.X, y2+c.Y))
	}
	return fillPath(c.Worker, path)
}
