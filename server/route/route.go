package route

import (
	"image"
	"image/png"
	"net/http"

	"github.com/fogleman/primitive/build"
	"github.com/fogleman/primitive/primitive"
)

// Config contains core settings one can
// configure to setup the primitive file processing route.
type Config struct {
	MaxUploadMb int64
	FileKey     string
}

// PrimitiveRoute is an http route than can be used to consume
// an uploaded file and execute primitive on it.
func PrimitiveRoute(conf Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files.
		var maxUpload int64 = 10
		if conf.MaxUploadMb > 0 {
			maxUpload = conf.MaxUploadMb
		}
		r.ParseMultipartForm(maxUpload << 20)

		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file
		fileKey := "file"
		if conf.FileKey != "" {
			fileKey = conf.FileKey
		}
		file, _, err := r.FormFile(fileKey)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		buildConfig := build.Config{
			Input:      img,
			Alpha:      128,
			Count:      100,
			OutputSize: 1024,
			Mode:       4,
		}
		model := build.Build(buildConfig, func(_ *primitive.Model, _ int, _ int) {})
		finalImg := model.Context.Image()

		// return final primitive image
		w.Header().Set("Content-Type", "image/png")
		png.Encode(w, finalImg)
	}
}
