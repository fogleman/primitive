package primitive

import (
	"fmt"
	"image/color"
	"regexp"
	"testing"
)

func TestNewModel(t *testing.T) {

	testingImage := createTestImage()

	// generate a small number of workers
	numWorkers := 4

	backgroundColor := MakeColor(color.NRGBA{uint8(32), uint8(17), uint8(202), uint8(240)})
	newSize := testingImage.Bounds().Dx()*testingImage.Bounds().Dy() - 22

	// run the function under test with generated values
	testModel := NewModel(testingImage, backgroundColor, newSize, numWorkers)

	// check scale values
	// (These were pre-calculated based on this image and this size)
	if int(testModel.Scale) != 102 {
		t.Error(fmt.Sprintf("Scale mismatch in NewModel: %d", int(testModel.Scale)))
	}

	if testModel.ScaledHeight != 10593 {
		t.Error(fmt.Sprintf(
			"Scaled height mismatch in NewModel: %d", testModel.ScaledHeight))
	}

	if testModel.ScaledWidth != 15428 {
		t.Error(
			fmt.Sprintf("Scaled width mismatch in NewModel: %d",
				testModel.ScaledWidth))
	}

	// check the target image for equivelency with the testing image
	testingImageRGBA := imageToRGBA(testingImage)
	if len(testModel.Target.Pix) != len(testingImageRGBA.Pix) {
		t.Error("Target Image size mismatch in NewModel")
	}

	testingHash := Hash(testModel.Target.Pix)
	comparePixHash := Hash(testingImageRGBA.Pix)

	if testingHash != comparePixHash {
		t.Error("Target Image byte mismatch in NewModel")
	}

	// check the initialized current image

	if len(testModel.Target.Pix) != len(testingImageRGBA.Pix) {
		t.Error("Current Image size mismatch in NewModel")
	}

	// check that a score is at least set
	if testModel.Score == 0 {
		t.Error("Score not set in NewModel")
	}

	// check that a context exists
	if testModel.Context == nil {
		t.Error("Context is not set in NewModel")
	}

	// check that we have the correct number of workers
	if len(testModel.Workers) != numWorkers {
		t.Error("Wrong number of workers in NewModel")
	}

}

func TestNewContext(t *testing.T) {
	testingModel := createTestModel()
	testingContext := testingModel.newContext()

	// Make sure dimenions were set correctly in the model
	if testingContext.Width() != testingModel.ScaledWidth {
		t.Error("Width mismatch in newContext")
	}
	if testingContext.Height() != testingModel.ScaledHeight {
		t.Error("Height mismatch in newContext")
	}

	// Make sure dimensions were set correctly in the context
	testingImage := testingContext.Image()
	testingRectangle := testingImage.Bounds()
	if testingRectangle.Dx() != testingModel.ScaledWidth {
		t.Error("Width mismatch for image in newContext")
	}
	if testingRectangle.Dy() != testingModel.ScaledHeight {
		t.Error("Height mismatch for image in newContext")
	}

	// Make sure backgrouind color is uniform in the context image

	nrgbaColor := testingModel.Background.NRGBA()
	nrgbaImage := imageToRGBA(testingImage)

	colorCheckFail := false
	for p := 0; p < len(nrgbaImage.Pix); p += 4 {
		if nrgbaImage.Pix[p] != nrgbaColor.R ||
			nrgbaImage.Pix[p+1] != nrgbaColor.G ||
			nrgbaImage.Pix[p+2] != nrgbaColor.B {
			colorCheckFail = true
			break
		}
	}
	if colorCheckFail {
		t.Error("Current Image color mismatch in NewModel")
	}
}

func (model *Model) runWorkersOverride(st ShapeType, alpha int, triesPerWorker int, age int, totalClimbes int) *State {
	state := NewState(nil, nil, 0)
	return state
}

func TestStep(t *testing.T) {

	testingModel := createTestModel()

	for i := 0; i < 9; i++ {

		alpha := 223
		notify := NewTestStringNotifier()
		// After the step, the model score should be lower
		beforeScore := testingModel.Score
		testingModel.Step(ShapeType(i), alpha, 10, notify)
		afterScore := testingModel.Score

		if beforeScore <= afterScore {
			t.Error("No score improvement after Model.Step")
		}

		// Check if we broke out of optimization before 11
		if notify.messages["Model.Add was called"] != 11 {
			if notify.messages["breaking out due to no optimization"] == 0 {
				t.Error("Premeture exit to optimization in Model.Step")
			}
		}
	}
}

func TestAdd(t *testing.T) {
	testingModel := createTestModel()
	notify := NewTestStringNotifier()
	for i := 0; i <= 9; i++ {

		alpha := 227
		beforeScore := testingModel.Score
		beforeShapesLen := len(testingModel.Shapes)
		beforeColorsLen := len(testingModel.Colors)
		beforeScoresLen := len(testingModel.Scores)
		resultState := testingModel.runWorkers(ShapeType(i), alpha, 1000, 100, 16)

		testingModel.Add(resultState.Shape, resultState.Alpha, notify)
		afterScore := testingModel.Score

		if len(testingModel.Shapes) != beforeShapesLen+1 {
			t.Error("Shape not added to model in Model.Add")
		}

		if len(testingModel.Colors) != beforeColorsLen+1 {
			t.Error("Color not added to model in Model.Add")
		}

		if len(testingModel.Scores) != beforeScoresLen+1 {
			t.Error("Score not added to modle in Model.Add")
		}

		if afterScore >= beforeScore {
			t.Error("No score improvement in model.Add")
		}
	}

	if notify.messages["Model.Add was called"] != 10 {
		t.Error("Model.Add not reporting")
	}
	if notify.messages["drawLines was called"] != 10 {
		t.Error("drawLines not called in Model.Add")
	}
}

func TestFrames(t *testing.T) {
	testingModel := createTestModel()

	for i := 0; i < 3; i++ {
		notify := NewTestStringNotifier()
		alpha := 255

		testingModel.Step(ShapeType(i), alpha, 2, notify)

		numShapes := len(testingModel.Shapes)
		minDelta := .001
		numQualifyingScores := 1
		previousScore := float64(10)

		for j, score := range testingModel.Scores {
			score = testingModel.Scores[j]
			delta := previousScore - score

			if delta >= minDelta {
				numQualifyingScores++
				previousScore = score
			}
		}

		testFrames := testingModel.Frames(0.001, notify)

		if notify.messages["Evaulating shape in Frames"] != numShapes {
			t.Error(fmt.Sprintf("Wrong number of shape evaluations in Frames: %d",
				notify.messages["Evaulating shape in Frames"]))
		}

		if notify.messages["Called Fill"] != numShapes {
			t.Error("Wrong number of fill executions in Frames")
		}

		if len(testFrames) != numQualifyingScores {
			t.Error(
				fmt.Sprintf("Mismatch in qualifing frames and returned frames in Frames. Got %d, expected %d.",
					numQualifyingScores, len(testFrames)))
		}
	}
}

func TestSVG(t *testing.T) {
	testingModel := createTestModel()

	for i := 0; i < 9; i++ {
		notify := NewTestStringNotifier()
		alpha := 200
		testingModel.Step(ShapeType(i), alpha, 10, notify)
	}

	svg := testingModel.SVG()

	svgRegex := regexp.MustCompile(
		`(?i)^\s*(?:<\?xml[^>]*>\s*)?(?:<!doctype svg[^>]*>\s*)?<svg[^>]*>[^*]*<\/svg>\s*$`)

	if !svgRegex.MatchString(svg) {
		t.Error("Malformed SVG returned by SVG")
	}

}
