package imageProcess

import (
	"testing"
)

func TestProccessImage(t *testing.T) {
	inputPath := "/Users/ericzhou/go/src/github.com/mojocn/primitive/examples/monalisa.png"
	outputPath := "/Users/ericzhou/go/src/github.com/mojocn/primitive/examples/monalisa%2d.png"

	ProccessImage(inputPath, outputPath, 0, 1024, 8, 1024)
}
