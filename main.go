package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fogleman/primitive/primitive"
	"github.com/nfnt/resize"
)

var (
	Input      string
	Outputs    flagArray
	Background string
	Number     int
	Alpha      int
	InputSize  int
	OutputSize int
	Mode       int
	Workers    int
	V, VV      bool
)

type flagArray []string

func (i *flagArray) String() string {
	return strings.Join(*i, ", ")
}

func (i *flagArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func init() {
	flag.StringVar(&Input, "i", "", "input image path")
	flag.Var(&Outputs, "o", "output image path")
	flag.StringVar(&Background, "bg", "", "background color (hex)")
	flag.IntVar(&Number, "n", 0, "number of primitives")
	flag.IntVar(&Alpha, "a", 128, "alpha value")
	flag.IntVar(&InputSize, "r", 256, "resize large input images to this size")
	flag.IntVar(&OutputSize, "s", 1024, "output image size")
	flag.IntVar(&Mode, "m", 1, "0=combo 1=triangle 2=rect 3=ellipse 4=circle 5=rotatedrect")
	flag.IntVar(&Workers, "j", 0, "number of parallel workers (default uses all cores)")
	flag.BoolVar(&V, "v", false, "verbose")
	flag.BoolVar(&VV, "vv", false, "very verbose")
}

func errorMessage(message string) bool {
	fmt.Fprintln(os.Stderr, message)
	return false
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// parse and validate arguments
	flag.Parse()
	ok := true
	if Input == "" {
		ok = errorMessage("ERROR: input argument required")
	}
	if len(Outputs) == 0 {
		ok = errorMessage("ERROR: output argument required")
	}
	if Number == 0 {
		ok = errorMessage("ERROR: number argument required")
	}
	if !ok {
		fmt.Println("Usage: primitive [OPTIONS] -i input -o output -n shape_count")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// set log level
	if V {
		primitive.LogLevel = 1
	}
	if VV {
		primitive.LogLevel = 2
	}

	// seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	// read input image
	primitive.Log(1, "reading %s\n", Input)
	input, err := primitive.LoadImage(Input)
	check(err)

	// scale down input image if needed
	size := uint(InputSize)
	input = resize.Thumbnail(size, size, input, resize.Bilinear)

	// determine background color
	var bg primitive.Color
	if Background == "" {
		bg = primitive.MakeColor(primitive.AverageImageColor(input))
	} else {
		bg = primitive.MakeHexColor(Background)
	}

	// run algorithm
	model := primitive.NewModel(input, bg, OutputSize)
	primitive.Log(1, "%d: t=%.3f, score=%.6f\n", 0, 0.0, model.Score)
	start := time.Now()
	for i := 1; i <= Number; i++ {
		// find optimal shape and add it to the model
		t := time.Now()
		n := model.Step(primitive.ShapeType(Mode), Alpha, Workers)
		nps := primitive.NumberString(float64(n) / time.Since(t).Seconds())
		elapsed := time.Since(start).Seconds()
		primitive.Log(1, "%d: t=%.3f, score=%.6f, n=%d, n/s=%s\n", i, elapsed, model.Score, n, nps)

		// write output image(s)
		if i%10 == 0 {
			primitive.SavePNG("out.png", model.Current)
		}
		for _, output := range Outputs {
			ext := strings.ToLower(filepath.Ext(output))
			saveFrames := strings.Contains(output, "%") && ext != ".gif"
			if saveFrames || i == Number {
				path := output
				if saveFrames {
					path = fmt.Sprintf(output, i)
				}
				primitive.Log(1, "writing %s\n", path)
				switch ext {
				default:
					check(fmt.Errorf("unrecognized file extension: %s", ext))
				case ".png":
					check(primitive.SavePNG(path, model.Context.Image()))
				case ".jpg", ".jpeg":
					check(primitive.SaveJPG(path, model.Context.Image(), 95))
				case ".svg":
					check(primitive.SaveFile(path, model.SVG()))
				case ".gif":
					frames := model.Frames(0.001)
					check(primitive.SaveGIFImageMagick(path, frames, 50, 250))
				}
			}
		}
	}
}
