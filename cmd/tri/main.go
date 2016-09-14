package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/fogleman/tri"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
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
	model := tri.NewModel(im)
	model.Run()
}
