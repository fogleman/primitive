package tri

import (
	"fmt"
	"image"
	"image/color"
)

func score(target, current *image.RGBA, t Triangle, c Color) float64 {
	im := copyRGBA(current)
	t.Draw(im, c)
	return DifferenceScore(target, im)
}

func Run(im image.Image) {
	size := im.Bounds().Size()
	w, h := size.X, size.Y

	target := imageToRGBA(im)
	current := uniformRGBA(im.Bounds(), color.Black)

	frame := 0
	for {
		bestScore := 1.0
		bestTriangle := Triangle{}
		bestColor := Color{}
		for i := 0; i < 1000; i++ {
			t := RandomTriangle(w, h)
			c := computeColor(target, current, t, 128)
			s := score(target, current, t, c)
			if s < bestScore {
				bestScore = s
				bestTriangle = t
				bestColor = c
			}
		}
		fmt.Println(frame, bestScore)
		bestTriangle.Draw(current, bestColor)
		if frame%1 == 0 {
			path := fmt.Sprintf("out%03d.png", frame)
			SavePNG(path, current)
		}
		frame++
	}
}
