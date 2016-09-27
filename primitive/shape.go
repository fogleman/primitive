package primitive

import (
	"math/rand"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Shape interface {
	Rasterize(r *raster.Rasterizer) []Scanline
	Copy() Shape
	Mutate(rnd *rand.Rand)
	Draw(dc *gg.Context, scale float64)
	SVG(attrs string) string
}

type ShapeType int

const (
	ShapeTypeAny ShapeType = iota
	ShapeTypeTriangle
	ShapeTypeRectangle
	ShapeTypeEllipse
	ShapeTypeCircle
	ShapeTypeRotatedRectangle
	ShapeTypePath
)
