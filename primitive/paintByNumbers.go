package primitive

// color type
type clr struct {
	name string
	c    Color
}

//REQ 2.0
//maps type clr with number
colorMap := make(map[clr]int)


//REQ 2.2 2.2.1
//makes the fill transparent
func clearFill(c *Shape) {

}

//Gets the fill
func getFill(c *Shape) {
	//Generates clr (color) based on fill color
	//Add to map
}

//REQ 2.1 2.1.1 2.1.2 2.1.2.1
//Writes the number in generated shape
func annotate() {
	//Gets fill then clearsFill
	//Writes int key corresponding color value
	//Note: Look into DrawString
	Font f = new Font(8)
}
