package screenshot

import (
	"image"
)

func ScreenRect() (image.Rectangle, error) {
	panic("ScreenRect is not supported on this platform.")
}

func CaptureScreen() (*image.RGBA, error) {
	panic("CaptureScreen is not supported on this platform.")
}

func CaptureRect(rect image.Rectangle) (*image.RGBA, error) {
	panic("CaptureRect is not supported on this platform.")
}
