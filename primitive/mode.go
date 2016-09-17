package primitive

type Mode int

const (
	ModeAny Mode = iota
	ModeTriangle
	ModeRectangle
	ModeEllipse
	ModeCircle
	ModeRotatedRectangle
)
