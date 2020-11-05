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
	"sync"

	"github.com/anthonynsimon/bild/transform"
)

type Image struct {
	image.Image
	mx sync.Mutex
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
	m := NewImage(i)
	if err = m.convert(width, height); err != nil {
		return nil, err
	}
	return encode(m.Image)
}

func NewImage(img image.Image) Image {
	return Image{img, sync.Mutex{}}
}

func (img *Image) convert(width int, height int) error {
	widthOrig := img.Bounds().Max.X
	heightOrig := img.Bounds().Max.Y
	if width <= 0 || height <= 0 {
		return errors.New("can't reduce toOrBelow zero")
	}
	sfOriginal, _ := sizeFactor(widthOrig, heightOrig)
	sfNew, _ := sizeFactor(width, height)

	switch {
	case sfOriginal > sfNew:
		calcWidth := int(float64(height) * sfOriginal)
		if err := img.resize(calcWidth, height); err != nil {
			return err
		}
		if err := img.crop(image.Point{X: (calcWidth - width) / 2, Y: 0}, image.Point{X: (calcWidth-width)/2 + width, Y: height}); err != nil {
			return err
		}
	case sfOriginal == sfNew:
		if err := img.resize(width, height); err != nil {
			return err
		}
	case sfOriginal < sfNew:
		calcHeight := int(float64(width) / sfOriginal)
		if err := img.resize(width, calcHeight); err != nil {
			return err
		}
		if err := img.crop(image.Point{X: 0, Y: (calcHeight - height) / 2}, image.Point{X: width, Y: (calcHeight-height)/2 + height}); err != nil {
			return err
		}
	}
	return nil
}

func (img *Image) resize(width, height int) error {
	img.mx.Lock()
	defer img.mx.Unlock()
	if width <= 0 || height <= 0 {
		return errors.New("can't resize to zero or negative value")
	}
	tmpImg := transform.Resize(img, width, height, transform.Linear)
	img.Image = tmpImg
	return nil
}

func (img *Image) crop(p1 image.Point, p2 image.Point) error {
	img.mx.Lock()
	defer img.mx.Unlock()
	if img.Image == nil {
		return errors.New("corrupted image")
	}
	if p1.X < 0 || p1.Y < 0 || p2.X > img.Image.Bounds().Max.X || p2.Y > img.Image.Bounds().Max.Y {
		return errors.New("not valid corner points")
	}
	b := image.Rect(0, 0, p2.X-p1.X, p2.Y-p1.Y)
	resImg := image.NewRGBA(b)
	draw.Draw(resImg, b, img.Image, p1, draw.Src)
	img.Image = resImg
	return nil
}

func sizeFactor(width int, height int) (float64, error) {
	if height == 0 {
		return 0, errors.New("can't divide by zero")
	}
	return float64(width) / float64(height), nil
}
