package primitive

import (
	"fmt"
	"image"
	"math/rand"
	"runtime"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

type Model struct {
	W, H       int
	Sw, Sh     int
	Scale      float64
	Background Color
	Target     *image.RGBA
	Current    *image.RGBA
	Buffer     *image.RGBA
	Context    *gg.Context
	Score      float64
	Shapes     []Shape
	Colors     []Color
	Scores     []float64
}

func NewModel(target image.Image, background Color, size int) *Model {
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
	model.W = w
	model.H = h
	model.Sw = sw
	model.Sh = sh
	model.Scale = scale
	model.Background = background
	model.Target = imageToRGBA(target)
	model.Current = uniformRGBA(target.Bounds(), background.NRGBA())
	model.Buffer = image.NewRGBA(target.Bounds())
	model.Score = differenceFull(model.Target, model.Current)
	model.Context = model.newContext()
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
	c := model.Background
	var lines []string
	lines = append(lines, fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" version=\"1.1\" width=\"%d\" height=\"%d\">", model.Sw, model.Sh))
	lines = append(lines, fmt.Sprintf("<rect x=\"0\" y=\"0\" width=\"%d\" height=\"%d\" fill=\"#%02x%02x%02x\" />", model.Sw, model.Sh, c.R, c.G, c.B))
	lines = append(lines, fmt.Sprintf("<g transform=\"scale(%f) translate(0.5 0.5)\">", model.Scale))
	for i, shape := range model.Shapes {
		c := model.Colors[i]
		attrs := "fill=\"#%02x%02x%02x\" fill-opacity=\"%f\""
		attrs = fmt.Sprintf(attrs, c.R, c.G, c.B, float64(c.A)/255)
		lines = append(lines, shape.SVG(attrs))
	}
	lines = append(lines, "</g>")
	lines = append(lines, "</svg>")
	return strings.Join(lines, "\n")
}

func (model *Model) Add(shape Shape, alpha int) {
	lines := shape.Rasterize()
	c := model.computeColor(lines, alpha)
	s := model.computeScore(lines, c, model.Buffer)
	Draw(model.Current, c, lines)

	model.Score = s
	model.Shapes = append(model.Shapes, shape)
	model.Colors = append(model.Colors, c)
	model.Scores = append(model.Scores, s)

	model.Context.SetRGBA255(c.R, c.G, c.B, c.A)
	shape.Draw(model.Context, model.Scale)
}

func (model *Model) Step(shapeType ShapeType, alpha, numWorkers int) {
	state := model.runWorkers(shapeType, alpha, numWorkers, 1000, 100, 16)
	state = HillClimb(state, 1000).(*State)
	model.Add(state.Shape, state.Alpha)
}

func (model *Model) runWorkers(t ShapeType, a, wn, n, age, m int) *State {
	if wn < 1 {
		wn = runtime.NumCPU()
	}
	ch := make(chan *State, wn)
	wm := m / wn
	if m%wn != 0 {
		wm++
	}
	for i := 0; i < wn; i++ {
		go model.runWorker(t, a, n, age, wm, ch)
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

func (model *Model) runWorker(t ShapeType, a, n, age, m int, ch chan *State) {
	buffer := image.NewRGBA(model.Target.Bounds())
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	state := model.BestHillClimbState(buffer, t, a, n, age, m, rnd)
	ch <- state
}

func (model *Model) BestHillClimbState(buffer *image.RGBA, t ShapeType, a, n, age, m int, rnd *rand.Rand) *State {
	var bestEnergy float64
	var bestState *State
	for i := 0; i < m; i++ {
		state := model.BestRandomState(buffer, t, a, n, rnd)
		before := state.Energy()
		state = HillClimb(state, age).(*State)
		energy := state.Energy()
		vv("%dx random: %.6f -> %dx hill climb: %.6f\n", n, before, age, energy)
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

func (model *Model) BestRandomState(buffer *image.RGBA, t ShapeType, a, n int, rnd *rand.Rand) *State {
	var bestEnergy float64
	var bestState *State
	for i := 0; i < n; i++ {
		state := model.RandomState(buffer, t, a, rnd)
		energy := state.Energy()
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

func (model *Model) RandomState(buffer *image.RGBA, t ShapeType, a int, rnd *rand.Rand) *State {
	switch t {
	default:
		return model.RandomState(buffer, ShapeType(rnd.Intn(5)+1), a, rnd)
	case ShapeTypeTriangle:
		return NewState(model, buffer, a, NewRandomTriangle(model.W, model.H, rnd), rnd)
	case ShapeTypeRectangle:
		return NewState(model, buffer, a, NewRandomRectangle(model.W, model.H, rnd), rnd)
	case ShapeTypeEllipse:
		return NewState(model, buffer, a, NewRandomEllipse(model.W, model.H, rnd), rnd)
	case ShapeTypeCircle:
		return NewState(model, buffer, a, NewRandomCircle(model.W, model.H, rnd), rnd)
	case ShapeTypeRotatedRectangle:
		return NewState(model, buffer, a, NewRandomRotatedRectangle(model.W, model.H, rnd), rnd)
	case ShapeTypePath:
		return NewState(model, buffer, a, NewRandomPath(model.W, model.H, rnd), rnd)
	}
}

func (model *Model) computeColor(lines []Scanline, alpha int) Color {
	var rsum, gsum, bsum, count int64
	a := 0x101 * 255 / alpha
	for _, line := range lines {
		i := model.Target.PixOffset(line.X1, line.Y)
		for x := line.X1; x <= line.X2; x++ {
			tr := int(model.Target.Pix[i])
			tg := int(model.Target.Pix[i+1])
			tb := int(model.Target.Pix[i+2])
			cr := int(model.Current.Pix[i])
			cg := int(model.Current.Pix[i+1])
			cb := int(model.Current.Pix[i+2])
			i += 4
			rsum += int64((tr-cr)*a + cr*0x101)
			gsum += int64((tg-cg)*a + cg*0x101)
			bsum += int64((tb-cb)*a + cb*0x101)
			count++
		}
	}
	if count == 0 {
		return Color{}
	}
	r := clampInt(int(rsum/count)>>8, 0, 255)
	g := clampInt(int(gsum/count)>>8, 0, 255)
	b := clampInt(int(bsum/count)>>8, 0, 255)
	return Color{r, g, b, alpha}
}

func (model *Model) computeScore(lines []Scanline, c Color, buffer *image.RGBA) float64 {
	Copy(buffer, model.Current, lines)
	Draw(buffer, c, lines)
	return differencePartial(model.Target, model.Current, buffer, model.Score, lines)
}

func (model *Model) Energy(alpha int, shape Shape, buffer *image.RGBA) float64 {
	lines := shape.Rasterize()
	c := model.computeColor(lines, alpha)
	s := model.computeScore(lines, c, buffer)
	return s
}
