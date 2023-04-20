package primitive

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"io"
	"strings"
	"testing"
)

const jpgData = `
/9j/4AAQSkZJRgABAQIAHAAcAAD/2wBDABALDA4MChAODQ4SERATGCgaGBYWGDEjJR0oOjM9PDkzODdA
SFxOQERXRTc4UG1RV19iZ2hnPk1xeXBkeFxlZ2P/2wBDARESEhgVGC8aGi9jQjhCY2NjY2NjY2NjY2Nj
Y2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2NjY2P/wAARCABnAJYDASIAAhEBAxEB/8QA
HwAAAQUBAQEBAQEAAAAAAAAAAAECAwQFBgcICQoL/8QAtRAAAgEDAwIEAwUFBAQAAAF9AQIDAAQRBRIh
MUEGE1FhByJxFDKBkaEII0KxwRVS0fAkM2JyggkKFhcYGRolJicoKSo0NTY3ODk6Q0RFRkdISUpTVFVW
V1hZWmNkZWZnaGlqc3R1dnd4eXqDhIWGh4iJipKTlJWWl5iZmqKjpKWmp6ipqrKztLW2t7i5usLDxMXG
x8jJytLT1NXW19jZ2uHi4+Tl5ufo6erx8vP09fb3+Pn6/8QAHwEAAwEBAQEBAQEBAQAAAAAAAAECAwQF
BgcICQoL/8QAtREAAgECBAQDBAcFBAQAAQJ3AAECAxEEBSExBhJBUQdhcRMiMoEIFEKRobHBCSMzUvAV
YnLRChYkNOEl8RcYGRomJygpKjU2Nzg5OkNERUZHSElKU1RVVldYWVpjZGVmZ2hpanN0dXZ3eHl6goOE
hYaHiImKkpOUlZaXmJmaoqOkpaanqKmqsrO0tba3uLm6wsPExcbHyMnK0tPU1dbX2Nna4uPk5ebn6Onq
8vP09fb3+Pn6/9oADAMBAAIRAxEAPwDlwKMD0pwzSiuK57QzGDxS7D6in8Y5ximnAPUfSlcq4m3ilUYp
2OKXHvRcVxnTtS7c07HNFK4DQPakC4PNOA+tOx70XAjK/So5gBGP94fzqfvUVx/qxx/EP51UXqRP4WSE
cmgjilP3jSEZqS0IO/NGDnpUiocDg/McDjvV6HTPOdVWYgsM5KcfzzQ2JySM2jp6VYu7SWzmMUwG4cgj
kMPUVBjjtTGtRu0Zopw+lFFxhinrGzuqqMsxAA9yaXFSRv5cqSEcIwYj6GpuZ30O30fSLKzhUpbpNMv3
5XGTn29BV28jt7pPLuIVljPBBFVreYx+VbqAjycgt3x14zRcNOxGyVFHQkIc/wA61exyKLbuzjdZ046d
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
`

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
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpgData))
	reader2 := base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpgData))
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
	reader = base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpgData))
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

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpgData))
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

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpgData))
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

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpgData))
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

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(jpgData))
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

	for _, numstr := range cases {
		if numstr.str != NumberString(numstr.num) {
			t.Error(fmt.Sprintf("bad conversion in NumberString on %f -> %s", numstr.num, numstr.str))
		}
	}
}
