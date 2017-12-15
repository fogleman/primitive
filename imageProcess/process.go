package imageProcess

import (
	"fmt"
	"github.com/mojocn/primitive/primitive"
	"github.com/nfnt/resize"
	"log"
	"math/rand"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//ProccessImage
//mode 0=combo 1=triangle 2=rect 3=ellipse 4=circle 5=rotatedrect 6=beziers 7=rotatedellipse 8=polygon
//frameCount
func ProccessImage(inputImgPath, outputImagePath string, mode, shapeCount, frameCount, outputSize int) {
	// seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())
	// determine worker count

	// read input image
	input, _ := primitive.LoadImage(inputImgPath)

	// scale down input image if needed
	size := uint(outputSize)
	if size > 0 {
		input = resize.Thumbnail(size, size, input, resize.Bilinear)
	}

	// determine background color
	bg := primitive.MakeColor(primitive.AverageImageColor(input))

	// run algorithm
	model := primitive.NewModel(input, bg, outputSize, runtime.NumCPU())

	ext := strings.ToLower(filepath.Ext(outputImagePath))
	percent := strings.Contains(outputImagePath, "%")

	frameDelta := shapeCount / frameCount

	frameDeltaBegin := shapeCount % frameCount
	//fmt.Println(frameDelta,"/n")
	start := time.Now()

	for i := 0; i < shapeCount; i++ {
		t := time.Now()

		// write output image(s)
		n := model.Step(primitive.ShapeType(mode), 0, 0)
		nps := primitive.NumberString(float64(n) / time.Since(t).Seconds())
		elapsed := time.Since(start).Seconds()
		fmt.Printf("%d: t=%.3f, score=%.6f, n=%d, n/s=%s\n", i, elapsed, model.Score, n, nps)

		isSaveFrame := percent && ext != ".gif"
		isSaveFrame = isSaveFrame && (i+1)%frameDelta == frameDeltaBegin

		isLastFrame := i == (shapeCount - 1)
		if isSaveFrame || isLastFrame {
			//设置output frame 函数
			path := outputImagePath
			if percent {
				path = fmt.Sprintf(outputImagePath, i/frameDelta+1)
			}

			switch ext {
			default:
				check(fmt.Errorf("unrecognized file extension: %s", ext))
			case ".png":
				check(primitive.SavePNG(path, model.Context.Image()))
			case ".jpg", ".jpeg":
				check(primitive.SaveJPG(path, model.Context.Image(), 95))
			case ".svg":
				check(primitive.SaveFile(path, model.SVG()))
			case ".gif":
				frames := model.FramesForGif(frameCount)
				fmt.Printf("%d,frames count", len(model.Shapes))
				check(primitive.SaveGIFImageMagick(path, frames, 50, 250))
			case ".mp4":
				frames := model.Frames(0.00000001)
				check(primitive.SaveMp4(path, frames))
			}
		}
	}
}

//ffmpeg -framerate 30 -i input_%05d.png -c:v libx264 -crf 23 -pix_fmt yuv420p output.mp4
//https://stackoverflow.com/questions/13163106/ffmpeg-converting-image-sequence-to-video-results-in-blank-video
//https://github.com/leafo/gifserver/blob/master/gifserver/gif.go

// ffmpeg -i "$pattern" -pix_fmt yuv420p -vf 'scale=trunc(in_w/2)*2:trunc(in_h/2)*2' "${out_base}.mp4"
func ConvertPngFramesToMP4(dir, imageNamePatern string) (string, error) {
	fmt.Print("Encoding ", dir, " to mp4")

	outFname := "out.mp4"
	cmd := exec.Command("ffmpeg",
		"-i", imageNamePatern,
		"-pix_fmt", "yuv420p",
		"-vf", "scale=trunc(in_w/2)*2:trunc(in_h/2)*2",
		outFname)

	cmd.Dir = dir
	err := cmd.Run()

	if err != nil {
		return "", err
	}

	return path.Join(dir, outFname), nil
}
