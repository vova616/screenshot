package screenshot

import (
	"code.google.com/p/x-go-binding/xgb"
	"image"
	"os"
)

func ScreenRect() (image.Rectangle, error) {
	c, err := xgb.Dial(os.Getenv("DISPLAY"))
	if err != nil {
		return image.Rectangle{}, err
	}
	defer c.Close()
	x := c.DefaultScreen().WidthInPixels
	y := c.DefaultScreen().HeightInPixels

	return image.Rect(0, 0, int(x), int(y)), nil
}

func CaptureScreen() (*image.RGBA, func(), error) {
	r, e := ScreenRect()
	if e != nil {
		return nil, nil, e
	}
	return CaptureRect(r)
}

func CaptureRect(rect image.Rectangle) (*image.RGBA, func(), error) {
	c, err := xgb.Dial(os.Getenv("DISPLAY"))
	if err != nil {
		return nil, nil, err
	}
	noError := false
	defer func() {
		if !noError {
			c.Close()
		}
	}()

	x, y := rect.Dx(), rect.Dy()

	xImg, err := c.GetImage(xgb.ImageFormatZPixmap, c.DefaultScreen().Root, uint16(rect.Min.X), uint16(rect.Min.Y), uint16(x), uint16(y), 0xffffffff)
	if err != nil {
		return nil, nil, err
	}

	data := xImg.Data

	for i := 0; i < len(data); i += 4 {
		data[i], data[i+2], data[i+3] = data[i+2], data[i], 255
	}

	img := &image.RGBA{data, 4 * x, image.Rect(0, 0, x, y)}
	noError = true
	return img, func() { c.Close() }, nil
}
