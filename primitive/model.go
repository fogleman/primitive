package primitive

import (
	"fmt"
	"image"
	"strings"

	"github.com/fogleman/gg"
)

type Model struct {
	Sw, Sh, Vw, Vh int
	Scale          float64
	Background     Color
	Target         *image.RGBA
	Current        *image.RGBA
	Context        *gg.Context
	Score          float64
	Shapes         []Shape
	Colors         []Color
	Scores         []float64
	Workers        []*Worker
}

func NewModel(target image.Image, background Color, size, numWorkers int) *Model {
	w := target.Bounds().Size().X
	h := target.Bounds().Size().Y
	aspect := float64(w) / float64(h)
	var sw, sh int
	var scale float64
	if aspect >= 1 {
		sw = size
		sh = int(float64(size) / aspect)
		scale = float64(size) / float64(w)
	} else {
		sw = int(float64(size) * aspect)
		sh = size
		scale = float64(size) / float64(h)
	}

	model := &Model{}
	model.Sw = sw
	model.Sh = sh
	model.Vw = sw / int(scale)
	model.Vh = sh / int(scale)
	model.Scale = scale
	model.Background = background
	model.Target = imageToRGBA(target)
	model.Current = uniformRGBA(target.Bounds(), background.NRGBA())
	model.Score = differenceFull(model.Target, model.Current)
	model.Context = model.newContext()
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(model.Target)
		model.Workers = append(model.Workers, worker)
	}
	return model
}

func (model *Model) newContext() *gg.Context {
	dc := gg.NewContext(model.Sw, model.Sh)
	dc.Scale(model.Scale, model.Scale)
	dc.Translate(0.5, 0.5)
	dc.SetColor(model.Background.NRGBA())
	dc.Clear()
	return dc
}

func (model *Model) Frames(scoreDelta float64) []image.Image {
	var result []image.Image
	dc := model.newContext()
	result = append(result, imageToRGBA(dc.Image()))
	previous := 10.0
	for i, shape := range model.Shapes {
		c := model.Colors[i]
		dc.SetRGBA255(c.R, c.G, c.B, c.A)
		shape.Draw(dc, model.Scale)
		dc.Fill()
		score := model.Scores[i]
		delta := previous - score
		if delta >= scoreDelta {
			previous = score
			result = append(result, imageToRGBA(dc.Image()))
		}
	}
	return result
}

func (model *Model) SVG() string {
	bg := model.Background
	fillA := model.Colors[0].A
	var lines []string
	lines = append(lines, fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" version=\"1.1\" viewBox=\"0 0 %d %d\">", model.Vw, model.Vh))
	lines = append(lines, fmt.Sprintf("<rect width=\"100%%\" height=\"100%%\" fill=\"#%02x%02x%02x\" />", bg.R, bg.G, bg.B))
	group := "<g fill-opacity=\"%f\">"
	lines = append(lines, fmt.Sprintf(group, float64(fillA)/255))
	for i, shape := range model.Shapes {
		var attrs []string
		c := model.Colors[i]
		fill := "fill=\"#%02x%02x%02x\""
		fill = fmt.Sprintf(fill, c.R, c.G, c.B)
		attrs = append(attrs, fill)
		if c.A != fillA {
			opacity := "fill-opacity=\"%f\""
			attrs = append(attrs, fmt.Sprintf(opacity, float64(c.A)/255))
		}
		lines = append(lines, shape.SVG(strings.Join(attrs, " ")))
	}
	lines = append(lines, "</g>")
	lines = append(lines, "</svg>")
	return strings.Join(lines, "\n")
}

func (model *Model) Add(shape Shape, alpha int) {
	before := copyRGBA(model.Current)
	lines := shape.Rasterize()
	color := computeColor(model.Target, model.Current, lines, alpha)
	drawLines(model.Current, color, lines)
	score := differencePartial(model.Target, before, model.Current, model.Score, lines)

	model.Score = score
	model.Shapes = append(model.Shapes, shape)
	model.Colors = append(model.Colors, color)
	model.Scores = append(model.Scores, score)

	model.Context.SetRGBA255(color.R, color.G, color.B, color.A)
	shape.Draw(model.Context, model.Scale)
}

func (model *Model) Step(shapeType ShapeType, alpha, repeat int) int {
	state := model.runWorkers(shapeType, alpha, 1000, 100, 16)
	// state = HillClimb(state, 1000).(*State)
	model.Add(state.Shape, state.Alpha)

	for i := 0; i < repeat; i++ {
		state.Worker.Init(model.Current, model.Score)
		a := state.Energy()
		state = HillClimb(state, 100).(*State)
		b := state.Energy()
		if a == b {
			break
		}
		model.Add(state.Shape, state.Alpha)
	}

	// for _, w := range model.Workers[1:] {
	// 	model.Workers[0].Heatmap.AddHeatmap(w.Heatmap)
	// }
	// SavePNG("heatmap.png", model.Workers[0].Heatmap.Image(0.5))

	counter := 0
	for _, worker := range model.Workers {
		counter += worker.Counter
	}
	return counter
}

func (model *Model) runWorkers(t ShapeType, a, n, age, m int) *State {
	wn := len(model.Workers)
	ch := make(chan *State, wn)
	wm := m / wn
	if m%wn != 0 {
		wm++
	}
	for i := 0; i < wn; i++ {
		worker := model.Workers[i]
		worker.Init(model.Current, model.Score)
		go model.runWorker(worker, t, a, n, age, wm, ch)
	}
	var bestEnergy float64
	var bestState *State
	for i := 0; i < wn; i++ {
		state := <-ch
		energy := state.Energy()
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

func (model *Model) runWorker(worker *Worker, t ShapeType, a, n, age, m int, ch chan *State) {
	ch <- worker.BestHillClimbState(t, a, n, age, m)
}
