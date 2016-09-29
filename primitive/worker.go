package primitive

import (
	"image"
	"math/rand"
	"time"

	"github.com/golang/freetype/raster"
)

type Worker struct {
	Target     *image.RGBA
	Current    *image.RGBA
	Buffer     *image.RGBA
	Rasterizer *raster.Rasterizer
	Lines      []Scanline
	Rnd        *rand.Rand
	Score      float64
}

func NewWorker(target, current *image.RGBA) *Worker {
	w := target.Bounds().Size().X
	h := target.Bounds().Size().Y
	worker := Worker{}
	worker.Target = target
	worker.Current = current
	worker.Buffer = image.NewRGBA(target.Bounds())
	worker.Rasterizer = raster.NewRasterizer(w, h)
	worker.Lines = make([]Scanline, 0, 8192)
	worker.Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	worker.Score = differenceFull(target, current)
	return &worker
}

func (worker *Worker) Score(shape Shape, alpha int) float64 {
	lines := shape.Rasterize(worker.Lines)
	color := computeColor(worker.Target, worker.Current, lines, alpha)
	Copy(worker.Buffer, worker.Current, lines)
	Draw(worker.Buffer, color, lines)
	return differencePartial(worker.Target, worker.Current, worker.Buffer, worker.Score, lines)
}
