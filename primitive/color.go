package primitive

import (
	"fmt"
	"image/color"
	"strings"
)

type Color struct {
	R, G, B, A int
}

// MakeColor takes a color.Color type from the Golang standard library
// and returns an equivelant Color type from this package.
func MakeColor(c color.Color) Color {
	r, g, b, a := c.RGBA()
	return Color{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

// MakeHexColor takes a hex color format string and returns a
// Color type struct of equivelant RGBA value
func MakeHexColor(x string) Color {
	x = strings.Trim(x, "#")
	var r, g, b, a int
	a = 255
	switch len(x) {
	case 3:
		fmt.Sscanf(x, "%1x%1x%1x", &r, &g, &b)
		r = (r << 4) | r
		g = (g << 4) | g
		b = (b << 4) | b
	case 4:
		fmt.Sscanf(x, "%1x%1x%1x%1x", &r, &g, &b, &a)
		r = (r << 4) | r
		g = (g << 4) | g
		b = (b << 4) | b
		a = (a << 4) | a
	case 6:
		fmt.Sscanf(x, "%02x%02x%02x", &r, &g, &b)
	case 8:
		fmt.Sscanf(x, "%02x%02x%02x%02x", &r, &g, &b, &a)
	}
	return Color{r, g, b, a}
}

// NRGBA, which is implemented by the Color type, converts the
// Color's RBGA value to an equivelant non-alpha pre-multiplied color.
func (c *Color) NRGBA() color.NRGBA {
	return color.NRGBA{uint8(c.R), uint8(c.G), uint8(c.B), uint8(c.A)}
}
