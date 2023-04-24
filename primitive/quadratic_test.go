package primitive

import (
	"fmt"
	"testing"

	"github.com/fogleman/gg"
)

func TestDrawQuadratic(t *testing.T) {

	worker := NewWorker(imageToRGBA(createTestImage()))
	context := gg.NewContext(100, 100)
	context.SetRGBA255(201, 97, 144, 185)

	my_quadratic := Quadratic{worker, 10.21, 20.54, 63.31,
		23.64, 37.76, 10.11, 8}

	notify := NewTestStringNotifier()

	my_quadratic.Draw(context, .9, notify)
	context.Fill()

	contextState := Hash(context.Image())

	// This value was pre-computed from static inputs
	if contextState != "e36e71bff2e58d03346cc6d26673154a" {
		t.Error(fmt.Sprintf("Incorect state after Draw in Quadratic: %s", contextState))
	}
}

func TestSVGQuadratic(t *testing.T) {
	worker := NewWorker(imageToRGBA(createTestImage2()))
	my_quadratic := Quadratic{worker, 10.21, 20.54, 63.31,
		23.64, 37.76, 10.11, 8}

	SVG := my_quadratic.SVG("myAttrs")

	// This value was pre-computed from static inputs
	if SVG !=
		"<path myAttrs fill=\"none\" d=\"M 10.210000 20.540000 Q 63.310000 "+
			"23.640000, 37.760000 10.110000\" stroke-width=\"8.000000\" />" {
		t.Error(fmt.Sprintf("Incorect SVG after SVG in Quadratic: %s\n Expected: "+
			"<path myAttrs fill=\"none\" d=\"M 10.210000 20.540000 Q 63.310000 "+
			"23.640000, 37.760000 10.110000\" stroke-width=\"8.000000\" />", SVG))
	}
}

func TestRasterizeQuadratic(t *testing.T) {
	worker := NewWorker(imageToRGBA(createTestImage()))
	context := gg.NewContext(100, 100)
	context.SetRGBA255(224, 117, 232, 187)

	my_quadratic := Quadratic{worker, 11.21, 20.54, 63.31,
		23.64, 37.76, 10.12, 8}

	lines := my_quadratic.Rasterize()

	linesState := Hash(lines)

	// This value was pre-computed from static inputs
	if linesState != "e74d96173314525b56d24d3a0aadc496" {
		t.Error(fmt.Sprintf("Incorect state after Rasterize in Quadraditc: %s", linesState))
	}
}
