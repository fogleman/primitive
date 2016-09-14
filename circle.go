package primitive

import (
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

type Circle struct {
	W, H    int
	X, Y, R int
}

func NewRandomCircle(w, h int) *Circle {
	x := rand.Intn(w)
	y := rand.Intn(h)
	r := rand.Intn(w)
	return &Circle{w, h, x, y, r}
}

func (c *Circle) Draw(dc *gg.Context) {
	dc.DrawCircle(pt(c.X), pt(c.Y), pt(c.R))
}

func (c *Circle) Copy() Shape {
	a := *c
	return &a
}

func (c *Circle) Mutate() {
	switch rand.Intn(2) {
	case 0:
		c.X = clampInt(c.X+rand.Intn(21)-10, 0, c.W-1)
		c.Y = clampInt(c.Y+rand.Intn(21)-10, 0, c.H-1)
	case 1:
		c.R = clampInt(c.R+rand.Intn(21)-10, 0, c.W-1)
	}
}

func (c *Circle) Rasterize() []Scanline {
	lines := make([]Scanline, 0, c.R*2)
	for dy := 0; dy < c.R; dy++ {
		y1 := c.Y - dy
		y2 := c.Y + dy
		if (y1 < 0 || y1 >= c.H) && (y2 < 0 || y2 >= c.H) {
			continue
		}
		w := int(math.Sqrt(float64(c.R*c.R - dy*dy)))
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
