package primitive

import (
	"fmt"
	"image"
	"strings"

	"github.com/fogleman/gg"
)

// Model contains the state for a transform job
type Model struct {
	ScaledWidth  int
	ScaledHeight int
	Scale        float64
	Background   Color
	Target       *image.RGBA
	Current      *image.RGBA
	Context      *gg.Context
	Score        float64
	Shapes       []Shape
	Colors       []Color
	Scores       []float64
	Workers      []*Worker
}

// NewModel creates a model which handles state for the operation of
// reading from an imput image and translating it into primitives.
// During this process, each iteration is a 'step' initiated by the Step method
func NewModel(target image.Image, background Color, size, numWorkers int) *Model {
	width := target.Bounds().Size().X
	height := target.Bounds().Size().Y
	aspect := float64(width) / float64(height)
	var scaledWidth, scaledHeight int
	var scale float64

	//If the image is wider than it is tall, the width should be set to 'size'
	if aspect >= 1 {
		scaledWidth = size
		scaledHeight = int(float64(size) / aspect)
		scale = float64(size) / float64(width)
	} else {

		// If the image is taller than it is wide, the height should instead be set to 'size'
		scaledWidth = int(float64(size) * aspect)
		scaledHeight = size
		scale = float64(size) / float64(height)
	}

	// Set up the model instance with:
	// * scale values
	// * background color
	// * the image we're operating on in RGB format
	// * the image we're going to generate set to an image of the same size as our input image
	//   initialized to our background color
	// * a reference to the input image
	// * a reference to the current state of the new image
	// * a score value, which will represent how closely a
	//   given shape describes the original image
	// * a context object, which will hold the state of the in-process image modification
	//   as well as provide helper methods for modifying the image
	// * an initialized set of workers in the number specified at the command line
	model := &Model{}
	model.ScaledWidth = scaledWidth
	model.ScaledHeight = scaledHeight
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

// Creates a new context from values in model
func (model *Model) newContext() *gg.Context {
	dc := gg.NewContext(model.ScaledWidth, model.ScaledHeight)
	dc.Scale(model.Scale, model.Scale)
	dc.Translate(0.5, 0.5)
	dc.SetColor(model.Background.NRGBA())
	dc.Clear()
	return dc
}

// Frames creates a sequence of frames for outputing to a GIF
func (model *Model) Frames(scoreDelta float64, notify Notifier) []image.Image {
	var result []image.Image
	dc := model.newContext()

	// Populate the initial value of the result with the current blank context image
	result = append(result, imageToRGBA(dc.Image()))

	// Initialize previous to a high value that should be easy to improve on.
	previous := 10.0

	// iterate over all of the shapes in the model
	for i, shape := range model.Shapes {
		notify.Notify("Evaulating shape in Frames")
		// find the color origionally associated with this shape
		c := model.Colors[i]

		// set the current context color to the color for this shape
		dc.SetRGBA255(c.R, c.G, c.B, c.A)

		// draw the shape and ifll it in with the current color
		shape.Draw(dc, model.Scale, notify)
		dc.Fill()
		notify.Notify("Called Fill")
		// Find the score associated with this shape
		score := model.Scores[i]
		delta := previous - score

		// If this shape improved the likeness by at least the passed in delta
		// Append this image to the result.
		if delta >= scoreDelta {
			previous = score
			result = append(result, imageToRGBA(dc.Image()))
		}
	}
	return result
}

// SVG outputs an svg string from the current state of the model
func (model *Model) SVG() string {
	bg := model.Background
	var lines []string
	svgTag := "<svg xmlns=\"http://www.w3.org/2000/svg\"" +
		" version=\"1.1\" width=\"%d\" height=\"%d\">"
	lines = append(lines,
		fmt.Sprintf(svgTag, model.ScaledWidth, model.ScaledHeight))
	rectTag := "<rect x=\"0\" y=\"0\" width=\"%d\" height=\"%d\" fill=\"#%02x%02x%02x\" />"
	lines = append(lines, fmt.Sprintf(rectTag, model.ScaledWidth, model.ScaledHeight, bg.R, bg.G, bg.B))
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

// Add adds a shape to the model
func (model *Model) Add(shape Shape, alpha int, notify Notifier) {
	notify.Notify("Model.Add was called")
	before := copyRGBA(model.Current)
	lines := shape.Rasterize()
	color := computeColor(model.Target, model.Current, lines, alpha)
	drawLines(model.Current, color, lines, notify)
	score := differencePartial(model.Target, before, model.Current, model.Score, lines)

	model.Score = score
	model.Shapes = append(model.Shapes, shape)
	model.Colors = append(model.Colors, color)
	model.Scores = append(model.Scores, score)

	model.Context.SetRGBA255(color.R, color.G, color.B, color.A)
	shape.Draw(model.Context, model.Scale, notify)
}

// Step kicks of a new shape creation for adding to the model
func (model *Model) Step(shapeType ShapeType, alpha, repeat int, notify Notifier) int {
	resultState := model.runWorkers(shapeType, alpha, 1000, 100, 16)
	// state = HillClimb(state, 1000).(*State)
	model.Add(resultState.Shape, resultState.Alpha, notify)

	//Optional additional optimizations and shape additions
	for i := 0; i < repeat; i++ {
		resultState.Worker.Init(model.Current, model.Score)
		beforeEnergy := resultState.Energy()
		resultState = HillClimb(resultState, 100).(*State)
		afterEnergy := resultState.Energy()

		//If we are no longer effectively optimizing, quit.
		if beforeEnergy == afterEnergy {
			notify.Notify("breaking out due to no optimization")
			break
		}
		model.Add(resultState.Shape, resultState.Alpha, notify)
	}

	counter := 0
	for _, worker := range model.Workers {
		counter += worker.Counter
	}
	return counter
}

func (model *Model) runWorkers(
	t ShapeType, alpha, triesPerWorker, age, totalClimbes int) *State {
	numberOfWorkers := len(model.Workers)
	workerChannel := make(chan *State, numberOfWorkers)

	climbesPerWorker := totalClimbes / numberOfWorkers

	//Err on the side of more climbes rather than less climbes
	if climbesPerWorker%numberOfWorkers != 0 {
		climbesPerWorker++
	}
	for i := 0; i < numberOfWorkers; i++ {
		worker := model.Workers[i]
		worker.Init(model.Current, model.Score)
		go model.runWorker(
			worker, t, alpha, triesPerWorker, age, climbesPerWorker, workerChannel)
	}
	var bestEnergy float64
	var bestState *State
	for i := 0; i < numberOfWorkers; i++ {
		state := <-workerChannel
		energy := state.Energy()
		if i == 0 || energy < bestEnergy {
			bestEnergy = energy
			bestState = state
		}
	}
	return bestState
}

func (model *Model) runWorker(
	worker *Worker, t ShapeType, alpha, triesPerWorker, age, climbs int, ch chan *State) {
	ch <- worker.BestHillClimbState(t, alpha, triesPerWorker, age, climbs)
}
