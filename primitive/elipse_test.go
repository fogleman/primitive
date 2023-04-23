package primitive

import (
	"fmt"
	"testing"

	"github.com/fogleman/gg"
)

func TestDraw(t *testing.T) {

	worker := NewWorker(imageToRGBA(createTestImage()))
	context := gg.NewContext(100, 100)
	context.SetRGBA255(224, 117, 232, 187)

	my_ellipse := Ellipse{worker, 50, 50, 10, 20, false}
	notify := NewTestStringNotifier()

	my_ellipse.Draw(context, 10, notify)
	context.Fill()

	contextState := Hash(context)

	// This value was pre-computed from static inputs
	if contextState != "99914b932bd37a50b983c5e7c90ae93b" {
		t.Error(fmt.Sprintf("Incorect state after Draw in Ellipse: %s", contextState))
	}
}

func TestSVGEllipse(t *testing.T) {
	worker := NewWorker(imageToRGBA(createTestImage2()))
	my_ellipse := Ellipse{worker, 123, 5234, 21, 55, false}
	SVG := my_ellipse.SVG("myAttrs")

	// This value was pre-computed from static inputs
	if SVG != "<ellipse myAttrs cx=\"123\" cy=\"5234\" rx=\"21\" ry=\"55\" />" {
		t.Error(fmt.Sprintf("Incorect SVG after SVG in Ellipse: %s", SVG))
	}
}

func TestRasterize(t *testing.T) {
	worker := NewWorker(imageToRGBA(createTestImage()))
	context := gg.NewContext(100, 100)
	context.SetRGBA255(224, 117, 232, 187)

	my_ellipse := Ellipse{worker, 40, 60, 12, 23, false}

	lines := my_ellipse.Rasterize()

	linesState := Hash(lines)

	// This value was pre-computed from static inputs
	if linesState != "958bacadcbf05504586f4922e5e999dd" {
		t.Error(fmt.Sprintf("Incorect state after Rasterize in Ellipse: %s", linesState))
	}
}
