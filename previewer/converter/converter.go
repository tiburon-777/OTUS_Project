package converter

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"

	"github.com/nfnt/resize"
)

type Image struct {
	image.Image
}

func SelectType(width int, height int, b []byte) ([]byte, error) {
	contentType := http.DetectContentType(b)
	var decode func(reader *bytes.Reader) (image.Image, error)
	var encode func(m image.Image) ([]byte, error)
	switch contentType {
	case "image/png":
		decode = func(reader *bytes.Reader) (image.Image, error) {
			return png.Decode(bytes.NewReader(b))
		}
		encode = func(m image.Image) (res []byte, err error) {
			err = png.Encode(bytes.NewBuffer(res), m)
			return res, err
		}
	case "image/gif":
		decode = func(reader *bytes.Reader) (image.Image, error) {
			return gif.Decode(bytes.NewReader(b))
		}
		encode = func(m image.Image) (res []byte, err error) {
			err = gif.Encode(bytes.NewBuffer(res), m, nil)
			return res, err
		}
	default:
		decode = func(reader *bytes.Reader) (image.Image, error) {
			return jpeg.Decode(bytes.NewReader(b))
		}
		encode = func(m image.Image) (res []byte, err error) {
			err = jpeg.Encode(bytes.NewBuffer(res), m, nil)
			return res, err
		}
	}
	i, err := decode(bytes.NewReader(b))
	m := Image{i}
	if err != nil {
		return nil, err
	}
	err = m.convert(width, height)
	if err != nil {
		return nil, err
	}
	res, err := encode(m)
	return res, err
}

func (img *Image) convert(width int, height int) error {
	widthOrig := img.Bounds().Max.X
	heightOrig := img.Bounds().Max.Y
	sfOriginal := sizeFactor(widthOrig, heightOrig)
	sfNew := sizeFactor(width, height)
	switch {
	case sfOriginal > sfNew:
		// Ресайз по одной высоте и кроп по ширине следом
		// Определение ширины кропа.
		img.resize(int(float64(widthOrig)*sfOriginal), height)
		if err := img.crop(image.Point{X: (widthOrig - width) / 2, Y: 0}, image.Point{X: (widthOrig-width)/2 + width, Y: height}); err != nil {
			return err
		}
	case sfOriginal == sfNew:
		img.resize(width, height)
	case sfOriginal < sfNew:
		// Ресайз по одной ширине и кроп по высоте следом
		img.resize(width, int(float64(heightOrig)*sfOriginal))
		if err := img.crop(image.Point{X: 0, Y: (heightOrig - height) / 2}, image.Point{X: width, Y: (heightOrig-height)/2 + height}); err != nil {
			return err
		}
	}
	return nil
}

func (img *Image) resize(width, height int) {
	img.Image = resize.Resize(uint(width), uint(height), img, resize.Bicubic)
}

func (img *Image) crop(p1 image.Point, p2 image.Point) error {
	if img == nil {
		return errors.New("corrupted image")
	}
	if p1.X < 0 || p1.Y < 0 || p2.X < 0 || p2.Y < 0 {
		return errors.New("not valid corner points")
	}
	b := image.Rect(0, 0, p2.X, p2.Y)
	resImg := image.NewRGBA(b)
	draw.Draw(resImg, b, img, p1, draw.Src)
	return nil
}

func sizeFactor(width int, height int) float64 {
	return float64(width) / float64(height)
}
