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
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/primitive/primitive"
	"github.com/nfnt/resize"
)

var InvalidCommand = errors.New("invalid command")

type Config struct {
	Model      *primitive.Model
	Background primitive.Color
	Image      image.Image
	Shape      primitive.ShapeType
	Alpha      int
	Repeat     int
	Workers    int
	Size       int
	Resize     int
	Dirty      bool
}

func NewConfig() *Config {
	c := &Config{}
	c.Alpha = 128
	c.Shape = primitive.ShapeTypeTriangle
	c.Resize = 256
	c.Size = 1024
	c.Dirty = true
	return c
}

func (c *Config) Step() {
	if c.Dirty {
		image := c.Image
		if c.Resize > 0 {
			size := uint(c.Resize)
			image = resize.Thumbnail(size, size, image, resize.Bilinear)
		}
		workers := c.Workers
		if workers < 1 {
			workers = runtime.NumCPU()
		}
		c.Model = primitive.NewModel(image, c.Background, c.Size, workers)
		c.Dirty = false
	}
	index := len(c.Model.Shapes)
	c.Model.Step(c.Shape, c.Alpha, c.Repeat)
	for _, shape := range c.Model.Shapes[index:] {
		fmt.Println(shape.Command())
	}
}

func (c *Config) Run(n int) {
	for i := 0; i < n; i++ {
		c.Step()
	}
}

func (c *Config) ParseLine(line string) error {
	line = strings.TrimSpace(line)
	line = strings.ToLower(line)
	args := strings.Split(line, " ")
	if len(args) == 0 {
		return InvalidCommand
	}
	command, args := args[0], args[1:]
	switch command {
	case "image":
		return c.parseImage(args)
	case "shape":
		return c.parseShape(args)
	case "size":
		return c.parseSize(args)
	case "resize":
		return c.parseResize(args)
	case "alpha":
		return c.parseAlpha(args)
	case "repeat":
		return c.parseRepeat(args)
	case "workers":
		return c.parseWorkers(args)
	case "background":
		return c.parseBackground(args)
	case "run":
		return c.parseRun(args)
	case "step":
		return c.parseStep(args)
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

func (c *Config) parseImage(args []string) error {
	if len(args) != 1 {
		return InvalidCommand
	}
	im, err := primitive.LoadImage(args[0])
	c.Image = im
	c.Dirty = true
	return err
}

func (c *Config) parseShape(args []string) error {
	if len(args) != 1 {
		return InvalidCommand
	}
	switch args[0] {
	case "0":
	case "any":
		c.Shape = primitive.ShapeTypeAny
	case "1":
	case "triangle":
		c.Shape = primitive.ShapeTypeTriangle
	case "2":
	case "rectangle":
		c.Shape = primitive.ShapeTypeRectangle
	case "3":
	case "ellipse":
		c.Shape = primitive.ShapeTypeEllipse
	case "4":
	case "circle":
		c.Shape = primitive.ShapeTypeCircle
	case "5":
	case "rotatedrectangle":
		c.Shape = primitive.ShapeTypeRotatedRectangle
	case "6":
	case "quadratic":
		c.Shape = primitive.ShapeTypeQuadratic
	case "7":
	case "rotatedellipse":
		c.Shape = primitive.ShapeTypeRotatedEllipse
	case "8":
	case "polygon":
		c.Shape = primitive.ShapeTypePolygon
	default:
		return InvalidCommand
	}
	return nil
}

func (c *Config) parseSize(args []string) error {
	size, err := c.parseInt(args, 1, math.MaxInt32)
	if err != nil {
		return err
	}
	c.Size = size
	c.Dirty = true
	return nil
}

func (c *Config) parseResize(args []string) error {
	resize, err := c.parseInt(args, 0, math.MaxInt32)
	if err != nil {
		return err
	}
	c.Resize = resize
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

func (c *Config) parseRun(args []string) error {
	n, err := c.parseInt(args, 0, math.MaxInt32)
	if err != nil {
		return err
	}
	c.Run(n)
	return nil
}

func (c *Config) parseStep(args []string) error {
	if len(args) != 0 {
		return InvalidCommand
	}
	c.Step()
	return nil
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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	config := NewConfig()
	reader := bufio.NewReaderSize(os.Stdin, 65536)
	for {
		line, err := readLine(reader)
		if err != nil && err != io.EOF {
			break
		}
		if len(line) == 0 {
			break
		}
		if err := config.ParseLine(line); err != nil {
			fmt.Println("err")
		} else {
			fmt.Println("ok")
		}
	}
}
