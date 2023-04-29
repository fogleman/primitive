package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/bensodenkamp/primitive_ent/primitive"
	"github.com/nfnt/resize"
)

// Define command line inputs
var (
	Input      string
	Outputs    flagArray
	Background string
	Configs    shapeConfigArray
	Alpha      int
	InputSize  int
	OutputSize int
	Mode       int
	Workers    int
	Nth        int
	Repeat     int
	V          bool
	VV         bool
)

type flagArray []string

// When interepreted as a string, a flagArray will print out all of its
// Constituent strings as a comma delimited list
func (i *flagArray) String() string {
	return strings.Join(*i, ", ")
}

// Calling 'Set' on a flagArray will add a string to the list.
// The purpose of being able to define multiple outputs is that
// for a given set of primitives, the result can be output as
// multiple filetypes based on extension. i.e. 'primitive -o result.png -o result.jpg'
func (i *flagArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type shapeConfig struct {
	Count  int
	Mode   int
	Alpha  int
	Repeat int
}

type shapeConfigArray []shapeConfig

// Flags must be able to interpret the value of Config as a string
func (i *shapeConfigArray) String() string {
	return ""
}

// Add to the shapeConfigArray 'Config' based on the value passed in at the command line.
// Set the 'count' attribute as the Int parsed from the input value provided and current values of
// Mode, Alpha, and Repeat.
// If multiple flags are set for -n.. i.e. 'primitive -n 100 -n 600', two elements will be generated.
func (i *shapeConfigArray) Set(value string) error {
	n, _ := strconv.ParseInt(value, 0, 0)
	*i = append(*i, shapeConfig{int(n), Mode, Alpha, Repeat})
	return nil
}

// Read in flags from command line
//
// note: command lines are evaluated in order, so for multiple values of n, different values can be
// specified for the other flags i.e. 'primitive -m 1 -n 20 -m 4 -n 100' would yeild a Config with a
// count or 20 and a mode of 1 followed by a count of 100 with a mode of 4.
func init() {
	flag.StringVar(&Input, "i", "", "input image path")
	flag.Var(&Outputs, "o", "output image path")
	flag.Var(&Configs, "n", "number of primitives")
	flag.StringVar(&Background, "bg", "", "background color (hex)")
	flag.IntVar(&Alpha, "a", 128, "alpha value")
	flag.IntVar(&InputSize, "r", 256, "resize large input images to this size")
	flag.IntVar(&OutputSize, "s", 1024, "output image size")
	flag.IntVar(&Mode, "m", 1, "0=combo 1=triangle 2=rect 3=ellipse 4=circle 5=rotatedrect 6=beziers 7=rotatedellipse 8=polygon")
	flag.IntVar(&Workers, "j", 0, "number of parallel workers (default uses all cores)")
	flag.IntVar(&Nth, "nth", 1, "save every Nth frame (put \"%d\" in path)")
	flag.IntVar(&Repeat, "rep", 0, "add N extra shapes per iteration with reduced search")
	flag.BoolVar(&V, "v", false, "verbose")
	flag.BoolVar(&VV, "vv", false, "very verbose")
}

// Send invocation errors to STDERR
func errorMessage(message string) bool {
	fmt.Fprintln(os.Stderr, message)
	return false
}

// Handle fatal errors from downstream code
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// parse and validate arguments
	flag.Parse()
	ok := true

	// Error if no input file is given
	if Input == "" {
		ok = errorMessage("ERROR: input image path required")
	}

	// Error if no output file is given
	if len(Outputs) == 0 {
		ok = errorMessage("ERROR: output image path required")
	}

	// Error if the number of primitives is not defined
	if len(Configs) == 0 {
		ok = errorMessage("ERROR: number of primitives required")
	}

	// If only one config was specified, assign mode, alpha, and repeat as they may have been included
	// as flags after the -n value and thus were not set when the config was created.
	if len(Configs) == 1 {
		Configs[0].Mode = Mode
		Configs[0].Alpha = Alpha
		Configs[0].Repeat = Repeat
	}

	for _, config := range Configs {
		if config.Count < 1 {
			ok = errorMessage("ERROR: number of primitives must be > 0")
		}
	}

	// If there was any error with the command invocation, show the usage
	if !ok {
		fmt.Println("Usage: primitive [OPTIONS] -i input -o output -n count")
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

	// determine worker count
	if Workers < 1 {
		Workers = runtime.NumCPU()
	}

	// read input image
	primitive.Log(1, "reading %s\n", Input)
	inputImage, err := primitive.LoadImage(Input)
	check(err)

	// scale down input image if needed
	size := uint(InputSize)
	if size > 0 {
		inputImage = resize.Thumbnail(size, size, inputImage, resize.Bilinear)
	}

	// determine background color
	var bg primitive.Color
	if Background == "" {
		bg = primitive.MakeColor(primitive.AverageImageColor(inputImage))
	} else {
		bg = primitive.MakeHexColor(Background)
	}

	// run algorithm
	model := primitive.NewModel(inputImage, bg, OutputSize, Workers)
	primitive.Log(1, "%d: t=%.3f, score=%.6f\n", 0, 0.0, model.Score)
	start := time.Now()
	frame := 0
	for j, config := range Configs {
		primitive.Log(1, "count=%d, mode=%d, alpha=%d, repeat=%d\n",
			config.Count, config.Mode, config.Alpha, config.Repeat)

		for i := 0; i < config.Count; i++ {
			frame++

			// find optimal shape and add it to the model
			t := time.Now()
			notify := primitive.NewTestStringNotifier()
			n := model.Step(
				primitive.ShapeType(config.Mode), config.Alpha, config.Repeat, notify)
			nps := primitive.NumberString(float64(n) / time.Since(t).Seconds())
			elapsed := time.Since(start).Seconds()
			primitive.Log(1, "%d: t=%.3f, score=%.6f, n=%d, n/s=%s\n", frame, elapsed, model.Score, n, nps)

			// write output image(s)
			for _, output := range Outputs {
				ext := strings.ToLower(filepath.Ext(output))
				if output == "-" {
					ext = ".svg"
				}
				percent := strings.Contains(output, "%")
				saveFrames := percent && ext != ".gif"
				saveFrames = saveFrames && frame%Nth == 0
				last := j == len(Configs)-1 && i == config.Count-1
				if saveFrames || last {
					path := output
					if percent {
						path = fmt.Sprintf(output, frame)
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
						// In the case of a gif, rather than saving an image, we save
						// a slice of images that are built from the existing model
						frames := model.Frames(0.001, notify)
						check(primitive.SaveGIFImageMagick(path, frames, 50, 250))
					}
				}
			}
		}
	}
}
