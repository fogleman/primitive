package primitive

import (
	"fmt"
	"testing"

	"github.com/fogleman/gg"
)

func TestDrawPolygon(t *testing.T) {

	worker := NewWorker(imageToRGBA(createTestImage()))
	context := gg.NewContext(100, 100)
	context.SetRGBA255(201, 97, 144, 185)

	myPolygon := Polygon{worker, 4, false,
		[]float64{10.20, 20.54, 63.31, 5.76},
		[]float64{23.64, 37.76, 10.11, 8},
	}
	notify := NewTestStringNotifier()

	myPolygon.Draw(context, 4, notify)
	context.Fill()

	contextState := Hash(context.Image())

	// This value was pre-computed from static inputs
	if contextState != "200b49359022c94a34393fba4bc18a03" {
		t.Error(fmt.Sprintf("Incorect state after Draw in Polygon: %s", contextState))
	}
}

func TestSVGPolygon(t *testing.T) {
	worker := NewWorker(imageToRGBA(createTestImage2()))
	myPolygon := Polygon{worker, 4, false,
		[]float64{10.20, 20.54, 63.31, 5.76},
		[]float64{23.64, 37.76, 10.11, 8},
	}
	SVG := myPolygon.SVG("myAttrs")

	// This value was pre-computed from static inputs
	if SVG !=
		"<polygon myAttrs points=\"10.200000,23.640000,20.540000,"+
			"37.760000,63.310000,10.110000,5.760000,8.000000\" />" {
		t.Error(fmt.Sprintf("Incorect SVG after SVG in Polygon: %s", SVG))
	}
}

func TestRasterizePolygon(t *testing.T) {
	worker := NewWorker(imageToRGBA(createTestImage()))
	context := gg.NewContext(100, 100)
	context.SetRGBA255(224, 117, 232, 187)

	myPolygon := Polygon{worker, 4, false,
		[]float64{15.20, 22.54, 62.31, 10.76},
		[]float64{23.64, 37.76, 10.11, 8},
	}

	lines := myPolygon.Rasterize()

	linesState := Hash(lines)

	// This value was pre-computed from static inputs
	if linesState != "a86c6ffb0b4617ed91149a4a4b82ce51" {
		t.Error(fmt.Sprintf("Incorect state after Rasterize in Polygon: %s", linesState))
	}
}
