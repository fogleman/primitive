package primitive

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

type Ellipse struct {
	W, H   int
	X, Y   int
	Rx, Ry int
	Circle bool
}

func NewRandomEllipse(w, h int, rnd *rand.Rand) *Ellipse {
	x := rnd.Intn(w)
	y := rnd.Intn(h)
	rx := rnd.Intn(w / 2)
	ry := rnd.Intn(h / 2)
	return &Ellipse{w, h, x, y, rx, ry, false}
}

func NewRandomCircle(w, h int, rnd *rand.Rand) *Ellipse {
	x := rnd.Intn(w)
	y := rnd.Intn(h)
	r := rnd.Intn(w / 4)
	return &Ellipse{w, h, x, y, r, r, true}
}

func (c *Ellipse) Draw(dc *gg.Context) {
	dc.DrawEllipse(float64(c.X), float64(c.Y), float64(c.Rx), float64(c.Ry))
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

func (c *Ellipse) Mutate(rnd *rand.Rand) {
	switch rnd.Intn(3) {
	case 0:
		c.X = clampInt(c.X+rnd.Intn(21)-10, 0, c.W-1)
		c.Y = clampInt(c.Y+rnd.Intn(21)-10, 0, c.H-1)
	case 1:
		c.Rx = clampInt(c.Rx+rnd.Intn(21)-10, 0, c.W-1)
		if c.Circle {
			c.Ry = c.Rx
		}
	case 2:
		c.Ry = clampInt(c.Ry+rnd.Intn(21)-10, 0, c.W-1)
		if c.Circle {
			c.Rx = c.Ry
		}
	}
}

func (c *Ellipse) Rasterize() []Scanline {
	lines := make([]Scanline, 0, c.Ry*2)
	aspect := float64(c.Rx) / float64(c.Ry)
	for dy := 0; dy < c.Ry; dy++ {
		y1 := c.Y - dy
		y2 := c.Y + dy
		if (y1 < 0 || y1 >= c.H) && (y2 < 0 || y2 >= c.H) {
			continue
		}
		w := int(math.Sqrt(float64(c.Ry*c.Ry-dy*dy)) * aspect)
		x1 := c.X - w
		x2 := c.X + w
		if x1 < 0 {
			x1 = 0
		}
		if x2 >= c.W {
			x2 = c.W - 1
		}
		if y1 >= 0 && y1 < c.H {
			lines = append(lines, Scanline{y1, x1, x2})
		}
		if y2 >= 0 && y2 < c.H && dy > 0 {
			lines = append(lines, Scanline{y2, x1, x2})
		}
	}
	return lines
}
