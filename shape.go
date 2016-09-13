package tri

type Shape interface {
	Rasterize() []Scanline
	Copy() Shape
	Mutate()
}
