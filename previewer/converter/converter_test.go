package converter

import (
	"github.com/stretchr/testify/require"
	"image"
	"image/color"
	"testing"
)

func TestResize(t *testing.T) {
	table := []struct {
		width     int
		height    int
		expectedX int
		expectedY int
		err       bool
		msg       string
	}{
		{
			width: 300, height: 200, expectedX: 300, expectedY: 200, err: false, msg: "Reducing the image size",
		},
		{
			width: 1600, height: 1200, expectedX: 1600, expectedY: 1200, err: false, msg: "Increasing the image size",
		},
		{
			width: 0, height: 0, expectedX: 800, expectedY: 600, err: true, msg: "Resize to zero",
		},
		{
			width: -1000, height: -1200, expectedX: 800, expectedY: 600, err: true, msg: "Use negative values",
		},
	}

	for _, dat := range table {
		t.Run(dat.msg, func(t *testing.T) {
			img := Image{Image: createImage(800, 600)}
			i := false
			err := img.resize(dat.width, dat.height)
			if err != nil {
				i = true
			}
			require.Equal(t, dat.err, i, dat.msg)
			require.Equal(t, dat.expectedX, img.Bounds().Max.X, dat.msg)
			require.Equal(t, dat.expectedY, img.Bounds().Max.Y, dat.msg)
		})
	}
}

func TestCrop(t *testing.T) {
	table := []struct {
		topLeft     image.Point
		bottomRight image.Point
		expectedX   int
		expectedY   int
		err         bool
		msg         string
	}{
		{
			topLeft: image.Point{X: 400, Y: 0}, bottomRight: image.Point{X: 600, Y: 1000}, expectedX: 200, expectedY: 1000, err: false, msg: "Vertical crop",
		},
		{
			topLeft: image.Point{X: 0, Y: 400}, bottomRight: image.Point{X: 1000, Y: 600}, expectedX: 1000, expectedY: 200, err: false, msg: "Horizontal crop",
		},
		{
			topLeft: image.Point{X: 100, Y: 100}, bottomRight: image.Point{X: 900, Y: 900}, expectedX: 800, expectedY: 800, err: false, msg: "Square crop",
		},
		{
			topLeft: image.Point{X: -100, Y: 0}, bottomRight: image.Point{X: 2000, Y: 1000}, expectedX: 1000, expectedY: 1000, err: true, msg: "Too wide crop with negative offset",
		},
		{
			topLeft: image.Point{X: 0, Y: -100}, bottomRight: image.Point{X: 1000, Y: 2000}, expectedX: 1000, expectedY: 1000, err: true, msg: "Too tall crop with negative offset",
		},
		{
			topLeft: image.Point{X: 100, Y: 0}, bottomRight: image.Point{X: 2000, Y: 1000}, expectedX: 1000, expectedY: 1000, err: true, msg: "Too wide crop with positive offset",
		},
		{
			topLeft: image.Point{X: 0, Y: 100}, bottomRight: image.Point{X: 1000, Y: 2000}, expectedX: 1000, expectedY: 1000, err: true, msg: "Too tall crop with positive offset",
		},
	}

	for _, dat := range table {
		t.Run(dat.msg, func(t *testing.T) {
			img := Image{Image: createImage(1000, 1000)}
			i := false
			err := img.crop(dat.topLeft, dat.bottomRight)
			if err != nil {
				i = true
			}
			require.Equal(t, dat.err, i, dat.msg)
			require.Equal(t, dat.expectedX, img.Image.Bounds().Max.X, dat.msg)
			require.Equal(t, dat.expectedY, img.Image.Bounds().Max.Y, dat.msg)
		})
	}
}

func TestConvert(t *testing.T) {
	originalAspect := 800.0 / 600.0
	releasedValue := 3000
	table := []struct {
		width     int
		height    int
		expectedX int
		expectedY int
		err       bool
		msg       string
	}{
		{
			width: 400, height: 600, expectedX: 400, expectedY: 600, err: false, msg: "Reducing the image size by horizontal",
		},
		{
			width: 800, height: 400, expectedX: 800, expectedY: 400, err: false, msg: "Reducing the image size by vertical",
		},
		{
			width: 400, height: int(400 / originalAspect), expectedX: 400, expectedY: int(400 / originalAspect), err: false, msg: "Resize to original aspect ratio",
		},
		{
			width: 1000, height: releasedValue, expectedX: 1000, expectedY: releasedValue, err: false, msg: "Increasing the image size by horizontal",
		},
		{
			width: releasedValue, height: 1000, expectedX: releasedValue, expectedY: 1000, err: false, msg: "Increasing the image size by vertical",
		},
		{
			width: 0, height: 0, expectedX: 800, expectedY: 600, err: true, msg: "Resize to zero",
		},
		{
			width: -1000, height: -1200, expectedX: 800, expectedY: 600, err: true, msg: "Use negative values",
		},
	}

	for _, dat := range table {
		t.Run(dat.msg, func(t *testing.T) {
			img := Image{Image: createImage(800, 600)}
			i := false
			err := img.convert(dat.width, dat.height)
			if err != nil {
				i = true
			}
			require.Equal(t, dat.err, i, dat.msg)
			require.Equal(t, dat.expectedX, img.Image.Bounds().Max.X, dat.msg)
			require.Equal(t, dat.expectedY, img.Image.Bounds().Max.Y, dat.msg)
		})
	}
}

func createImage(w, h int) image.Image {
	res := image.NewRGBA(image.Rectangle{Min: image.Point{X: 0, Y: 0}, Max: image.Point{X: w, Y: h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			switch {
			case x%10 == 0 || y%10 == 00: // upper left quadrant
				res.Set(x, y, color.Black)
			default:
				res.Set(x, y, color.White)
			}
		}
	}
	return res
}
