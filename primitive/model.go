package primitive

import (
	"fmt"
	"image"
	"time"

	"github.com/fogleman/gg"
)

const Scale = 4

type Model struct {
	W, H    int
	Target  *image.RGBA
	Current *image.RGBA
	Buffer  *image.RGBA
	Score   float64
	Context *gg.Context
}

func NewModel(target image.Image) *Model {
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
	model.Context = gg.NewContext(model.W*Scale, model.H*Scale)
	model.Context.Scale(Scale, Scale)
	model.Context.SetColor(c)
	model.Context.Clear()
	return model
}

func (model *Model) Run() {
	frame := 1
	start := time.Now()
	for {
		model.Step()
		elapsed := time.Since(start).Seconds()
		fmt.Printf("%d, %.3f, %.6f\n", frame, elapsed, model.Score)
		if frame%1 == 0 {
			path := fmt.Sprintf("out%03d.png", frame)
			SavePNG("out.png", model.Current)
			model.Context.SavePNG(path)
		}
		frame++
	}
}

func (model *Model) Step() {
	state := model.CreateState()
	// state := model.RandomState()
	// fmt.Println(PreAnneal(state, 10000))
	state = Anneal(state, 0.1, 0.00001, 30000).(*State)
	// state = HillClimb(state, 1000).(*State)
	model.Add(state.Shape)
}

func (model *Model) CreateState() *State {
	var bestEnergy float64
	var bestState *State
	for i := 0; i < 100; i++ {
		state := model.RandomState()
		energy := state.Energy()
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

func (model *Model) RandomState() *State {
	return NewState(model, NewRandomTriangle(model.W, model.H))
	// return NewState(model, NewRandomRectangle(model.W, model.H))
	// return NewState(model, NewRandomCircle(model.W, model.H))
	// return NewState(model, NewRandomEllipse(model.W, model.H))
	// switch rand.Intn(2) {
	// case 0:
	// 	return NewState(model, NewRandomRectangle(model.W, model.H))
	// case 1:
	// 	return NewState(model, NewRandomEllipse(model.W, model.H))
	// }
	// return nil
}

func (model *Model) Add(shape Shape) {
	lines := shape.Rasterize()
	c := model.computeColor(lines, 128)
	s := model.computeScore(lines, c)
	Draw(model.Current, c, lines)
	model.Score = s
	model.Context.SetRGBA255(c.R, c.G, c.B, c.A)
	shape.Draw(model.Context)
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

func (model *Model) computeScore(lines []Scanline, c Color) float64 {
	copy(model.Buffer.Pix, model.Current.Pix)
	Draw(model.Buffer, c, lines)
	return differencePartial(model.Target, model.Current, model.Buffer, model.Score, lines)
}
