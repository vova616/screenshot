package screenshot

import (
	"fmt"
	"github.com/AllenDang/w32"
	"image"
	"reflect"
	"unsafe"
)

func ScreenRect() (image.Rectangle, error) {
	hDC := w32.GetDC(0)
	if hDC == 0 {
		return image.Rectangle{}, fmt.Errorf("Could not Get primary display err:%d\n", w32.GetLastError())
	}
	defer w32.ReleaseDC(0, hDC)
	x := w32.GetDeviceCaps(hDC, w32.HORZRES)
	y := w32.GetDeviceCaps(hDC, w32.VERTRES)
	return image.Rect(0, 0, x, y), nil
}

func CaptureScreen() (*image.RGBA, func(), error) {
	r, e := ScreenRect()
	if e != nil {
		return nil, nil, e
	}
	return CaptureRect(r)
}

func CaptureRect(rect image.Rectangle) (*image.RGBA, func(), error) {
	hDC := w32.GetDC(0)
	if hDC == 0 {
		return nil, nil, fmt.Errorf("Could not Get primary display err:%d.\n", w32.GetLastError())
	}
	defer w32.ReleaseDC(0, hDC)

	m_hDC := w32.CreateCompatibleDC(hDC)
	if m_hDC == 0 {
		return nil, nil, fmt.Errorf("Could not Create Compatible DC err:%d.\n", w32.GetLastError())
	}
	noError := false
	defer func() {
		if !noError {
			w32.DeleteDC(m_hDC)
		}
	}()

	x, y := rect.Dx(), rect.Dy()

	bt := w32.BITMAPINFO{}
	bt.BmiHeader.BiSize = uint(reflect.TypeOf(bt.BmiHeader).Size())
	bt.BmiHeader.BiWidth = x
	bt.BmiHeader.BiHeight = -y
	bt.BmiHeader.BiPlanes = 1
	bt.BmiHeader.BiBitCount = 32
	bt.BmiHeader.BiCompression = w32.BI_RGB

	ptr := unsafe.Pointer(uintptr(0))

	m_hBmp := w32.CreateDIBSection(m_hDC, &bt, w32.DIB_RGB_COLORS, &ptr, 0, 0)
	if m_hBmp == 0 {
		return nil, nil, fmt.Errorf("Could not Create DIB Section err:%d.\n", w32.GetLastError())
	}
	if m_hBmp == w32.InvalidParameter {
		return nil, nil, fmt.Errorf("One or more of the input parameters is invalid while calling CreateDIBSection.\n")
	}
	defer func() {
		if !noError {
			w32.DeleteObject(w32.HGDIOBJ(m_hBmp))
		}
	}()

	obj := w32.SelectObject(m_hDC, w32.HGDIOBJ(m_hBmp))
	if obj == 0 {
		return nil, nil, fmt.Errorf("error occurred and the selected object is not a region err:%d.\n", w32.GetLastError())
	}
	if obj == 0xffffffff { //GDI_ERROR 
		return nil, nil, fmt.Errorf("GDI_ERROR while calling SelectObject err:%d.\n", w32.GetLastError())
	}
	defer func() {
		if !noError {
			w32.DeleteObject(obj)
		}
	}()

	//Note:BitBlt contains bad error handling, we will just assume it works and if it doesn't it will panic :x
	w32.BitBlt(m_hDC, 0, 0, x, y, hDC, rect.Min.X, rect.Min.Y, w32.SRCCOPY)

	var slice []byte
	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	hdrp.Data = uintptr(ptr)
	hdrp.Len = x * y * 4
	hdrp.Cap = x * y * 4

	for i := 0; i < len(slice); i += 4 {
		slice[i], slice[i+2] = slice[i+2], slice[i]
	}

	img := &image.RGBA{slice, 4 * x, image.Rect(0, 0, x, y)}
	noError = true
	return img, func() {
		w32.DeleteDC(m_hDC)
		w32.DeleteObject(w32.HGDIOBJ(m_hBmp))
		w32.DeleteObject(obj)
	}, nil
}
