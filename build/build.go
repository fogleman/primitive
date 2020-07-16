package build

import (
	"image"
	"runtime"
	"time"

	"github.com/fogleman/primitive/primitive"
	"github.com/nfnt/resize"
)

// Config are options needed to build a primitive image
type Config struct {
	Input      image.Image
	Background string
	Alpha      int
	Count      int
	OutputSize int
	Mode       int
	Repeat     int
	V, VV      bool
}

// Build builds a primitive image from available options
func Build(conf Config, stepCallback func(*primitive.Model, int, int)) *primitive.Model {
	start := time.Now()

	// set log level
	if conf.V {
		primitive.LogLevel = 1
	}
	if conf.VV {
		primitive.LogLevel = 2
	}

	// determine background color
	var bg primitive.Color
	if conf.Background == "" {
		bg = primitive.MakeColor(primitive.AverageImageColor(conf.Input))
	} else {
		bg = primitive.MakeHexColor(conf.Background)
	}

	input := resize.Thumbnail(uint(256), uint(256), conf.Input, resize.Bilinear)
	model := primitive.NewModel(input, bg, conf.OutputSize, runtime.NumCPU())
	primitive.Log(1, "%d: t=%.3f, score=%.6f\n", 0, 0.0, model.Score)

	frame := 0
	for i := 0; i < conf.Count; i++ {
		frame++

		// find optimal shape and add it to the model
		t := time.Now()
		n := model.Step(primitive.ShapeType(conf.Mode), conf.Alpha, conf.Repeat)
		nps := primitive.NumberString(float64(n) / time.Since(t).Seconds())
		elapsed := time.Since(start).Seconds()
		primitive.Log(1, "%d: t=%.3f, score=%.6f, n=%d, n/s=%s\n", frame, elapsed, model.Score, n, nps)
		stepCallback(model, frame, i)
	}

	return model
}
