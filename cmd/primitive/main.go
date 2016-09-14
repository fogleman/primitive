package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/fogleman/primitive"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("Usage: primitive input")
		return
	}
	im, err := primitive.LoadImage(args[0])
	if err != nil {
		panic(err)
	}
	model := primitive.NewModel(im)
	model.Run()
}
