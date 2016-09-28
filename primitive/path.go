package primitive

import (
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Path struct {
	W, H   int
	X1, Y1 int
	X2, Y2 int
	X3, Y3 int
	Width  int
}

func NewRandomPath(w, h int, rnd *rand.Rand) *Path {
	x1 := rnd.Intn(w)
	y1 := rnd.Intn(h)
	x2 := x1 + rnd.Intn(41) - 20
	y2 := y1 + rnd.Intn(41) - 20
	x3 := x2 + rnd.Intn(41) - 20
	y3 := y2 + rnd.Intn(41) - 20
	width := 1
	return &Path{w, h, x1, y1, x2, y2, x3, y3, width}
}

func (p *Path) Draw(dc *gg.Context, scale float64) {
	dc.MoveTo(float64(p.X1), float64(p.Y1))
	dc.QuadraticTo(float64(p.X2), float64(p.Y2), float64(p.X3), float64(p.Y3))
	dc.SetLineWidth(float64(p.Width) * scale)
	dc.Stroke()
}

func (p *Path) SVG(attrs string) string {
	return ""
}

func (p *Path) Copy() Shape {
	a := *p
	return &a
}

func (p *Path) Mutate(rnd *rand.Rand) {
	const m = 16
	switch rnd.Intn(3) {
	case 0:
		p.X1 = clampInt(p.X1+rnd.Intn(21)-10, -m, p.W-1+m)
		p.Y1 = clampInt(p.Y1+rnd.Intn(21)-10, -m, p.H-1+m)
	case 1:
		p.X2 = clampInt(p.X2+rnd.Intn(21)-10, -m, p.W-1+m)
		p.Y2 = clampInt(p.Y2+rnd.Intn(21)-10, -m, p.H-1+m)
	case 2:
		p.X3 = clampInt(p.X3+rnd.Intn(21)-10, -m, p.W-1+m)
		p.Y3 = clampInt(p.Y3+rnd.Intn(21)-10, -m, p.H-1+m)
	case 3:
		p.Width = clampInt(p.Width+rnd.Intn(3)-1, 1, 16)
	}
}

func (p *Path) Rasterize() []Scanline {
	var path raster.Path
	p1 := fixp(float64(p.X1), float64(p.Y1))
	p2 := fixp(float64(p.X2), float64(p.Y2))
	p3 := fixp(float64(p.X3), float64(p.Y3))
	path.Start(p1)
	path.Add2(p2, p3)
	width := fix(float64(p.Width))
	return strokePath(p.W, p.H, path, width, raster.RoundCapper, raster.RoundJoiner)
}
