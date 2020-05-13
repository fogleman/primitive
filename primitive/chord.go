package primitive

import (
	"fmt"
	"strings"
	"math"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/raster"
)

type Chord struct {
	Worker *Worker
	A1, A2 float64
	//X1, Y1 float64
	//X2, Y2 float64
}

func NewRandomChord(worker *Worker) *Chord {
	rnd := worker.Rnd

	// These will be on (-1, 1)
	a1 := rnd.Float64() * float64(360)
	a2 := rnd.Float64() * float64(360)

	q := &Chord{worker, a1, a2}
	q.Mutate()
	return q
}

func (q *Chord) GetPositions() (float64, float64, float64, float64) {
	// So we need to make them on a unit circle centered around worker.W / 2, worker.H / 2
	// Of width worker.W / 2 and height worker.H / 2
	half_width := float64(q.Worker.W) / 2
	half_height:= float64(q.Worker.H) / 2

	x1 := math.Cos(q.A1)
	y1 := math.Sin(q.A1)
	x2 := math.Cos(q.A2)
	y2 := math.Sin(q.A2)

	x1 = (x1 * half_width) + half_width
	y1 = (y1 * half_height) + half_height
	x2 = (x2 * half_width) + half_width
	y2 = (y2 * half_height) + half_height
	return x1, y1, x2, y2
}

func (q *Chord) Draw(dc *gg.Context, scale float64) {
	x1, y1, x2, y2 := q.GetPositions()

	dc.MoveTo(x1, y1)
	dc.MoveTo(x2, y2)
	dc.SetLineWidth(1)
	dc.Stroke()
}

func (q *Chord) SVG(attrs string) string {
	// TODO: this is a little silly
	x1, y1, x2, y2 := q.GetPositions()
	attrs = strings.Replace(attrs, "fill", "stroke", -1)
	return fmt.Sprintf(
		"<line %s x1=\"%f\" y1=\"%f\" x2=\"%f\" y2=\"%f\" stroke-width=\"1\" />",
		attrs, x1, y1, x2, y2)
}

func (q *Chord) Copy() Shape {
	a := *q
	return &a
}

func (q *Chord) Mutate() {
	const m = 16
	rnd := q.Worker.Rnd
	for {
		switch rnd.Intn(2) {
		case 0:
			q.A1 = clamp(q.A1 + rnd.NormFloat64() * m, 0, 360)
		case 1:
			q.A2 = clamp(q.A2 + rnd.NormFloat64() * m, 0, 360)
		}
		if q.Valid() {
			break
		}
	}
}

func (q *Chord) Valid() bool {
	return true;
}

func (q *Chord) Rasterize() []Scanline {
	x1, y1, x2, y2 := q.GetPositions()

	var path raster.Path
	p1 := fixp(x1, y1)
	p2 := fixp(x2, y2)
	path.Start(p1)
	path.Add1(p2)
	width := fix(1)
	return strokePath(q.Worker, path, width, raster.RoundCapper, raster.RoundJoiner)
}
