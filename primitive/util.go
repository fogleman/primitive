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
	"io"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// Interface types for mocking os.File
type closableWriter interface {
	Write(p []byte) (n int, err error)
	WriteString(s string) (n int, err error)
	Close() error
}

type closableReader interface {
	Read(p []byte) (n int, err error)
	Close() error
}

type runable interface {
	Run() error
}

// alias downstream functions to enable mocking for unit test
var osStdin = func() io.Reader { return os.Stdin }
var osStdout = func() io.Writer { return os.Stdout }
var osOpen = func(p string) (closableReader, error) { return os.Open(p) }
var osCreate = func(p string) (closableWriter, error) { return os.Create(p) }
var execCommand = func(cmd string, arg ...string) runable { return exec.Command(cmd, arg...) }
var pngEncode = png.Encode

var osRemoveAll = os.RemoveAll
var imageDecode = image.Decode
var fmtFprint = fmt.Fprint
var jpegEncode = jpeg.Encode
var imageNewPaletted = image.NewPaletted
var drawDraw = draw.Draw

var gifEncodeAll = gif.EncodeAll
var imageNewRGBA = image.NewRGBA
var ioutilTempDir = ioutil.TempDir

// LoadImage decodes an image from an image file into raw format
func LoadImage(path string) (image.Image, error) {
	if path == "-" {
		im, _, err := imageDecode(osStdin())
		return im, err
	}
	file, err := osOpen(path)
	if err != nil {
		return nil, err
	} else {
		defer file.Close()
		im, _, err := imageDecode(file)
		return im, err
	}
}

// SaveFile saves the value of the 'contents' string in the specified path
func SaveFile(path, contents string) error {
	if path == "-" {
		_, err := fmtFprint(osStdout(), contents)
		return err
	}
	file, err := osCreate(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(contents)
	return err
}

// SavePNG saves the passed in image as a png
func SavePNG(path string, im image.Image) error {
	file, err := osCreate(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return pngEncode(file, im)
}

// SaveJPG saves the passed in image as a jpeg
func SaveJPG(path string, im image.Image, quality int) error {
	file, err := osCreate(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return jpegEncode(file, im, &jpeg.Options{Quality: quality})
}

// SaveGIF saves the passed in images slice as a gif
func SaveGIF(path string, frames []image.Image, delay, lastDelay int) error {
	g := gif.GIF{}
	for i, src := range frames {
		dst := imageNewPaletted(src.Bounds(), palette.Plan9)
		drawDraw(dst, dst.Rect, src, image.ZP, draw.Src)
		g.Image = append(g.Image, dst)
		if i == len(frames)-1 {
			g.Delay = append(g.Delay, lastDelay)
		} else {
			g.Delay = append(g.Delay, delay)
		}
	}
	file, err := osCreate(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return gifEncodeAll(file, &g)
}

// SaveGIFImageMagick saves the passed in slice of images
func SaveGIFImageMagick(path string, frames []image.Image, delay, lastDelay int) error {
	dir, err := ioutilTempDir("", "")
	if err != nil {
		return err
	}
	for i, im := range frames {
		path := filepath.Join(dir, fmt.Sprintf("%06d.png", i))
		err = SavePNG(path, im)
		if err != nil {
			return err
		}
	}
	args := []string{
		"-loop", "0",
		"-delay", fmt.Sprint(delay),
		filepath.Join(dir, "*.png"),
		"-delay", fmt.Sprint(lastDelay - delay),
		filepath.Join(dir, fmt.Sprintf("%06d.png", len(frames)-1)),
		path,
	}
	cmd := execCommand("convert", args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return osRemoveAll(dir)
}

// NumberString translates a float64 into a filesize string
func NumberString(bf float64) string {
	for _, unit := range []string{"", "K", "M", "G", "T", "P", "E", "Z"} {
		if math.Abs(bf) < 1000.0 {
			formatted := strconv.FormatFloat(bf, 'f', 2, 64) + unit + "B"
			if len(formatted) >= 9 {
				//Seeing it say 1000.00KB instead of 1.00MB was making me itch
				bf /= 1000
				continue
			}
			return (formatted)
		}
		bf /= 1000
	}
	return strconv.FormatFloat(bf, 'f', 2, 64) + "YB"
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
	dst := imageNewRGBA(src.Bounds())
	drawDraw(dst, dst.Rect, src, image.ZP, draw.Src)
	return dst
}

func copyRGBA(src *image.RGBA) *image.RGBA {
	dst := imageNewRGBA(src.Bounds())
	copy(dst.Pix, src.Pix)
	return dst
}

func uniformRGBA(r image.Rectangle, c color.Color) *image.RGBA {
	im := imageNewRGBA(r)
	drawDraw(im, im.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)
	return im
}

// AverageImageColor Takes the average of all the colors in an
// image to create a good default background color
func AverageImageColor(im image.Image) color.NRGBA {
	rgba := imageToRGBA(im)
	size := rgba.Bounds().Size()
	w, h := size.X, size.Y
	var r, g, b int

	// Scan through every pixel at every x,y location and add up their
	// total r, g, and b values
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := rgba.RGBAAt(x, y)
			r += int(c.R)
			g += int(c.G)
			b += int(c.B)
		}
	}
	// Divide r, g, and b by the total number of pixles to get the
	// average r, g, and b
	r /= w * h
	g /= w * h
	b /= w * h
	//Return the averages in NRGAB format with an alpha of 255
	return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
}
