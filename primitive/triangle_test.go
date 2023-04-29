package primitive

import (
	"fmt"
	"testing"

	"github.com/fogleman/gg"
)

func TestDrawTriangle(t *testing.T) {

	worker := NewWorker(imageToRGBA(createTestImage()))
	context := gg.NewContext(100, 100)
	context.SetRGBA255(201, 97, 134, 185)

	myTriangle := Triangle{worker, 10, 21, 63,
		23, 37, 10}

	notify := NewTestStringNotifier()

	myTriangle.Draw(context, .9, notify)
	context.Fill()

	contextState := Hash(context.Image())

	// This value was pre-computed from static inputs
	if contextState != "7a73188233ca5f4bd90d64863c3d69e6" {
		t.Error(fmt.Sprintf("Incorect state after Draw in Triangle: %s", contextState))
	}
}

func TestSVGTriangle(t *testing.T) {
	worker := NewWorker(imageToRGBA(createTestImage2()))
	myTriangle := Triangle{worker, 10, 20, 63,
		23, 37, 10}

	SVG := myTriangle.SVG("myAttrs")

	// This value was pre-computed from static inputs
	if SVG !=
		"<polygon myAttrs points=\"10,20 63,23 37,10\" />" {
		t.Error(fmt.Sprintf("Incorect SVG after SVG in Triangle: %s", SVG))
	}
}

func TestRasterizeTriangle(t *testing.T) {
	worker := NewWorker(imageToRGBA(createTestImage()))
	context := gg.NewContext(100, 100)
	context.SetRGBA255(224, 117, 232, 187)

	myTriangle := Triangle{worker, 11, 21, 63,
		23, 37, 10}

	lines := myTriangle.Rasterize()

	linesState := Hash(lines)

	// This value was pre-computed from static inputs
	if linesState != "e29bd06177bfc772af8077878696682d" {
		t.Error(fmt.Sprintf("Incorect state after Rasterize in Quadraditc: %s", linesState))
	}
}
