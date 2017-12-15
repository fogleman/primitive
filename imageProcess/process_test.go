package imageProcess

import (
	"testing"
)

func TestProccessImage(t *testing.T) {
	inputPath := "/Users/ericzhou/go/src/github.com/mojocn/primitive/examples/lenna.png"
	outputPath := "/Users/ericzhou/go/src/github.com/mojocn/primitive/examples/lenna.mp4"

	ProccessImage(inputPath, outputPath, 0, 100, 4, 1024)
}
