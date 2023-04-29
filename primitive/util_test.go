package primitive

import (
	"encoding/base64"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"io"
	"strings"
	"testing"
)

type mockClosableReader struct {
	reader io.Reader
	state  string
	path   string
}

type mockClosableWriter struct {
	writer io.Writer
	state  string
	path   string
}

func (cr *mockClosableReader) Read(p []byte) (n int, err error) {
	return cr.reader.Read(p)
}

func (cw *mockClosableWriter) Write(p []byte) (n int, err error) {
	return cw.writer.Write(p)
}

func (cw *mockClosableWriter) WriteString(s string) (int, error) {
	cw.writer.Write([]byte(s))
	return len(s), nil
}

type mockRunable struct {
	command string
	args    []string
	err     error
}

func (mr *mockRunable) Run() error {
	return mr.err
}

func (cr *mockClosableReader) Close() error {
	cr.state = "closed"
	return nil
}

func (cw *mockClosableWriter) Close() error {
	cw.state = "closed"
	return nil
}

func TestLoadImage(t *testing.T) {

	// Path of "-" takes image from stdin
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpgImageTestData))
	reader2 := base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpgImageTestData))
	compareImage, _, _ := image.Decode(reader2)
	osStdin = func() io.Reader { return reader }

	responseImage, err := LoadImage("-")

	if !responseImage.Bounds().Eq(compareImage.Bounds()) ||
		responseImage.ColorModel() != compareImage.ColorModel() {

		t.Error("Error decoding from stdin in Load Image")
	}

	if err != nil {
		t.Error("Unexpected non-null value in err in Load Image")
	}

	// filepath is opened with osOpen
	testReader := mockClosableReader{reader: strings.NewReader(""), state: "open"}
	osOpen = func(p string) (closableReader, error) {
		testReader.path = p
		return &testReader, nil
	}
	_, err = LoadImage("This is a path!")

	if testReader.path != "This is a path!" {
		t.Error("Path not properly set in Load Image")
	}

	// file is eventually closed
	if testReader.state != "closed" {
		t.Error("File was not closed in Load Image")
	}

	// osOpen error is handled
	osOpen = func(p string) (closableReader, error) {
		return nil, errors.New("Test Error")
	}

	_, err = LoadImage("This is a path!")

	if err.Error() != "Test Error" {
		t.Error("osOpen error not passed back with input from stdin in Load Image")
	}

	// imageDecode is called
	reader = base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpgImageTestData))
	testReader = mockClosableReader{reader: reader, state: "open"}

	osOpen = func(p string) (closableReader, error) {
		return &testReader, nil
	}

	responseImage, err = LoadImage("testPath")

	if !responseImage.Bounds().Eq(compareImage.Bounds()) ||
		responseImage.ColorModel() != compareImage.ColorModel() {
		t.Error("Error decoding image from path in Load Image")
	}

	// imageDecode error is handled
	osOpen = func(p string) (closableReader, error) {
		return nil, errors.New("Test Error2")
	}

	responseImage, err = LoadImage("testPath")

	if err.Error() != "Test Error2" {
		t.Error("osOpen error not passed back in Load Image")
	}

}

func TestSaveFile(t *testing.T) {

	// Path of "-" saves to stdout
	writer := new(strings.Builder)
	testWriter := mockClosableWriter{writer: writer, state: "open"}
	osStdout = func() io.Writer { return &testWriter }
	SaveFile("-", "test content")
	if writer.String() != "test content" {
		t.Error("SafeFile not writing to stdout in SaveFile")
	}

	// filepath is opened with osOpen
	osCreate = func(p string) (closableWriter, error) {
		testWriter.path = p
		return &testWriter, nil
	}

	writer = new(strings.Builder)
	testWriter = mockClosableWriter{writer: writer, state: "open"}

	SaveFile("test path", "test content2")

	if testWriter.path != "test path" {
		t.Error("path not set in SaveFile")
	}

	// file is eventually closed
	if testWriter.state != "closed" {
		t.Error("File was not closed in Save File")
	}

	// contents are written
	if writer.String() != "test content2" {
		t.Error("SaveFile did not write contents: " + writer.String())
	}

	// osCreate error is handled
	osCreate = func(p string) (closableWriter, error) {
		testWriter.path = p
		return nil, errors.New("Test Error3")
	}

	err := SaveFile("test path", "test content3")

	if err.Error() != "Test Error3" {
		t.Error("SaveFile did not handle osCreate error")
	}
}

