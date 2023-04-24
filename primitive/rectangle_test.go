package primitive

import (
	"fmt"
	"testing"

	"github.com/fogleman/gg"
)

func TestDrawRectangle(t *testing.T) {

	worker := NewWorker(imageToRGBA(createTestImage()))
	context := gg.NewContext(100, 100)
	context.SetRGBA255(224, 117, 232, 187)

	my_rectangle := Rectangle{worker, 11, 22, 49, 59}
	notify := NewTestStringNotifier()

	my_rectangle.Draw(context, 10, notify)
	context.Fill()

	contextState := Hash(context.Image())

	// This value was pre-computed from static inputs
	if contextState != "4a75200d8cd9008322c87c3252bae353" {
		t.Error(fmt.Sprintf("Incorect state after Draw in Rectangle: %s", contextState))
	}
}

func TestSVGRectangle(t *testing.T) {
	worker := NewWorker(imageToRGBA(createTestImage2()))
	my_rectangle := Rectangle{worker, 123, 134, 21, 55}
	SVG := my_rectangle.SVG("myAttrs")

	// This value was pre-computed from static inputs
	if SVG != "<rect myAttrs x=\"21\" y=\"55\" width=\"103\" height=\"80\" />" {
		t.Error(fmt.Sprintf("Incorect SVG after SVG in Rectangle: %s", SVG))
	}
}

func TestRasterizeRectangle(t *testing.T) {
	worker := NewWorker(imageToRGBA(createTestImage()))
	context := gg.NewContext(100, 100)
	context.SetRGBA255(224, 117, 232, 187)

	my_rectangle := Rectangle{worker, 47, 60, 12, 23}

	lines := my_rectangle.Rasterize()

	linesState := Hash(lines)

	// This value was pre-computed from static inputs
	if linesState != "40335c30caee2d75ba58425ce47375de" {
		t.Error(fmt.Sprintf("Incorect state after Rasterize in Rectangle: %s", linesState))
	}
}
