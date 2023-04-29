package primitive

import "github.com/fogleman/gg"

// Shape is an interface that models a shope configuration such as rectangle
// or triangle
type Shape interface {
	Rasterize() []Scanline
	Copy() Shape
	Mutate()
	Draw(dc *gg.Context, scale float64, notify Notifier)
	SVG(attrs string) string
}

// ShapeType specifies the type of shape that is modeled
type ShapeType int

const (
	shapeTypeAny ShapeType = iota
	shapeTypeTriangle
	shapeTypeRectangle
	shapeTypeEllipse
	shapeTypeCircle
	shapeTypeRotatedRectangle
	shapeTypeQuadratic
	shapeTypeRotatedEllipse
	shapeTypePoligon
)
