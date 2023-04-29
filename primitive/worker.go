package primitive

import (
	"image"
	"math/rand"
	"time"

	"github.com/golang/freetype/raster"
)

// Worker models a worker thread for image transforms
type Worker struct {
	W, H       int
	Target     *image.RGBA
	Current    *image.RGBA
	Buffer     *image.RGBA
	Rasterizer *raster.Rasterizer
	Lines      []Scanline
	Heatmap    *Heatmap
	Rnd        *rand.Rand
	Score      float64
	Counter    int
}

// NewWorker creates a worker
func NewWorker(target *image.RGBA) *Worker {
	w := target.Bounds().Size().X
	h := target.Bounds().Size().Y
	worker := Worker{}
	worker.W = w
	worker.H = h
	worker.Target = target
	worker.Buffer = image.NewRGBA(target.Bounds())
	worker.Rasterizer = raster.NewRasterizer(w, h)
	worker.Lines = make([]Scanline, 0, 4096) // TODO: based on height
	worker.Heatmap = NewHeatmap(w, h)
	worker.Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	return &worker
}

// Init sets default values for a worker and associates it with an image
func (worker *Worker) Init(current *image.RGBA, score float64) {
	worker.Current = current
	worker.Score = score
	worker.Counter = 0
	worker.Heatmap.Clear()
}

// Energy returns the energy calculation for the worker's shape
func (worker *Worker) Energy(shape Shape, alpha int) float64 {
	worker.Counter++
	lines := shape.Rasterize()
	// worker.Heatmap.Add(lines)
	color := computeColor(worker.Target, worker.Current, lines, alpha)
	copyLines(worker.Buffer, worker.Current, lines)
	notify := NewTestStringNotifier()
	drawLines(worker.Buffer, color, lines, notify)
	return differencePartial(worker.Target, worker.Current, worker.Buffer, worker.Score, lines)
}

// BestHillClimbState returns the best solution the worker has found
func (worker *Worker) BestHillClimbState(
	t ShapeType, alpha, triesPerWorker, age, climbes int) *State {
	var bestEnergy float64
	var bestState *State
	for i := 0; i < climbes; i++ {
		testState := worker.BestRandomState(t, alpha, triesPerWorker)
		beforeEnergy := testState.Energy()
		testState = HillClimb(testState, age).(*State)
		climedEnergy := testState.Energy()
		vv("%dx random: %.6f -> %dx hill climb: %.6f\n", triesPerWorker, beforeEnergy, age, climedEnergy)
		if i == 0 || climedEnergy < bestEnergy {
			bestEnergy = climedEnergy
			bestState = testState
		}
	}
	return bestState
}

// BestRandomState returns the best solution found based on a random
// optimization
func (worker *Worker) BestRandomState(t ShapeType, alpha, triesPerWorker int) *State {
	var bestEnergy float64
	var bestState *State
	for i := 0; i < triesPerWorker; i++ {
		testState := worker.RandomState(t, alpha)
		energy := testState.Energy()
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = testState
		}
	}
	return bestState
}

// RandomState returns a random shape wrapped in a state
func (worker *Worker) RandomState(t ShapeType, alpha int) *State {
	switch t {
	default:
		return worker.RandomState(ShapeType(worker.Rnd.Intn(8)+1), alpha)
	case shapeTypeTriangle:
		return NewState(worker, NewRandomTriangle(worker), alpha)
	case shapeTypeRectangle:
		return NewState(worker, NewRandomRectangle(worker), alpha)
	case shapeTypeEllipse:
		return NewState(worker, NewRandomEllipse(worker), alpha)
	case shapeTypeCircle:
		return NewState(worker, NewRandomCircle(worker), alpha)
	case shapeTypeRotatedRectangle:
		return NewState(worker, NewRandomRotatedRectangle(worker), alpha)
	case shapeTypeQuadratic:
		return NewState(worker, NewRandomQuadratic(worker), alpha)
	case shapeTypeRotatedEllipse:
		return NewState(worker, NewRandomRotatedEllipse(worker), alpha)
	case shapeTypePoligon:
		return NewState(worker, NewRandomPolygon(worker, 4, false), alpha)
	}
}
