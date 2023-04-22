package primitive

import (
	"image"
	"image/color"
	"math/rand"
	"time"
)

func createTestModel() *Model {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	size := rnd.Intn(400) + 10 // Avoid passing 0 to Intn
	width_modifier := rnd.Intn(int(size / 10))
	height_modifier := rnd.Intn(int(size / 10))

	modify_width_negative := rnd.Intn(1)
	modify_height_negative := rnd.Intn(1)

	if modify_width_negative == 1 {
		width_modifier = width_modifier * -1
	}

	if modify_height_negative == 1 {
		height_modifier = height_modifier * -1
	}
	width := size + width_modifier
	height := size + height_modifier

	rect := image.Rect(0, 0, width, height)
	pix := make([]uint8, rect.Dx()*rect.Dy()*4)

	rand.Read(pix)

	testingImage := &image.NRGBA{
		Pix:    pix,
		Stride: rect.Dx() * 4,
		Rect:   rect,
	}

	// generate sizes which are significant fractions of the original image size
	// size_modifier := rnd.Intn(size / 2)
	// scaled_size := size - size_modifier

	// generate a small random number of workers
	num_workers := rnd.Intn(4)

	// generate a random background color
	r := rnd.Intn(255)
	g := rnd.Intn(255)
	b := rnd.Intn(255)

	nrgba_color := color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
	background_color := MakeColor(nrgba_color)

	// run the function under test with generated values
	return NewModel(testingImage, background_color, size, num_workers)
}

// Notifier for unit testing
type Notifier interface {
	Notify(message string)
}

type NullNotify struct {
}

func (n *NullNotify) Notify(message string) {
	//do nothing
}

func NewTestStringNotifier() *testStringNotifier {

	notifier := new(testStringNotifier)
	notifier.messages = make(map[string]int)

	return notifier
}

type testStringNotifier struct {
	messages map[string]int
}

func (tsn *testStringNotifier) Notify(message string) {
	tsn.messages[message] += 1
}
