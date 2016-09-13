package main

import (
	"flag"
	"fmt"

	"github.com/fogleman/tri"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: tri input")
		return
	}
	im, err := tri.LoadImage(args[0])
	if err != nil {
		panic(err)
	}
	tri.Run(im)
}
