package primitive

import (
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
	"golang.org/x/image/math/fixed"
)

type painter struct {
	Lines []Scanline
}

func (p *painter) Paint(spans []raster.Span, done bool) {
	for _, span := range spans {
		line := Scanline{span.Y, span.X0, span.X1 - 1}
		p.Lines = append(p.Lines, line)
	}
}

func fillPath(w, h int, path raster.Path) []Scanline {
	r := raster.NewRasterizer(w, h)
	r.UseNonZeroWinding = true
	r.AddPath(path)
	var p painter
	r.Rasterize(&p)
	return p.Lines
}

func strokePath(w, h int, path raster.Path, width fixed.Int26_6, cr raster.Capper, jr raster.Joiner) []Scanline {
	r := raster.NewRasterizer(w, h)
	r.UseNonZeroWinding = true
	r.AddStroke(path, width, cr, jr)
	var p painter
	r.Rasterize(&p)
	return p.Lines
}

type Path struct {
	W, H   int
	X1, Y1 int
	X2, Y2 int
	X3, Y3 int
	Width  int
	rnd    *rand.Rand
}

func NewRandomPath(w, h int, rnd *rand.Rand) *Path {
	x1 := rnd.Intn(w)
	y1 := rnd.Intn(h)
	x2 := x1 + rnd.Intn(21) - 10
	y2 := y1 + rnd.Intn(21) - 10
	x3 := x2 + rnd.Intn(21) - 10
	y3 := y2 + rnd.Intn(21) - 10
	width := 4 //rnd.Intn(16) + 1
	return &Path{w, h, x1, y1, x2, y2, x3, y3, width, rnd}
}

func (p *Path) Draw(dc *gg.Context) {
	dc.MoveTo(float64(p.X1), float64(p.Y1))
	dc.QuadraticTo(float64(p.X2), float64(p.Y2), float64(p.X3), float64(p.Y3))
	dc.SetLineWidth(float64(p.Width * 6))
	dc.Stroke()
}

func (p *Path) SVG(attrs string) string {
	return ""
}

func (p *Path) Copy() Shape {
	a := *p
	return &a
}

func (p *Path) Mutate() {
	const m = 16
	rnd := p.rnd
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
		p.Width = clampInt(p.Width+rnd.Intn(5)-2, 1, 64)
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
	lines := strokePath(p.W, p.H, path, width, raster.RoundCapper, raster.RoundJoiner)
	return cropScanlines(lines, p.W, p.H)
}
