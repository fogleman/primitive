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
	num_workers := 4

	background_color := MakeColor(color.NRGBA{uint8(32), uint8(17), uint8(202), uint8(240)})
	new_size := testingImage.Bounds().Dx()*testingImage.Bounds().Dy() - 22

	// run the function under test with generated values
	testModel := NewModel(testingImage, background_color, new_size, num_workers)

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

	fail_image_eq := false
	for key, val := range testingImageRGBA.Pix {
		if testModel.Target.Pix[key] != val {
			fail_image_eq = true
			break
		}
	}
	if fail_image_eq {
		t.Error("Target Image byte mismatch in NewModel")
	}

	// check the initialized current image

	if len(testModel.Target.Pix) != len(testingImageRGBA.Pix) {
		t.Error("Current Image size mismatch in NewModel")
	}

	// check that a score is at leaset set
	if testModel.Score == 0 {
		t.Error("Score not set in NewModel")
	}

	// check that a context exists
	if testModel.Context == nil {
		t.Error("Context is not set in NewModel")
	}

	// check that we have the correct number of workers
	if len(testModel.Workers) != num_workers {
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

	nrgba_color := testingModel.Background.NRGBA()
	nrgba_image := imageToRGBA(testingImage)

	color_check_fail := false
	for p := 0; p < len(nrgba_image.Pix); p += 4 {
		if nrgba_image.Pix[p] != nrgba_color.R ||
			nrgba_image.Pix[p+1] != nrgba_color.G ||
			nrgba_image.Pix[p+2] != nrgba_color.B {
			fmt.Println(fmt.Sprintf("r: %d %d g: %d %d b %d %d",
				nrgba_image.Pix[p], nrgba_color.R,
				nrgba_image.Pix[p+1], nrgba_color.B,
				nrgba_image.Pix[p+2], nrgba_color.G))
			color_check_fail = true
			break
		}
	}
	if color_check_fail {
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
		before_score := testingModel.Score
		testingModel.Step(ShapeType(i), alpha, 10, notify)
		after_score := testingModel.Score

		if before_score <= after_score {
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
		before_score := testingModel.Score
		before_shapes_len := len(testingModel.Shapes)
		before_colors_len := len(testingModel.Colors)
		before_scores_len := len(testingModel.Scores)
		resultState := testingModel.runWorkers(ShapeType(i), alpha, 1000, 100, 16)

		testingModel.Add(resultState.Shape, resultState.Alpha, notify)
		after_score := testingModel.Score

		if len(testingModel.Shapes) != before_shapes_len+1 {
			t.Error("Shape not added to model in Model.Add")
		}

		if len(testingModel.Colors) != before_colors_len+1 {
			t.Error("Color not added to model in Model.Add")
		}

		if len(testingModel.Scores) != before_scores_len+1 {
			t.Error("Score not added to modle in Model.Add")
		}

		if after_score >= before_score {
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

		num_shapes := len(testingModel.Shapes)
		min_delta := .001
		num_qualifying_scores := 1
		previous_score := float64(10)

		for j, score := range testingModel.Scores {
			score = testingModel.Scores[j]
			delta := previous_score - score

			if delta >= min_delta {
				num_qualifying_scores += 1
				previous_score = score
			}
		}

		testFrames := testingModel.Frames(0.001, notify)

		if notify.messages["Evaulating shape in Frames"] != num_shapes {
			t.Error(fmt.Sprintf("Wrong number of shape evaluations in Frames: %d",
				notify.messages["Evaulating shape in Frames"]))
		}

		if notify.messages["Called Fill"] != num_shapes {
			t.Error("Wrong number of fill executions in Frames")
		}

		if len(testFrames) != num_qualifying_scores {
			t.Error(
				fmt.Sprintf("Mismatch in qualifing frames and returned frames in Frames. Got %d, expected %d.",
					num_qualifying_scores, len(testFrames)))
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
