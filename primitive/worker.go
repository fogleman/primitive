package primitive

import (
	"image"
	"math/rand"
	"time"

	"github.com/golang/freetype/raster"
)

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

func (worker *Worker) Init(current *image.RGBA, score float64) {
	worker.Current = current
	worker.Score = score
	worker.Counter = 0
	worker.Heatmap.Clear()
}

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

func (worker *Worker) RandomState(t ShapeType, alpha int) *State {
	switch t {
	default:
		return worker.RandomState(ShapeType(worker.Rnd.Intn(8)+1), alpha)
	case ShapeTypeTriangle:
		return NewState(worker, NewRandomTriangle(worker), alpha)
	case ShapeTypeRectangle:
		return NewState(worker, NewRandomRectangle(worker), alpha)
	case ShapeTypeEllipse:
		return NewState(worker, NewRandomEllipse(worker), alpha)
	case ShapeTypeCircle:
		return NewState(worker, NewRandomCircle(worker), alpha)
	case ShapeTypeRotatedRectangle:
		return NewState(worker, NewRandomRotatedRectangle(worker), alpha)
	case ShapeTypeQuadratic:
		return NewState(worker, NewRandomQuadratic(worker), alpha)
	case ShapeTypeRotatedEllipse:
		return NewState(worker, NewRandomRotatedEllipse(worker), alpha)
	case ShapeTypePolygon:
		return NewState(worker, NewRandomPolygon(worker, 4, false), alpha)
	}
}