func TestSavePNG(t *testing.T) {

	// osCreate error is handled
	writer := new(strings.Builder)
	testWriter := mockClosableWriter{writer: writer, state: "open"}

	osCreate = func(p string) (closableWriter, error) {
		testWriter.path = p
		return nil, errors.New("Test Error4")
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpgImageTestData))
	sampleImage, _, _ := image.Decode(reader)
	err := SavePNG("test path", sampleImage)

	if err.Error() != "Test Error4" {
		t.Error("Error not handled for osCreate in SavePNG")
	}

	// file is eventually closed
	writer = new(strings.Builder)
	testWriter = mockClosableWriter{writer: writer, state: "open"}

	osCreate = func(p string) (closableWriter, error) {
		testWriter.path = p
		return &testWriter, nil
	}

	SavePNG("test path", sampleImage)

	if testWriter.state != "closed" {
		t.Error("File was not closed in SavePNG")
	}

	//PNGEncode is called
	writer = new(strings.Builder)
	wasCalled := false
	testWriter = mockClosableWriter{writer: writer, state: "open"}
	pngEncode = func(wr io.Writer, img image.Image) error {
		*&wasCalled = true
		return nil
	}

	SavePNG("test path", sampleImage)

	if wasCalled != true {
		t.Error("PNGEncode was never falled in SavePNG")
	}

}

func TestSaveJPG(t *testing.T) {

	// osCreate error is handled
	writer := new(strings.Builder)
	testWriter := mockClosableWriter{writer: writer, state: "open"}

	osCreate = func(p string) (closableWriter, error) {
		testWriter.path = p
		return nil, errors.New("Test Error5")
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpgImageTestData))
	sampleImage, _, _ := image.Decode(reader)
	err := SaveJPG("test path", sampleImage, 10)

	if err.Error() != "Test Error5" {
		t.Error("Error not handled for osCreate in SaveJPG")
	}

	// file is eventually closed
	writer = new(strings.Builder)
	testWriter = mockClosableWriter{writer: writer, state: "open"}

	osCreate = func(p string) (closableWriter, error) {
		testWriter.path = p
		return &testWriter, nil
	}

	SaveJPG("test path", sampleImage, 10)

	if testWriter.state != "closed" {
		t.Error("File was not closed in SaveJPG")
	}
	//jpegEncode is called
	writer = new(strings.Builder)
	wasCalled := false
	testWriter = mockClosableWriter{writer: writer, state: "open"}

	jpegEncode = func(wr io.Writer, img image.Image, ops *jpeg.Options) error {
		*&wasCalled = true
		return nil
	}

	SaveJPG("test path", sampleImage, 1)

	if wasCalled != true {
		t.Error("jpegEncode was never falled in SaveJPG")
	}
}

func TestSaveGIF(t *testing.T) {

	// osCreate error is handled
	writer := new(strings.Builder)
	testWriter := mockClosableWriter{writer: writer, state: "open"}

	osCreate = func(p string) (closableWriter, error) {
		testWriter.path = p
		return nil, errors.New("Test Error6")
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpgImageTestData))
	sampleImage, _, _ := image.Decode(reader)
	var sampleImageArr = make([]image.Image, 1)
	sampleImageArr[0] = sampleImage
	err := SaveGIF("test path", sampleImageArr, 1, 1)

	if err.Error() != "Test Error6" {
		t.Error("Error not handled for osCreate in SaveGIF")
	}

	// file is eventually closed
	writer = new(strings.Builder)
	testWriter = mockClosableWriter{writer: writer, state: "open"}

	osCreate = func(p string) (closableWriter, error) {
		testWriter.path = p
		return &testWriter, nil
	}

	SaveGIF("test path", sampleImageArr, 1, 1)

	if testWriter.state != "closed" {
		t.Error("File was not closed in SavGIF")
	}

	//gifEncodeAll is called
	writer = new(strings.Builder)
	wasCalled := false
	testWriter = mockClosableWriter{writer: writer, state: "open"}
	gifEncodeAll = func(wr io.Writer, g *gif.GIF) error {
		*&wasCalled = true
		return nil
	}

	SaveGIF("test path", sampleImageArr, 1, 1)

	if wasCalled != true {
		t.Error("PNGEncode was never falled in SavePNG")
	}
}

