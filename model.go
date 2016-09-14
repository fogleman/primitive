package primitive

import (
	"fmt"
	"image"
	"time"

	"github.com/fogleman/gg"
)

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
	size := target.Bounds().Size()
	model := &Model{}
	model.W = size.X
	model.H = size.Y
	model.Target = imageToRGBA(target)
	model.Current = uniformRGBA(target.Bounds(), c)
	model.Buffer = uniformRGBA(target.Bounds(), c)
	model.Score = differenceFull(model.Target, model.Current)
	model.Context = gg.NewContext(model.W*4, model.H*4)
	model.Context.Scale(4, 4)
	model.Context.SetColor(c)
	model.Context.Clear()
	return model
}

func (model *Model) Run() {
	frame := 0
	start := time.Now()
	for {
		model.Step()
		elapsed := time.Since(start).Seconds()
		fmt.Printf("%d, %.3f, %.6f\n", frame, elapsed, model.Score)
		if frame%1 == 0 {
			path := fmt.Sprintf("out%03d.png", frame)
			// SavePNG(path, model.Current)
			model.Context.SavePNG(path)
		}
		frame++
	}
}

func (model *Model) Step() {
	// state := NewState(model, NewRandomTriangle(model.W, model.H))
	// state := NewState(model, NewRandomRectangle(model.W, model.H))
	state := NewState(model, NewRandomCircle(model.W, model.H))
	// fmt.Println(PreAnneal(state, 10000))
	state = Anneal(state, 0.2, 0.0001, 10000).(*State)
	model.Add(state.Shape)
}

func (model *Model) Add(shape Shape) {
	lines := shape.Rasterize()
	c := model.computeColor(lines, 128)
	s := model.computeScore(lines, c)
	Draw(model.Current, c, lines)
	model.Score = s
	shape.Draw(model.Context)
	model.Context.SetRGBA255(c.R, c.G, c.B, c.A)
	model.Context.Fill()
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
