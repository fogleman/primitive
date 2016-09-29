package primitive

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
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
