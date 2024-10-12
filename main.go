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

	"github.com/fogleman/primitive/primitive"
	"github.com/nfnt/resize"
)

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
	Seed       int64
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

type shapeConfig struct {
	Count  int
	Mode   int
	Alpha  int
	Repeat int
}

type shapeConfigArray []shapeConfig

func (i *shapeConfigArray) String() string {
	return ""
}

func (i *shapeConfigArray) Set(value string) error {
	n, _ := strconv.ParseInt(value, 0, 0)
	*i = append(*i, shapeConfig{int(n), Mode, Alpha, Repeat})
	return nil
}

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
	flag.Func("seed", "seed to feed the RNG", func(value string) (err error) {
		Seed, err = strconv.ParseInt(value, 10, 64)
		return
	})
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
	if len(Configs) == 0 {
		ok = errorMessage("ERROR: number argument required")
	}
	if len(Configs) == 1 {
		Configs[0].Mode = Mode
		Configs[0].Alpha = Alpha
		Configs[0].Repeat = Repeat
	}
	for _, config := range Configs {
		if config.Count < 1 {
			ok = errorMessage("ERROR: number argument must be > 0")
		}
	}
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
	if Seed == 0 {
		Seed = time.Now().UTC().UnixNano()
	}
	rand.Seed(Seed)

	// determine worker count
	if Workers < 1 {
		Workers = runtime.NumCPU()
	}

	// read input image
	primitive.Log(1, "reading %s\n", Input)
	input, err := primitive.LoadImage(Input)
	check(err)

	// scale down input image if needed
	size := uint(InputSize)
	if size > 0 {
		input = resize.Thumbnail(size, size, input, resize.Bilinear)
	}

	// determine background color
	var bg primitive.Color
	if Background == "" {
		bg = primitive.MakeColor(primitive.AverageImageColor(input))
	} else {
		bg = primitive.MakeHexColor(Background)
	}

	// run algorithm
	model := primitive.NewModel(input, bg, OutputSize, Workers, Seed+1)
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
			n := model.Step(primitive.ShapeType(config.Mode), config.Alpha, config.Repeat)
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
						frames := model.Frames(0.001)
						check(primitive.SaveGIFImageMagick(path, frames, 50, 250))
					}
				}
			}
		}
	}
}
