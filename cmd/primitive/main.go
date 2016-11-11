package main

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/primitive/primitive"
	"github.com/nfnt/resize"
)

var InvalidCommand = errors.New("invalid command")

type Config struct {
	Model       *primitive.Model
	BigModel    *primitive.Model
	Size        int
	BigSize     int
	Background  primitive.Color
	Image       image.Image
	Shape       primitive.ShapeType
	Alpha       int
	Repeat      int
	Workers     int
	StrokeWidth float64
	Scale       float64
	Dirty       bool
	Timestamp   time.Time
}

func NewConfig() *Config {
	c := &Config{}
	c.Alpha = 128
	c.Background = primitive.Color{}
	c.Shape = primitive.ShapeTypeTriangle
	c.Size = 256
	c.BigSize = 512
	c.StrokeWidth = 1
	c.Dirty = true
	c.Timestamp = time.Now()
	return c
}

func (c *Config) Step() {
	if c.Dirty {
		image := c.Image
		if c.Size > 0 {
			size := uint(c.Size)
			image = resize.Thumbnail(size, size, image, resize.Bilinear)
		}
		bigImage := c.Image
		if c.BigSize > 0 {
			size := uint(c.BigSize)
			bigImage = resize.Thumbnail(size, size, bigImage, resize.Bilinear)
		}
		c.Scale = (float64(bigImage.Bounds().Size().X) /
			float64(image.Bounds().Size().X))
		workers := c.Workers
		if workers < 1 {
			workers = runtime.NumCPU()
		}
		background := c.Background
		blank := primitive.Color{}
		if background == blank {
			background = primitive.MakeColor(primitive.AverageImageColor(image))
		}
		c.Model = primitive.NewModel(image, background, 1024, workers)
		c.BigModel = primitive.NewModel(bigImage, background, 1024, workers)
		c.Dirty = false
		size := bigImage.Bounds().Size()
		println(fmt.Sprintf("size %d %d", size.X, size.Y))
		println(fmt.Sprintf("background %d %d %d %d",
			background.R, background.G, background.B, background.A))
	}
	c.Model.StrokeWidth = c.StrokeWidth
	c.BigModel.StrokeWidth = c.StrokeWidth * c.Scale
	for i := 0; i <= c.Repeat; i++ {
		if i == 0 {
			state := c.Model.GlobalSearch(c.Shape, c.Alpha)
			state, _ = c.BigModel.LocalSearch(state.Shape.Scale(c.Scale), c.Alpha)
			c.BigModel.AddState(state, 1)
			c.Model.AddState(state, 1/c.Scale)
		} else {
			shape := c.BigModel.Shapes[len(c.BigModel.Shapes)-1]
			state, ok := c.BigModel.LocalSearch(shape, c.Alpha)
			if !ok {
				break
			}
			c.BigModel.AddState(state, 1)
			c.Model.AddState(state, 1/c.Scale)
		}
		color := c.BigModel.Colors[len(c.BigModel.Colors)-1]
		shape := c.BigModel.Shapes[len(c.BigModel.Shapes)-1]
		println(fmt.Sprintf("color %d %d %d %d",
			color.R, color.G, color.B, color.A))
		println(shape.Command())
	}
	println(fmt.Sprintf("score %f", c.BigModel.Score))
	// primitive.SavePNG("1.png", c.Model.Current)
	// primitive.SavePNG("2.png", c.BigModel.Current)
}

func (c *Config) Run(n int) {
	for i := 0; i < n; i++ {
		c.Step()
	}
}

func (c *Config) ParseLine(line string) error {
	line = strings.TrimSpace(line)
	args := strings.Split(line, " ")
	if len(args) == 0 {
		return InvalidCommand
	}
	command, args := strings.ToLower(args[0]), args[1:]
	remainder := strings.TrimSpace(line[len(command):])
	switch command {
	case "keepalive":
		return nil
	case "image":
		return c.parseImage(remainder)
	case "shape":
		return c.parseShape(args)
	case "size":
		return c.parseSize(args)
	case "alpha":
		return c.parseAlpha(args)
	case "repeat":
		return c.parseRepeat(args)
	case "workers":
		return c.parseWorkers(args)
	case "background":
		return c.parseBackground(args)
	case "clear":
		return c.parseClear(args)
	case "run":
		return c.parseRun(args)
	case "step":
		return c.parseStep(args)
	case "save":
		return c.parseSave(args)
	case "strokewidth":
		return c.parseStrokeWidth(args)
	}
	return InvalidCommand
}

func (c *Config) parseInt(args []string, min, max int) (int, error) {
	if len(args) != 1 {
		return 0, InvalidCommand
	}
	x, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, err
	}
	if x < min || x > max {
		return 0, InvalidCommand
	}
	return x, nil
}

