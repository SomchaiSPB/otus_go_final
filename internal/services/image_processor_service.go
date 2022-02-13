package services

import (
	"bytes"
	"context"
	"errors"
	"image"
	"io"
	"log"
	"net/http"
	"net/url"

	_ "image/jpeg" // Side effect import.
)

var (
	ErrBrokenURL          = errors.New("url is broken or invalid")
	ErrNotFound           = errors.New("target URL not found")
	ErrInternalServer     = errors.New("target server internal error")
	ErrNotImageType       = errors.New("target file is not an image")
	ErrImageSizeViolation = errors.New("target image size is less than wanted")
	ErrServerNotExists    = errors.New("target server does not exist")
)

type ImageProperty struct {
	width     int
	height    int
	targetURL string
	header    http.Header
}

func NewImageProperty(width int, height int, target string, h http.Header) *ImageProperty {
	return &ImageProperty{
		width:     width,
		height:    height,
		targetURL: target,
		header:    h,
	}
}

type ImageProcessService struct {
	Resizer        Resizer
	InputProps     *ImageProperty
	OutputHeaders  string
	OutputImage    []byte
	OriginalHeader http.Header
	Client         *http.Client
	ResponseCode   int
}

func NewProcessService(props *ImageProperty, r *JpegResizer) *ImageProcessService {
	return &ImageProcessService{
		Resizer:    r,
		InputProps: props,
		Client:     &http.Client{},
	}
}

func (s *ImageProperty) GetWidth() int {
	return s.width
}

func (s *ImageProperty) GetHeight() int {
	return s.height
}

func (s *ImageProcessService) Invoke() ([]byte, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}

	img, err := s.ProxyRequest()
	if err != nil {
		return nil, err
	}

	im, _, err := image.DecodeConfig(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}

	if im.Width <= s.InputProps.width || im.Height <= s.InputProps.height {
		s.ResponseCode = http.StatusBadRequest
		return nil, ErrImageSizeViolation
	}

	m, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return nil, err
	}

	result, err := s.Resizer.Resize(m)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ImageProcessService) Validate() error {
	validURL, err := url.ParseRequestURI(s.InputProps.targetURL)
	if err != nil {
		s.ResponseCode = http.StatusInternalServerError
		return err
	}

	s.InputProps.targetURL = validURL.String()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, s.InputProps.targetURL, nil)
	if err != nil {
		s.ResponseCode = http.StatusInternalServerError
		return ErrServerNotExists
	}

	res, err := s.Client.Do(req)
	if err != nil {
		s.ResponseCode = http.StatusInternalServerError
		return ErrServerNotExists
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	s.ResponseCode = res.StatusCode

	switch c := res.StatusCode; {
	case c >= 400 && c <= 499:
		return ErrNotFound
	case c >= 500 && c <= 599:
		return ErrInternalServer
	default:
		return nil
	}
}

func (s *ImageProcessService) ProxyRequest() ([]byte, error) {
	var err error

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, s.InputProps.targetURL, nil)
	if err != nil {
		log.Println(err)
		s.ResponseCode = http.StatusBadRequest
		return nil, err
	}

	for key, val := range s.OriginalHeader {
		req.Header.Set(key, val[0])
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		log.Println(err)
		s.ResponseCode = resp.StatusCode
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrBrokenURL
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		s.ResponseCode = http.StatusBadRequest
		return nil, err
	}

	fileType := http.DetectContentType(data)

	switch fileType {
	case "image/jpeg":
		// case "image/png":
		return data, nil
	default:
		s.ResponseCode = http.StatusBadRequest
	}

	return nil, ErrNotImageType
}
