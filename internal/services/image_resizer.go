package services

import (
	"bytes"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
)

type JpegResizer struct {
	props *ImageProperty
}

type Resizer interface {
	Resize(image image.Image) ([]byte, error)
}

func NewImageResizer(p *ImageProperty) *JpegResizer {
	return &JpegResizer{
		props: p,
	}
}

func (p *JpegResizer) Resize(image image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)

	newImage := resize.Resize(uint(p.props.GetWidth()), uint(p.props.GetHeight()), image, resize.Lanczos3)

	err := jpeg.Encode(buf, newImage, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