func (c *Config) parseFloat(args []string, min, max float64) (float64, error) {
	if len(args) != 1 {
		return 0, InvalidCommand
	}
	x, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return 0, err
	}
	if x < min || x > max {
		return 0, InvalidCommand
	}
	return x, nil
}

func (c *Config) parseImage(path string) error {
	im, err := primitive.LoadImage(path)
	c.Image = im
	c.Dirty = true
	return err
}

func (c *Config) parseShape(args []string) error {
	if len(args) != 1 {
		return InvalidCommand
	}
	switch args[0] {
	case "0", "any":
		c.Shape = primitive.ShapeTypeAny
	case "1", "triangle":
		c.Shape = primitive.ShapeTypeTriangle
	case "2", "rectangle":
		c.Shape = primitive.ShapeTypeRectangle
	case "3", "ellipse":
		c.Shape = primitive.ShapeTypeEllipse
	case "4", "circle":
		c.Shape = primitive.ShapeTypeCircle
	case "5", "rotatedrectangle":
		c.Shape = primitive.ShapeTypeRotatedRectangle
	case "6", "quadratic":
		c.Shape = primitive.ShapeTypeQuadratic
	case "7", "rotatedellipse":
		c.Shape = primitive.ShapeTypeRotatedEllipse
	case "8", "polygon":
		c.Shape = primitive.ShapeTypePolygon
	default:
		return InvalidCommand
	}
	return nil
}

func (c *Config) parseStrokeWidth(args []string) error {
	strokeWidth, err := c.parseFloat(args, 0, math.MaxFloat64)
	if err != nil {
		return err
	}
	c.StrokeWidth = strokeWidth
	return nil
}

func (c *Config) parseSize(args []string) error {
	if len(args) != 2 {
		return InvalidCommand
	}
	size, err := c.parseInt(args[:1], 0, math.MaxInt32)
	if err != nil {
		return err
	}
	bigSize, err := c.parseInt(args[1:], 0, math.MaxInt32)
	if err != nil {
		return err
	}
	c.Size = size
	c.BigSize = bigSize
	c.Dirty = true
	return nil
}

func (c *Config) parseAlpha(args []string) error {
	alpha, err := c.parseInt(args, 0, 255)
	if err != nil {
		return err
	}
	c.Alpha = alpha
	return nil
}

func (c *Config) parseRepeat(args []string) error {
	repeat, err := c.parseInt(args, 0, math.MaxInt32)
	if err != nil {
		return err
	}
	c.Repeat = repeat
	return nil
}

func (c *Config) parseWorkers(args []string) error {
	workers, err := c.parseInt(args, 0, math.MaxInt32)
	if err != nil {
		return err
	}
	c.Workers = workers
	c.Dirty = true
	return nil
}

func (c *Config) parseBackground(args []string) error {
	if len(args) != 1 {
		return InvalidCommand
	}
	c.Background = primitive.MakeHexColor(args[0])
	c.Dirty = true
	return nil
}

func (c *Config) parseClear(args []string) error {
	if len(args) != 0 {
		return InvalidCommand
	}
	c.Dirty = true
	return nil
}

func (c *Config) parseRun(args []string) error {
	if c.Image == nil {
		return InvalidCommand
	}
	n, err := c.parseInt(args, 0, math.MaxInt32)
	if err != nil {
		return err
	}
	c.Run(n)
	return nil
}

func (c *Config) parseStep(args []string) error {
	if len(args) != 0 || c.Image == nil {
		return InvalidCommand
	}
	c.Step()
	return nil
}

func (c *Config) parseSave(args []string) error {
	if len(args) != 1 || c.Model == nil {
		return InvalidCommand
	}
	path := args[0]
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png":
		return primitive.SavePNG(path, c.Model.Context.Image())
	case ".jpg", ".jpeg":
		return primitive.SaveJPG(path, c.Model.Context.Image(), 95)
	case ".svg":
		return primitive.SaveFile(path, c.Model.SVG())
	}
	return InvalidCommand
}

func readLine(reader *bufio.Reader) (string, error) {
	result := ""
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			return result, err
		}
		result += string(line)
		if !isPrefix {
			return result, nil
		}
	}
}

func println(x string) (int, error) {
	return os.Stdout.Write([]byte(x + "\n"))
}

func watchdog(config *Config) {
	for {
		time.Sleep(time.Second * 5)
		age := time.Since(config.Timestamp)
		if age > time.Second*60 {
			os.Exit(1)
		}
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	config := NewConfig()
	go watchdog(config)
	reader := bufio.NewReaderSize(os.Stdin, 65536)
	for {
		line, err := readLine(reader)
		if err != nil && err != io.EOF {
			break
		}
		if len(line) == 0 {
			break
		}
		config.Timestamp = time.Now()
		if err := config.ParseLine(line); err != nil {
			println("err")
		} else {
			println("ok")
		}
	}
}
