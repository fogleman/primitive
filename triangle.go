package tri

import (
	"image"
	"image/color"
	"math"
	"math/rand"
)

type Triangle struct {
	X1, Y1 int
	X2, Y2 int
	X3, Y3 int
}

func RandomTriangle(w, h int) Triangle {
	x1 := rand.Intn(w)
	y1 := rand.Intn(h)
	x2 := rand.Intn(w)
	y2 := rand.Intn(h)
	x3 := rand.Intn(w)
	y3 := rand.Intn(h)
	return Triangle{x1, y1, x2, y2, x3, y3}
}

func (t *Triangle) BoundingBox() (x1, y1, x2, y2 int) {
	x1 = int(math.Min(float64(t.X1), math.Min(float64(t.X2), float64(t.X3))))
	x2 = int(math.Max(float64(t.X1), math.Max(float64(t.X2), float64(t.X3))))
	y1 = int(math.Min(float64(t.Y1), math.Min(float64(t.Y2), float64(t.Y3))))
	y2 = int(math.Max(float64(t.Y1), math.Max(float64(t.Y2), float64(t.Y3))))
	return
}

func (t *Triangle) Draw(im *image.RGBA, c Color) {
	fillTriangle(im, c.Color(), t.X1, t.Y1, t.X2, t.Y2, t.X3, t.Y3)
}

func fillTriangle(im *image.RGBA, c color.Color, x1, y1, x2, y2, x3, y3 int) {
	if y1 > y3 {
		x1, x3 = x3, x1
		y1, y3 = y3, y1
	}
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	if y2 > y3 {
		x2, x3 = x3, x2
		y2, y3 = y3, y2
	}
	if y2 == y3 {
		flatBottom(im, c, x1, y1, x2, y2, x3, y3)
	} else if y1 == y2 {
		flatTop(im, c, x1, y1, x2, y2, x3, y3)
	} else {
		x4 := x1 + int((float64(y2-y1)/float64(y3-y1))*float64(x3-x1))
		y4 := y2
		flatBottom(im, c, x1, y1, x2, y2, x4, y4)
		flatTop(im, c, x2, y2, x4, y4, x3, y3)
	}
}

func flatBottom(im *image.RGBA, c color.Color, x1, y1, x2, y2, x3, y3 int) {
	const m = 1<<16 - 1
	sr, sg, sb, sa := c.RGBA()
	aa := (m - sa) * 0x101
	s1 := float64(x2-x1) / float64(y2-y1)
	s2 := float64(x3-x1) / float64(y3-y1)
	ax := float64(x1)
	bx := float64(x1)
	for y := y1; y <= y2; y++ {
		a := int(ax)
		b := int(bx)
		if a > b {
			a, b = b, a
		}
		i := im.PixOffset(a, y)
		for x := a; x <= b; x++ {
			dr := &im.Pix[i+0]
			dg := &im.Pix[i+1]
			db := &im.Pix[i+2]
			da := &im.Pix[i+3]
			i += 4
			*dr = uint8((uint32(*dr)*aa/m + sr) >> 8)
			*dg = uint8((uint32(*dg)*aa/m + sg) >> 8)
			*db = uint8((uint32(*db)*aa/m + sb) >> 8)
			*da = uint8((uint32(*da)*aa/m + sa) >> 8)
		}
		ax += s1
		bx += s2
	}
}

func flatTop(im *image.RGBA, c color.Color, x1, y1, x2, y2, x3, y3 int) {
	const m = 1<<16 - 1
	sr, sg, sb, sa := c.RGBA()
	aa := (m - sa) * 0x101
	s1 := float64(x3-x1) / float64(y3-y1)
	s2 := float64(x3-x2) / float64(y3-y2)
	ax := float64(x3)
	bx := float64(x3)
	for y := y3; y > y1; y-- {
		ax -= s1
		bx -= s2
		a := int(ax)
		b := int(bx)
		if a > b {
			a, b = b, a
		}
		i := im.PixOffset(a, y)
		for x := a; x <= b; x++ {
			dr := &im.Pix[i+0]
			dg := &im.Pix[i+1]
			db := &im.Pix[i+2]
			da := &im.Pix[i+3]
			i += 4
			*dr = uint8((uint32(*dr)*aa/m + sr) >> 8)
			*dg = uint8((uint32(*dg)*aa/m + sg) >> 8)
			*db = uint8((uint32(*db)*aa/m + sb) >> 8)
			*da = uint8((uint32(*da)*aa/m + sa) >> 8)
		}
	}
}
