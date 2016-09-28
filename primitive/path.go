package primitive

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Path struct {
	W, H   int
	X1, Y1 float64
	X2, Y2 float64
	X3, Y3 float64
	Width  float64
}

func NewRandomPath(w, h int, rnd *rand.Rand) *Path {
	x1 := rnd.Float64() * float64(w)
	y1 := rnd.Float64() * float64(h)
	x2 := x1 + rnd.Float64()*40 - 20
	y2 := y1 + rnd.Float64()*40 - 20
	x3 := x2 + rnd.Float64()*40 - 20
	y3 := y2 + rnd.Float64()*40 - 20
	width := 1.0
	p := &Path{w, h, x1, y1, x2, y2, x3, y3, width}
	p.Mutate(rnd)
	return p
}

func (p *Path) Draw(dc *gg.Context, scale float64) {
	dc.MoveTo(float64(p.X1), float64(p.Y1))
	dc.QuadraticTo(float64(p.X2), float64(p.Y2), float64(p.X3), float64(p.Y3))
	dc.SetLineWidth(float64(p.Width) * scale)
	dc.Stroke()
}

func (p *Path) SVG(attrs string) string {
	// TODO: this is a little silly
	attrs = strings.Replace(attrs, "fill", "stroke", -1)
	return fmt.Sprintf(
		"<path %s fill=\"none\" d=\"M %f %f Q %f %f, %f %f\" stroke-width=\"%f\" />",
		attrs, p.X1, p.Y1, p.X2, p.Y2, p.X3, p.Y3, p.Width)
}

func (p *Path) Copy() Shape {
	a := *p
	return &a
}

func (p *Path) Mutate(rnd *rand.Rand) {
	const m = 16
	for {
		switch rnd.Intn(3) {
		case 0:
			p.X1 = clamp(p.X1+rnd.Float64()*21-10, -m, float64(p.W-1+m))
			p.Y1 = clamp(p.Y1+rnd.Float64()*21-10, -m, float64(p.H-1+m))
		case 1:
			p.X2 = clamp(p.X2+rnd.Float64()*21-10, -m, float64(p.W-1+m))
			p.Y2 = clamp(p.Y2+rnd.Float64()*21-10, -m, float64(p.H-1+m))
		case 2:
			p.X3 = clamp(p.X3+rnd.Float64()*21-10, -m, float64(p.W-1+m))
			p.Y3 = clamp(p.Y3+rnd.Float64()*21-10, -m, float64(p.H-1+m))
		case 3:
			p.Width = clamp(p.Width+rnd.Float64()*2-1, 1, 16)
		}
		if p.Valid() {
			break
		}
	}
}

func (p *Path) Valid() bool {
	dx12 := int(p.X1 - p.X2)
	dy12 := int(p.Y1 - p.Y2)
	dx23 := int(p.X2 - p.X3)
	dy23 := int(p.Y2 - p.Y3)
	dx13 := int(p.X1 - p.X3)
	dy13 := int(p.Y1 - p.Y3)
	d12 := dx12*dx12 + dy12*dy12
	d23 := dx23*dx23 + dy23*dy23
	d13 := dx13*dx13 + dy13*dy13
	return d13 > d12 && d13 > d23
}

func (p *Path) Rasterize(buf []Scanline) []Scanline {
	var path raster.Path
	p1 := fixp(p.X1, p.Y1)
	p2 := fixp(p.X2, p.Y2)
	p3 := fixp(p.X3, p.Y3)
	path.Start(p1)
	path.Add2(p2, p3)
	width := fix(p.Width)
	return strokePath(p.W, p.H, path, width, raster.RoundCapper, raster.RoundJoiner, buf)
}
