package primitive

import (
	"fmt"
	"image"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

const SaveFrames = false

type ShapeType int

const (
	ShapeTypeAny ShapeType = iota
	ShapeTypeTriangle
	ShapeTypeRectangle
	ShapeTypeEllipse
	ShapeTypeCircle
	ShapeTypeRotatedRectangle
)

type Model struct {
	W, H    int
	Target  *image.RGBA
	Current *image.RGBA
	Buffer  *image.RGBA
	Context *gg.Context
	Score   float64
	Alpha   int
	Scale   int
	Mode    ShapeType
	Shapes  []Shape
}

func NewModel(target image.Image, alpha, scale int, mode ShapeType) *Model {
	c := averageImageColor(target)
	// c = color.White
	// c = color.Black
	size := target.Bounds().Size()
	model := &Model{}
	model.W = size.X
	model.H = size.Y
	model.Target = imageToRGBA(target)
	model.Current = uniformRGBA(target.Bounds(), c)
	model.Buffer = uniformRGBA(target.Bounds(), c)
	model.Score = differenceFull(model.Target, model.Current)
	model.Context = gg.NewContext(model.W*scale, model.H*scale)
	model.Context.Scale(float64(scale), float64(scale))
	model.Context.SetColor(c)
	model.Context.Clear()
	model.Alpha = alpha
	model.Scale = scale
	model.Mode = mode
	return model
}

func (model *Model) Run(n int) image.Image {
	start := time.Now()
	for i := 1; i <= n; i++ {
		model.Step()
		elapsed := time.Since(start).Seconds()
		fmt.Printf("%d, %.3f, %.6f\n", i, elapsed, model.Score)
		if SaveFrames {
			SavePNG("out.png", model.Current)
			model.Context.SavePNG(fmt.Sprintf("out%03d.png", i))
		}
	}
	return model.Context.Image()
}

func (model *Model) Step() {
	state := model.BestRandomState(model.Buffer, model.Mode, 3000)
	// state = Anneal(state, 0.1, 0.00001, 25000).(*State)
	state = HillClimb(state, 1000).(*State)
	model.Add(state.Shape)
}

func (model *Model) worker(i int, ch chan *State) {
	buffer := image.NewRGBA(model.Target.Bounds())
	state := model.BestRandomState(buffer, ShapeType(i+1), 3000)
	// state = Anneal(state, 0.1, 0.00001, 25000).(*State)
	state = HillClimb(state, 1000).(*State)
	ch <- state
}

func (model *Model) GoStep() {
	n := 3
	ch := make(chan *State, n)
	for i := 0; i < n; i++ {
		go model.worker(i, ch)
	}
	var bestShape Shape
	var bestEnergy float64
	for i := 0; i < n; i++ {
		state := <-ch
		shape := state.Shape
		energy := model.Energy(shape, model.Buffer)
		fmt.Printf("%.6f\n", energy)
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestShape = shape
		}
	}
	model.Add(bestShape)
}

func (model *Model) BestRandomState(buffer *image.RGBA, t ShapeType, n int) *State {
	var bestEnergy float64
	var bestState *State
	for i := 0; i < n; i++ {
		state := model.RandomState(buffer, t)
		energy := state.Energy()
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

func (model *Model) RandomState(buffer *image.RGBA, t ShapeType) *State {
	switch t {
	default:
		return model.RandomState(buffer, ShapeType(rand.Intn(4)+1))
	case ShapeTypeTriangle:
		return NewState(model, buffer, NewRandomTriangle(model.W, model.H))
	case ShapeTypeRectangle:
		return NewState(model, buffer, NewRandomRectangle(model.W, model.H))
	case ShapeTypeEllipse:
		return NewState(model, buffer, NewRandomEllipse(model.W, model.H))
	case ShapeTypeCircle:
		return NewState(model, buffer, NewRandomCircle(model.W, model.H))
	case ShapeTypeRotatedRectangle:
		return NewState(model, buffer, NewRandomRotatedRectangle(model.W, model.H))
	}
}

func (model *Model) Add(shape Shape) {
	lines := shape.Rasterize()
	c := model.computeColor(lines, model.Alpha)
	s := model.computeScore(lines, c, model.Buffer)
	Draw(model.Current, c, lines)
	model.Score = s
	model.Context.SetRGBA255(c.R, c.G, c.B, c.A)
	shape.Draw(model.Context)
	model.Shapes = append(model.Shapes, shape)
}

func (model *Model) computeColor(lines []Scanline, alpha int) Color {
	var count int
	var rsum, gsum, bsum float64
	a := float64(alpha) / 255
	for _, line := range lines {
		i := model.Target.PixOffset(line.X1, line.Y)
		for x := line.X1; x <= line.X2; x++ {
			count++
			tr := float64(model.Target.Pix[i])
			tg := float64(model.Target.Pix[i+1])
			tb := float64(model.Target.Pix[i+2])
			cr := float64(model.Current.Pix[i])
			cg := float64(model.Current.Pix[i+1])
			cb := float64(model.Current.Pix[i+2])
			i += 4
			rsum += (a*cr - cr + tr) / a
			gsum += (a*cg - cg + tg) / a
			bsum += (a*cb - cb + tb) / a
		}
	}
	if count == 0 {
		return Color{}
	}
	r := clampInt(int(rsum/float64(count)), 0, 255)
	g := clampInt(int(gsum/float64(count)), 0, 255)
	b := clampInt(int(bsum/float64(count)), 0, 255)
	return Color{r, g, b, alpha}
}

func (model *Model) computeScore(lines []Scanline, c Color, buffer *image.RGBA) float64 {
	copy(buffer.Pix, model.Current.Pix)
	Draw(buffer, c, lines)
	return differencePartial(model.Target, model.Current, buffer, model.Score, lines)
}

func (model *Model) Energy(shape Shape, buffer *image.RGBA) float64 {
	lines := shape.Rasterize()
	c := model.computeColor(lines, model.Alpha)
	s := model.computeScore(lines, c, buffer)
	return s
}
