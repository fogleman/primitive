package tri

import "image/color"

type Color struct {
	R, G, B, A int
}

func (c *Color) NRGBA() color.NRGBA {
	return color.NRGBA{uint8(c.R), uint8(c.G), uint8(c.B), uint8(c.A)}
}
