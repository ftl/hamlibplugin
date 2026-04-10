package graphic

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	"github.com/vincent-petithory/dataurl"
)

const (
	imageSize = 64
)

var (
	Red  = color.RGBA{255, 0, 0, 255}
	Blue = color.RGBA{0, 255, 0, 255}
)

func GenerateSimpleImageURL(clr color.RGBA) (string, error) {
	img := generateSimpleImage(clr)
	pngBytes, err := imageToPNG(img)
	if err != nil {
		return "", err
	}
	return pngToDataURL(pngBytes)
}

func generateSimpleImage(clr color.RGBA) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))
	draw.Draw(img, img.Bounds(), image.NewUniform(clr), image.Point{}, draw.Src)
	return img
}

func imageToPNG(img image.Image) ([]byte, error) {
	pngBytes := make([]byte, 0, 1024)
	buffer := bytes.NewBuffer(pngBytes)
	err := png.Encode(buffer, img)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func pngToDataURL(pngBytes []byte) (string, error) {
	dataURL := dataurl.New(pngBytes, "image/png")
	urlBytes, err := dataURL.MarshalText()
	if err != nil {
		return "", err
	}
	return string(urlBytes), nil
}
