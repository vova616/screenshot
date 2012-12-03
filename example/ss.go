package main

import (
	"github.com/vova616/screenshot"
	"image/png"
	"os"
)

func main() {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		panic(err)
	}
	f, err := os.Create("./ss.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	f.Close()
}
