package tri

import (
	"image"
	"image/color"
)

type Color struct {
	R, G, B, A int
}

func (c *Color) Color() color.Color {
	return color.NRGBA{uint8(c.R), uint8(c.G), uint8(c.B), uint8(c.A)}
}

func computeColor(target, current *image.RGBA, t Triangle, alpha int) Color {
	a := float64(alpha) / 255
	mask := image.NewRGBA(target.Bounds())
	t.Draw(mask, Color{255, 255, 255, 255})
	x1, y1, x2, y2 := t.BoundingBox()
	var count int
	var rsum, gsum, bsum float64
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			if mask.Pix[mask.PixOffset(x, y)] == 0 {
				continue
			}
			count++
			t := target.RGBAAt(x, y)
			c := current.RGBAAt(x, y)
			rsum += (a*float64(c.R) - float64(c.R) + float64(t.R)) / a
			gsum += (a*float64(c.G) - float64(c.G) + float64(t.G)) / a
			bsum += (a*float64(c.B) - float64(c.B) + float64(t.B)) / a
		}
	}
	if count == 0 {
		return Color{}
	}
	r := ClampInt(int(rsum/float64(count)), 0, 255)
	g := ClampInt(int(gsum/float64(count)), 0, 255)
	b := ClampInt(int(bsum/float64(count)), 0, 255)
	return Color{r, g, b, alpha}
}
