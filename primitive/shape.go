package primitive

import "github.com/fogleman/gg"

type Shape interface {
	Rasterize(worker *Worker) []Scanline
	Copy() Shape
	Scale(s float64) Shape
	Mutate(worker *Worker)
	Draw(dc *gg.Context, scale float64)
	SVG(attrs string) string
	Command() string
}

type ShapeType int

const (
	ShapeTypeAny ShapeType = iota
	ShapeTypeTriangle
	ShapeTypeRectangle
	ShapeTypeEllipse
	ShapeTypeCircle
	ShapeTypeRotatedRectangle
	ShapeTypeQuadratic
	ShapeTypeRotatedEllipse
	ShapeTypePolygon
)
