package services

import (
	"bytes"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
)

type ImageProcessor struct {
	format string
	image  image.Image
	props  *ImageProperty
}

type Resizer interface {
	Resize() ([]byte, error)
}

func NewImageProcessor(format string, image image.Image, p *ImageProperty) *ImageProcessor {
	return &ImageProcessor{
		format: format,
		image:  image,
		props:  p,
	}
}

func (p *ImageProcessor) Resize() ([]byte, error) {
	buf := new(bytes.Buffer)

	newImage := resize.Resize(uint(p.props.GetWidth()), uint(p.props.GetHeight()), p.image, resize.Lanczos3)

	err := jpeg.Encode(buf, newImage, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
