package main

import (
	"github.com/vova616/screenshot"
	"image/png"
	"os"
)

func main() {
	img, closer, err := screenshot.CaptureScreen()
	if err != nil {
		panic(err)
	}
	defer closer()
	f, e := os.Create("./ss.png")
	if e != nil {
		panic(e)
	}
	png.Encode(f, img)
	f.Close()
}
