package primitive

import (
	"fmt"
	"testing"
)

type testingColor struct {
	r uint32
	g uint32
	b uint32
	a uint32
}

func (c *testingColor) RGBA() (r, g, b, a uint32) {
	c.r = 62535
	c.g = 0
	c.b = 65400
	c.a = 22456

	return c.r, c.g, c.b, c.a
}

func TestMakeColor(t *testing.T) {
	testColor := new(testingColor)
	color := MakeColor(testColor)

	if color.R != int(testColor.r/257) ||
		color.G != int(testColor.g/257) ||
		color.B != int(testColor.b/257) ||
		color.A != int(testColor.a/257) {
		t.Error("bad color computation in MakeColor")
	}
}

type hexStrColor struct {
	hexstr string
	r      int
	g      int
	b      int
	a      int
}

func TestMakeHexColor(t *testing.T) {

	cases := make([]hexStrColor, 25)

	// 3 digit hex values: a hex value of a will become aa, which is 170 decimal.
	// a hex value of 2 will become 22 which is 34 decimal ..etc.
	// the value of data member a is always 255 if ommited (if ther is no fourth value)
	cases[0] = hexStrColor{hexstr: "#aaa", r: 170, g: 170, b: 170, a: 255}
	cases[1] = hexStrColor{hexstr: "bc3", r: 187, g: 204, b: 51, a: 255}
	cases[2] = hexStrColor{hexstr: "fff", r: 255, g: 255, b: 255, a: 255}
	cases[3] = hexStrColor{hexstr: "#fff", r: 255, g: 255, b: 255, a: 255}
	cases[4] = hexStrColor{hexstr: "000", r: 0, g: 0, b: 0, a: 255}

	// 4 digit hex values: 2 becomes 22 the same way as 3 digit hex values,
	// but member value a is sepcified
	cases[5] = hexStrColor{hexstr: "#aaaa", r: 170, g: 170, b: 170, a: 170}
	cases[6] = hexStrColor{hexstr: "bc3c", r: 187, g: 204, b: 51, a: 204}
	cases[7] = hexStrColor{hexstr: "fffd", r: 255, g: 255, b: 255, a: 221}
	cases[8] = hexStrColor{hexstr: "#fff9", r: 255, g: 255, b: 255, a: 153}
	cases[9] = hexStrColor{hexstr: "0000", r: 0, g: 0, b: 0, a: 0}

	// 6 digit hex values: direct translation without transforming 2 into 22,
	// member value a is always 255
	cases[10] = hexStrColor{hexstr: "#aaaaaa", r: 170, g: 170, b: 170, a: 255}
	cases[11] = hexStrColor{hexstr: "bc3cbc", r: 188, g: 60, b: 188, a: 255}
	cases[12] = hexStrColor{hexstr: "fffdff", r: 255, g: 253, b: 255, a: 255}
	cases[13] = hexStrColor{hexstr: "#fff987", r: 255, g: 249, b: 135, a: 255}
	cases[14] = hexStrColor{hexstr: "000000", r: 0, g: 0, b: 0, a: 255}

	// 8 digit hex values. Member value a is specified
	cases[15] = hexStrColor{hexstr: "#aaaaaabc", r: 170, g: 170, b: 170, a: 188}
	cases[16] = hexStrColor{hexstr: "bc3cbcdf", r: 188, g: 60, b: 188, a: 223}
	cases[17] = hexStrColor{hexstr: "fffdff3c", r: 255, g: 253, b: 255, a: 60}
	cases[18] = hexStrColor{hexstr: "#fff98701", r: 255, g: 249, b: 135, a: 1}
	cases[19] = hexStrColor{hexstr: "00000000", r: 0, g: 0, b: 0, a: 0}

	// Out of range values that should fall through
	cases[20] = hexStrColor{hexstr: "i3gunklsjdfh", r: 0, g: 0, b: 0, a: 255}
	cases[21] = hexStrColor{hexstr: "", r: 0, g: 0, b: 0, a: 255}
	cases[22] = hexStrColor{hexstr: "56454545454545454", r: 0, g: 0, b: 0, a: 255}
	cases[23] = hexStrColor{hexstr: "########", r: 0, g: 0, b: 0, a: 255}
	cases[24] = hexStrColor{hexstr: "#a", r: 0, g: 0, b: 0, a: 255}

	for _, testcase := range cases {
		color := MakeHexColor(testcase.hexstr)
		if color.R != testcase.r ||
			color.G != testcase.g ||
			color.B != testcase.b ||
			color.A != testcase.a {
			t.Error(fmt.Sprintf("bad conversion in MakeHexColor on %s", testcase.hexstr))
		}
	}
}
