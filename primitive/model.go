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
	Background Color
	Target     *image.RGBA
	Current    *image.RGBA
	Buffer     *image.RGBA
	Context    *gg.Context
	Score      float64
	Alpha      int
	Size       int
	Mode       Mode
	Shapes     []Shape
	Colors     []Color
	Scores     []float64
	SVGs       []string
}

func NewModel(target image.Image, background Color, alpha, size int, mode Mode) *Model {
	model := &Model{}
	model.W = target.Bounds().Size().X
	model.H = target.Bounds().Size().Y
	model.Background = background
	model.Alpha = alpha
	model.Size = size
	model.Mode = mode
	model.Target = imageToRGBA(target)
	model.Current = uniformRGBA(target.Bounds(), background.NRGBA())
	model.Buffer = uniformRGBA(target.Bounds(), background.NRGBA())
	model.Score = differenceFull(model.Target, model.Current)
	model.Context = model.newContext()
	return model
}

func (model *Model) sizeAndScale() (w, h int, scale float64) {
	aspect := float64(model.W) / float64(model.H)
	if aspect >= 1 {
		w = model.Size
		h = int(float64(model.Size) / aspect)
		scale = float64(model.Size) / float64(model.W)
	} else {
		w = int(float64(model.Size) * aspect)
		h = model.Size
		scale = float64(model.Size) / float64(model.H)
	}
	return
}

func (model *Model) newContext() *gg.Context {
	w, h, scale := model.sizeAndScale()
	dc := gg.NewContext(w, h)
	dc.Scale(scale, scale)
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
		shape.Draw(dc)
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
	w, h, scale := model.sizeAndScale()
	c := model.Background
	var lines []string
	lines = append(lines, fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" version=\"1.1\" width=\"%d\" height=\"%d\">", w, h))
	lines = append(lines, fmt.Sprintf("<rect x=\"0\" y=\"0\" width=\"%d\" height=\"%d\" fill=\"#%02x%02x%02x\" />", w, h, c.R, c.G, c.B))
	lines = append(lines, fmt.Sprintf("<g transform=\"scale(%f) translate(0.5 0.5)\">", scale))
	lines = append(lines, model.SVGs...)
	lines = append(lines, "</g>")
	lines = append(lines, "</svg>")
	return strings.Join(lines, "\n")
}

func (model *Model) Add(shape Shape) {
	lines := shape.Rasterize()
	c := model.computeColor(lines, model.Alpha)
	s := model.computeScore(lines, c, model.Buffer)
	Draw(model.Current, c, lines)

	attrs := "fill=\"#%02x%02x%02x\" fill-opacity=\"%f\""
	attrs = fmt.Sprintf(attrs, c.R, c.G, c.B, float64(c.A)/255)
	svg := shape.SVG(attrs)

	model.Score = s
	model.Shapes = append(model.Shapes, shape)
	model.Colors = append(model.Colors, c)
	model.Scores = append(model.Scores, s)
	model.SVGs = append(model.SVGs, svg)

	model.Context.SetRGBA255(c.R, c.G, c.B, c.A)
	shape.Draw(model.Context)
	model.Context.Fill()
}

func (model *Model) Step() {
	state := model.runWorkers(model.Mode, 100, 100, 8)
	state = HillClimb(state, 1000).(*State)
	model.Add(state.Shape)
}

func (model *Model) runWorkers(t Mode, n, age, m int) *State {
	wn := runtime.GOMAXPROCS(0)
	ch := make(chan *State, wn)
	wm := m / wn
	if m%wn != 0 {
		wm++
	}
	for i := 0; i < wn; i++ {
		go model.runWorker(t, n, age, wm, ch)
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

func (model *Model) runWorker(t Mode, n, age, m int, ch chan *State) {
	buffer := image.NewRGBA(model.Target.Bounds())
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	state := model.BestHillClimbState(buffer, t, n, age, m, rnd)
	ch <- state
}

func (model *Model) BestHillClimbState(buffer *image.RGBA, t Mode, n, age, m int, rnd *rand.Rand) *State {
	var bestEnergy float64
	var bestState *State
	for i := 0; i < m; i++ {
		state := model.BestRandomState(buffer, t, n, rnd)
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

func (model *Model) BestRandomState(buffer *image.RGBA, t Mode, n int, rnd *rand.Rand) *State {
	var bestEnergy float64
	var bestState *State
	for i := 0; i < n; i++ {
		state := model.RandomState(buffer, t, rnd)
		energy := state.Energy()
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

func (model *Model) RandomState(buffer *image.RGBA, t Mode, rnd *rand.Rand) *State {
	switch t {
	default:
		return model.RandomState(buffer, Mode(rnd.Intn(5)+1), rnd)
	case ModeTriangle:
		return NewState(model, buffer, NewRandomTriangle(model.W, model.H, rnd))
	case ModeRectangle:
		return NewState(model, buffer, NewRandomRectangle(model.W, model.H, rnd))
	case ModeEllipse:
		return NewState(model, buffer, NewRandomEllipse(model.W, model.H, rnd))
	case ModeCircle:
		return NewState(model, buffer, NewRandomCircle(model.W, model.H, rnd))
	case ModeRotatedRectangle:
		return NewState(model, buffer, NewRandomRotatedRectangle(model.W, model.H, rnd))
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

func (model *Model) Energy(shape Shape, buffer *image.RGBA) float64 {
	lines := shape.Rasterize()
	c := model.computeColor(lines, model.Alpha)
	s := model.computeScore(lines, c, buffer)
	return s
}
