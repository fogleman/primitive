package tri

import (
	"image"
	"math"
)

func differenceFull(a, b *image.RGBA) float64 {
	size := a.Bounds().Size()
	w, h := size.X, size.Y
	var total uint64
	for y := 0; y < h; y++ {
		i := a.PixOffset(0, y)
		for x := 0; x < w; x++ {
			ar := int(a.Pix[i])
			ag := int(a.Pix[i+1])
			ab := int(a.Pix[i+2])
			br := int(b.Pix[i])
			bg := int(b.Pix[i+1])
			bb := int(b.Pix[i+2])
			i += 4
			dr := ar - br
			dg := ag - bg
			db := ab - bb
			total += uint64(dr*dr + dg*dg + db*db)
		}
	}
	return math.Sqrt(float64(total)/float64(w*h*3)) / 255
}

func differencePartial(target, before, after *image.RGBA, score float64, lines []Scanline) float64 {
	size := target.Bounds().Size()
	w, h := size.X, size.Y
	total := uint64(math.Pow(score*255, 2) * float64(w*h*3))
	for _, line := range lines {
		i := target.PixOffset(line.X1, line.Y)
		for x := line.X1; x <= line.X2; x++ {
			tr := int(target.Pix[i])
			tg := int(target.Pix[i+1])
			tb := int(target.Pix[i+2])
			br := int(before.Pix[i])
			bg := int(before.Pix[i+1])
			bb := int(before.Pix[i+2])
			ar := int(after.Pix[i])
			ag := int(after.Pix[i+1])
			ab := int(after.Pix[i+2])
			i += 4
			dr1 := tr - br
			dg1 := tg - bg
			db1 := tb - bb
			dr2 := tr - ar
			dg2 := tg - ag
			db2 := tb - ab
			total -= uint64(dr1*dr1 + dg1*dg1 + db1*db1)
			total += uint64(dr2*dr2 + dg2*dg2 + db2*db2)
		}
	}
	return math.Sqrt(float64(total)/float64(w*h*3)) / 255
}
