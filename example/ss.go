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
	f, e := os.Create("./ss.png")
	if e != nil {
		panic(e)
	}
	png.Encode(f, img)
	f.Close()
}
