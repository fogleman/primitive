package primitive

type Scanline struct {
	Y, X1, X2 int
	Alpha     uint32
}

// cropScanLines reduces the scanning width of a set of Scanline
// objects so that they don't scane outside of the area defined by
// w and h.
// Returns a slice containing the original Scanlines, now all cropped 
func cropScanlines(lines []Scanline, w, h int) []Scanline {
	i := 0
	for _, line := range lines {
		if line.Y < 0 || line.Y >= h {
			continue
		}
		if line.X1 >= w {
			continue
		}
		if line.X2 < 0 {
			continue
		}
		line.X1 = clampInt(line.X1, 0, w-1)
		line.X2 = clampInt(line.X2, 0, w-1)
		if line.X1 > line.X2 {
			continue
		}
		lines[i] = line
		i++
	}
	return lines[:i]
}
