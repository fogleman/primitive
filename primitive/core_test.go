package primitive

import (
	"fmt"
	"image/color"
	"testing"
)

func TestComputeColor(t *testing.T) {
	rgbaImage1 := imageToRGBA(createTestImage())
	rgbaImage2 := imageToRGBA(createTestImage2())

	model := createTestModel()
	state := model.runWorkers(3, 200, 1000, 100, 16)
	shape := state.Shape
	lines := shape.Rasterize()

	color := computeColor(rgbaImage1, rgbaImage2, lines, 200)

	// These values are pre-computed for the static inputs
	if color.R != 152 || color.G != 96 || color.B != 70 || color.A != 200 {
		t.Error(fmt.Sprintf(
			"Wrong color generated in ComputeColor.) R=%d, G=%d, B=%d, A=%d",
			color.R, color.G, color.B, color.A,
		))
	}
}

func TestCopyLines(t *testing.T) {

	rgbaImage1 := imageToRGBA(createTestImage())
	rgbaImage2 := imageToRGBA(createTestImage2())

	model := createTestModel()
	state := model.runWorkers(3, 193, 1000, 100, 16)
	shape := state.Shape
	lines := shape.Rasterize()

	copyLines(rgbaImage1, rgbaImage2, lines)

	imageHash := Hash(rgbaImage1)

	// This value was pre-computed from static inputs
	if imageHash != "fe247f8d1d4eaf12f57605a73129b3c3" {
		t.Error(fmt.Sprintf("Bad result image after copy in CopyLines: %s", imageHash))
	}
}

func TestDrawLines(t *testing.T) {

	rgbaImage2 := imageToRGBA(createTestImage2())
	lines := getStaticScanLines()
	color := MakeColor(color.NRGBA{uint8(217), uint8(44), uint8(143), uint8(230)})
	notify := NewTestStringNotifier()

	drawLines(rgbaImage2, color, lines, notify)

	pixTotal := 0
	for _, pix_val := range rgbaImage2.Pix {
		pixTotal += int(pix_val)
	}

	fmt.Printf("%#+v", lines)

	// This value was pre-computed from static inputs
	if pixTotal != 9353011 {
		t.Error(fmt.Sprintf("Bad result image after draw in DrawLines: %d", pixTotal))
	}
}

func TestDifferenceFull(t *testing.T) {

	rgbaImage1 := imageToRGBA(createTestImage())
	rgbaImage2 := imageToRGBA(createTestImage2())

	result := int(differenceFull(rgbaImage1, rgbaImage2) * 1000000)

	// This value was pre-computed from static inputs
	if result != 357714 {
		t.Error(fmt.Sprintf("Bad result for DifferenceFull: %d", result))
	}
}

func TestDifferencePartial(t *testing.T) {

	rgbaImage1 := imageToRGBA(createTestImage())
	rgbaImage2 := imageToRGBA(createTestImage2())

	notify := NewTestStringNotifier()
	model := createTestModel()
	state := model.runWorkers(3, 193, 1000, 100, 16)
	shape := state.Shape
	model.Add(shape, 200, notify)
	lines := shape.Rasterize()

	result := differencePartial(rgbaImage2, rgbaImage1, model.Current, 10, lines)
	intResult := int(result * 1000000)

	// This value was pre-computed from static inputs
	if intResult != 9996133 {
		t.Error(fmt.Sprintf("Bad result for DifferenceFull: %d", intResult))
	}

}
