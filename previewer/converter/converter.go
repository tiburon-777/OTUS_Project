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
	var decode func(b []byte) (image.Image, error)
	var encode func(m image.Image) ([]byte, error)
	switch contentType {
	case "image/jpeg":
		decode = func(b []byte) (image.Image, error) {
			tb := bytes.NewBuffer(b)
			return jpeg.Decode(tb)
		}
		encode = func(m image.Image) ([]byte, error) {
			res := bytes.NewBuffer([]byte{})
			err := jpeg.Encode(res, m, &jpeg.Options{Quality: 80})
			return res.Bytes(), err
		}
	case "image/png":
		decode = func(b []byte) (image.Image, error) {
			tb := bytes.NewBuffer(b)
			return png.Decode(tb)
		}
		encode = func(m image.Image) ([]byte, error) {
			res := bytes.NewBuffer([]byte{})
			err := png.Encode(res, m)
			return res.Bytes(), err
		}
	case "image/gif":
		decode = func(b []byte) (image.Image, error) {
			tb := bytes.NewBuffer(b)
			return gif.Decode(tb)
		}
		encode = func(m image.Image) ([]byte, error) {
			res := bytes.NewBuffer([]byte{})
			err := gif.Encode(res, m, nil)
			return res.Bytes(), err
		}
	default:
		decode = func(b []byte) (image.Image, error) {
			return nil, errors.New("unknown format")
		}
		encode = func(m image.Image) ([]byte, error) {
			return nil, errors.New("unknown format")
		}
	}
	i, err := decode(b)
	if err != nil {
		return nil, err
	}
	m := Image{i}
	if err = m.convert(width, height); err != nil {
		return nil, err
	}
	return encode(m.Image)
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
		calcWidth := int(float64(height) * sfOriginal)
		img.resize(calcWidth, height)
		if err := img.crop(image.Point{X: (calcWidth - width) / 2, Y: 0}, image.Point{X: (calcWidth-width)/2 + width, Y: height}); err != nil {
			return err
		}
	case sfOriginal == sfNew:
		img.resize(width, height)
	case sfOriginal < sfNew:
		// Ресайз по одной ширине и кроп по высоте следом
		calcHeight := int(float64(width) / sfOriginal)
		img.resize(width, calcHeight)
		if err := img.crop(image.Point{X: 0, Y: (calcHeight - height) / 2}, image.Point{X: width, Y: (calcHeight-height)/2 + height}); err != nil {
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
	b := image.Rect(0, 0, p2.X-p1.X, p2.Y-p1.Y)
	resImg := image.NewRGBA(b)
	draw.Draw(resImg, b, img.Image, p1, draw.Src)
	img.Image = resImg
	return nil
}

func sizeFactor(width int, height int) float64 {
	return float64(width) / float64(height)
}
