package primitive

type Scanline struct {
	Y, X1, X2 int
	Alpha     uint32
}

func cropScanlines(lines []Scanline, w, h int) []Scanline {
	result := make([]Scanline, 0, len(lines))
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
		result = append(result, line)
	}
	return result
}
