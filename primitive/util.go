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

// LoadImage opens the file at the location given by the path argument,
// decodes the file into an image and return it along with an error status response.
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	im, _, err := image.Decode(file)
	return im, err
}

// SaveFile takes a path (save location) and a string which encodes an image
// in SVG format, saving that string encoded image to file.
// Returns an error response, or nil if no error.
func SaveFile(path, contents string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(contents)
	return err
}

// SavePNG takes a path (save location) and a standard lib image.Image type,
// saving that image object as a PNG.
// Returns an error response, or nil if no error.
func SavePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, im)
}

// SaveJPG takes a path (save location) and a standard lib image.Image type,
// saving that image object in JPG encoding.
// Returns an error response, or nil if no error.
func SaveJPG(path string, im image.Image, quality int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return jpeg.Encode(file, im, &jpeg.Options{quality})
}

// Note: This function is currently not used in project.
// SaveGIF, takes a path (save location), a number of image.Image type 'frames',
// and some GIF parameters. It processes all frames according to the GIF parameters
// and saves a GIF encoded file.
// Returns and error response, or nil if no error.
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

// SaveGIF, takes a path (save location), a number of image.Image type 'frames',
// and some GIF parameters. It processes all frames according to the GIF parameters
// and saves a GIF encoded file.
// Returns and error response, or nil if no error.
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

// NumberString returns a human-readable string representation of float x.
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

// radians converts a degrees angle value to a radians angle value.
func radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// degrees converts a radians angle value to a degrees angle value.
func degrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

// clamp takes a float value, x, and a number range, and if x falls outside
// of that range returns the range extremity that is closest to x in value.
func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

// clampInt takes an integer value, x, and a integer range, and if x falls outside
// of that integer range returns the range boundary value that is closest to x.
func clampInt(x, lo, hi int) int {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

// minInt returns the minimum of two integer values.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// maxInt returns the maximum of two integer values.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// rotate performs a rotation of angle theta on a coordinate pair.
// Returns the rotated coordinate pair.
func rotate(x, y, theta float64) (rx, ry float64) {
	rx = x*math.Cos(theta) - y*math.Sin(theta)
	ry = x*math.Sin(theta) + y*math.Cos(theta)
	return
}

// imageToRGBA converts an image.Image (grid of color.Color values) to
// an equivelant image.RGBA (image of color.RGBA values).
func imageToRGBA(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Rect, src, image.ZP, draw.Src)
	return dst
}

// copyRGBA returns a copy of the src image.RGBA.
func copyRGBA(src *image.RGBA) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	copy(dst.Pix, src.Pix)
	return dst
}

// uniformRGBA creates and returns a 'background' which is all one RBGA value color.
func uniformRGBA(r image.Rectangle, c color.Color) *image.RGBA {
	im := image.NewRGBA(r)
	draw.Draw(im, im.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)
	return im
}

// AverageImageColor takes an image.Image object and iterates over
// each pixel value in the image to calculate the average RGBA value,
// which it returns as a non-alpha-premultiplied color.
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
