package primitive

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path/filepath"
)

func LoadImage(path string) (image.Image, error) {
	if path == "-" {
		im, _, err := image.Decode(os.Stdin)
		return im, err
	} else {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		im, _, err := image.Decode(file)
		return im, err
	}
}

func SaveFile(path, contents string) error {
	if path == "-" {
		_, err := fmt.Fprint(os.Stdout, contents)
		return err
	} else {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = file.WriteString(contents)
		return err
	}
}

func SavePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, im)
}

func SaveJPG(path string, im image.Image, quality int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return jpeg.Encode(file, im, &jpeg.Options{quality})
}

func SaveGIF(path string, frames []image.Image, delay, lastDelay int) error {
	g := gif.GIF{}
	for i, src := range frames {
		dst := image.NewPaletted(src.Bounds(), palette.Plan9)
		draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)
		g.Image = append(g.Image, dst)
		if i == len(frames)-1 {
			g.Delay = append(g.Delay, lastDelay)
		} else {
			g.Delay = append(g.Delay, delay)
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return gif.EncodeAll(file, &g)
}

func SaveGIFImageMagick(path string, frames []image.Image, delay, lastDelay int) error {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	for i, im := range frames {
		path := filepath.Join(dir, fmt.Sprintf("%06d.png", i))
		SavePNG(path, im)
	}
	args := []string{
		"-loop", "0",
		"-delay", fmt.Sprint(delay),
		filepath.Join(dir, "*.png"),
		"-delay", fmt.Sprint(lastDelay - delay),
		filepath.Join(dir, fmt.Sprintf("%06d.png", len(frames)-1)),
		path,
	}
	cmd := exec.Command("convert", args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return os.RemoveAll(dir)
}

func NumberString(x float64) string {
	suffixes := []string{"", "k", "M", "G"}
	for _, suffix := range suffixes {
		if x < 1000 {
			return fmt.Sprintf("%.1f%s", x, suffix)
		}
		x /= 1000
	}
	return fmt.Sprintf("%.1f%s", x, "T")
}

func radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func degrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func clampInt(x, lo, hi int) int {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func rotate(x, y, theta float64) (rx, ry float64) {
	rx = x*math.Cos(theta) - y*math.Sin(theta)
	ry = x*math.Sin(theta) + y*math.Cos(theta)
	return
}

func imageToRGBA(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)
	return dst
}

func copyRGBA(src *image.RGBA) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	copy(dst.Pix, src.Pix)
	return dst
}

func uniformRGBA(r image.Rectangle, c color.Color) *image.RGBA {
	im := image.NewRGBA(r)
	draw.Draw(im, im.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)
	return im
}

func AverageImageColor(im image.Image) color.NRGBA {
	rgba := imageToRGBA(im)
	size := rgba.Bounds().Size()
	w, h := size.X, size.Y
	var r, g, b int
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := rgba.RGBAAt(x, y)
			r += int(c.R)
			g += int(c.G)
			b += int(c.B)
		}
	}
	r /= w * h
	g /= w * h
	b /= w * h
	return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
}

//REQ 2.3
func printKey(colorMap *map) {

}
