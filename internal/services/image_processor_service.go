package services

import (
	"bytes"
	"context"
	"errors"
	"image"
	_ "image/jpeg" //nolint
	"io"
	"log"
	"net/http"
	"net/url"

	"otus_go_final/internal"
)

const (
	HTTP  = "http"
	HTTPS = "https"
)

var (
	ErrBrokenURL          = errors.New("url is broken or invalid")
	ErrRedirectResponse   = errors.New("target URL redirect status")
	ErrNotFound           = errors.New("target URL not found")
	ErrInternalServer     = errors.New("target server internal error")
	ErrNotImageType       = errors.New("target file is not an image")
	ErrImageSizeViolation = errors.New("target image size is less than wanted")
)

type ImageProperty struct {
	width     int
	height    int
	targetURL string
}

func NewImageProperty(width int, height int, target string) *ImageProperty {
	return &ImageProperty{
		width:     width,
		height:    height,
		targetURL: target,
	}
}

type ImageProcessService struct {
	InputProps     *ImageProperty
	OutputHeaders  string
	OutputImage    []byte
	OriginalHeader http.Header
	Client         *http.Client
	ResponseCode   int
}

func NewProcessService(props *ImageProperty, headers http.Header) *ImageProcessService {
	return &ImageProcessService{
		InputProps:     props,
		OriginalHeader: headers,
		Client:         &http.Client{},
	}
}

func (s *ImageProcessService) Invoke() ([]byte, error) {
	code, err := s.Validate()
	s.ResponseCode = code

	if err != nil {
		log.Println(code, err)
		return nil, err
	}

	img, err := s.ProxyRequest()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	im, _, err := image.DecodeConfig(bytes.NewReader(img))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if im.Width <= s.InputProps.width || im.Height <= s.InputProps.height {
		return nil, ErrImageSizeViolation
	}

	m, format, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	processor := internal.NewImageProcessor(format, m, s.InputProps.width, s.InputProps.height)

	result, err := processor.Resize()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func (s *ImageProcessService) Validate() (int, error) {
	validURL, err := url.Parse(s.InputProps.targetURL)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}

	validURL.Scheme = HTTP

	u, err := url.ParseRequestURI(validURL.String())
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, err
	}

	s.InputProps.targetURL = u.String()

	res, err := s.Client.Get(s.InputProps.targetURL)

	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Println(err)
			return
		}
	}()

	return res.StatusCode, err
}

func (s *ImageProcessService) ProxyRequest() ([]byte, error) {
	var err error

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, s.InputProps.targetURL, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for key, val := range s.OriginalHeader {
		req.Header.Set(key, val[0])
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		log.Println(err)
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
		return nil, err
	}

	fileType := http.DetectContentType(data)

	if fileType != "image/jpeg" {
		return nil, ErrNotImageType
	}

	return data, nil
}