func TestSaveGIFImageMagick(t *testing.T) {

	// osCreate error is handled
	writer := new(strings.Builder)
	testWriter := mockClosableWriter{writer: writer, state: "open"}

	osCreate = func(p string) (closableWriter, error) {
		testWriter.path = p
		return nil, errors.New("Test Error7")
	}

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpgImageTestData))
	sampleImage, _, _ := image.Decode(reader)
	var sampleImageArr = make([]image.Image, 1)
	sampleImageArr[0] = sampleImage
	err := SaveGIFImageMagick("test path", sampleImageArr, 1, 1)

	if err.Error() != "Test Error7" {
		t.Error("Error not handled for osCreate in SaveGIFImageMagick")
	}

	// exec error is handled and exec command is called
	mrError := errors.New("test exec error")
	command := mockRunable{err: mrError}
	execCommand = func(cmd string, arg ...string) runable {
		return &command
	}

	osCreate = func(p string) (closableWriter, error) {
		testWriter.path = p
		return &testWriter, nil
	}

	err = SaveGIFImageMagick("test path", sampleImageArr, 1, 1)

	if err.Error() != "test exec error" {
		t.Error("Error not handled for cmd.Run() in SaveGIFImageMagick")
	}

	// file is eventually closed
	writer = new(strings.Builder)
	testWriter = mockClosableWriter{writer: writer, state: "open"}

	osCreate = func(p string) (closableWriter, error) {
		testWriter.path = p
		return &testWriter, nil
	}

	SaveGIFImageMagick("test path", sampleImageArr, 1, 1)

	if testWriter.state != "closed" {
		t.Error("File was not closed in SavGIFImageMagick")
	}

}

type numstr struct {
	num float64
	str string
}

func TestNumberString(t *testing.T) {
	cases := make([]numstr, 26)

	cases[0] = numstr{num: 1, str: "1.00B"}
	cases[1] = numstr{num: 1.1, str: "1.10B"}
	cases[2] = numstr{num: 999, str: "999.00B"}
	cases[3] = numstr{num: 999.4, str: "999.40B"}
	cases[4] = numstr{num: 1000, str: "1.00KB"}
	cases[5] = numstr{num: 1020, str: "1.02KB"}
	cases[6] = numstr{num: 1350, str: "1.35KB"}
	cases[7] = numstr{num: 9999, str: "10.00KB"}
	cases[8] = numstr{num: 9999.999, str: "10.00KB"}
	cases[9] = numstr{num: 999999, str: "1.00MB"}
	cases[10] = numstr{num: 999999.2323234, str: "1.00MB"}
	cases[11] = numstr{num: 1000000, str: "1.00MB"}
	cases[12] = numstr{num: 1000000.1, str: "1.00MB"}
	cases[13] = numstr{num: 1190000.18, str: "1.19MB"}
	cases[14] = numstr{num: 10400001.17, str: "10.40MB"}
	cases[15] = numstr{num: 10999000, str: "11.00MB"}
	cases[16] = numstr{num: 10449000.87, str: "10.45MB"}
	cases[17] = numstr{num: -1000, str: "-1.00KB"}
	cases[18] = numstr{num: -1000000.1, str: "-1.00MB"}
	cases[19] = numstr{num: -999999999, str: "-1.00GB"}
	cases[20] = numstr{num: 3123123123123, str: "3.12TB"}
	cases[21] = numstr{num: 2342342342342345, str: "2.34PB"}
	cases[22] = numstr{num: 3454876748574948496, str: "3.45EB"}
	cases[23] = numstr{num: 6985776947869476957407, str: "6.99ZB"}
	cases[24] = numstr{num: 7590732498760209384908757, str: "7.59YB"}
	cases[25] = numstr{num: 999999999999999999999999999999, str: "1000000.00YB"}

	for _, nstr := range cases {
		if nstr.str != NumberString(nstr.num) {
			t.Errorf("bad conversion in NumberString on %f -> %s", nstr.num, nstr.str)
		}
	}
}

// Construct a random image, adding up all the r,g, and b values
// to find the average, then pass the image into AverageImageColor
// to make sure the averages match.
func TestAverageImageColor(t *testing.T) {

	testingImage := createTestImage()
	testAVG := AverageImageColor(testingImage)

	if uint8(testImageAverageRed()) != testAVG.R {
		t.Errorf("wrong average red in AverageImageColor: %d", testAVG.R)
	}

	if uint8(testImageAverageGreen()) != testAVG.G {
		t.Errorf("wrong average green in AverageImageColor: %d", testAVG.G)
	}

	if uint8(testImageAverageBlue()) != testAVG.B {
		t.Errorf("wrong average blue in AverageImageColor: %d", testAVG.B)
	}
}
