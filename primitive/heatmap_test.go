package primitive

import (
	"fmt"
	"testing"
)

func TestAddHeatmap(t *testing.T) {
	testHeatmap := getTestHeatmap()

	testHeatmap.Add(getStaticScanLines())
	heatMapState := Hash(testHeatmap)

	if heatMapState != "37d9cb6898968c7d97d0e00307da9b73" {
		t.Error(fmt.Sprintf("Incorect state after Add in Heatmap: %s", heatMapState))
	}
}

func TestHeatmapImage(t *testing.T) {
	testHeatmap := getTestHeatmap()
	image := testHeatmap.Image(235.8698)
	imageState := Hash(image)
	if imageState != "7292d90fa7fef9c6914f8409a593412b" {
		t.Error(fmt.Sprintf("Incorect state after Imagein Heatmap: %s", imageState))
	}
}
