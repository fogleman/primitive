package primitive

import (
	"image"
	"image/color"
	"math"
)

// Heatmap models a heatmap as a slice of uint64 with dimensions
type Heatmap struct {
	W, H  int
	Count []uint64
}

// NewHeatmap models a heatmap with a specific width and height
func NewHeatmap(w, h int) *Heatmap {
	count := make([]uint64, w*h)
	return &Heatmap{w, h, count}
}

// Clear sets all the values in the heatmap to 0
func (h *Heatmap) Clear() {
	for i := range h.Count {
		h.Count[i] = 0
	}
}

// Add adds a slice of scanlines to the current heatmap
func (h *Heatmap) Add(lines []Scanline) {
	for _, line := range lines {
		i := line.Y*h.W + line.X1
		for x := line.X1; x <= line.X2; x++ {
			h.Count[i] += uint64(line.Alpha)
			i++
		}
	}
}

// AddHeatmap adds the values of one heatmap to another heatmap
func (h *Heatmap) AddHeatmap(a *Heatmap) {
	for i, x := range a.Count {
		h.Count[i] += x
	}
}

// Image returns an image which was derived from the current heatmap
func (h *Heatmap) Image(gamma float64) *image.Gray16 {
	im := image.NewGray16(image.Rect(0, 0, h.W, h.H))
	var hi uint64
	for _, h := range h.Count {
		if h > hi {
			hi = h
		}
	}
	i := 0
	for y := 0; y < h.H; y++ {
		for x := 0; x < h.W; x++ {
			p := float64(h.Count[i]) / float64(hi)
			p = math.Pow(p, gamma)
			im.SetGray16(x, y, color.Gray16{uint16(p * 0xffff)})
			i++
		}
	}
	return im
}
