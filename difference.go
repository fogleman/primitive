package tri

import (
	"image"
	"math"
)

func DifferenceScore(a, b *image.RGBA) float64 {
	size := a.Bounds().Size()
	w, h := size.X, size.Y
	var total uint64
	for y := 0; y < h; y++ {
		i := a.PixOffset(0, y)
		for x := 0; x < w; x++ {
			ar := int(a.Pix[i+0])
			ag := int(a.Pix[i+1])
			ab := int(a.Pix[i+2])
			br := int(b.Pix[i+0])
			bg := int(b.Pix[i+1])
			bb := int(b.Pix[i+2])
			i += 4
			dr := ar - br
			dg := ag - bg
			db := ab - bb
			total += uint64(dr*dr + dg*dg + db*db)
		}
	}
	rmse := math.Sqrt(float64(total) / float64(w*h*3))
	return rmse / 255
}
