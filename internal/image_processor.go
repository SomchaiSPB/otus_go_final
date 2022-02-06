package internal

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/nfnt/resize"
)

type ImageProcessor struct {
	format string
	image  image.Image
	width  int
	height int
}

func NewImageProcessor(format string, image image.Image, width int, height int) *ImageProcessor {
	return &ImageProcessor{
		format: format,
		image:  image,
		width:  width,
		height: height,
	}
}

func (p *ImageProcessor) Resize() ([]byte, error) {
	buf := new(bytes.Buffer)

	newImage := resize.Resize(uint(p.width), uint(p.height), p.image, resize.Lanczos3)

	err := jpeg.Encode(buf, newImage, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
