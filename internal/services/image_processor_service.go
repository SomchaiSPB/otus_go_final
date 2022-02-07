package services

import (
	"bytes"
	"context"
	"errors"
	_ "image/jpeg" // Register jpeg package.
	"io"
	"log"
	"net/http"
	"net/url"

	"image"
	"otus_go_final/internal"
)

var ErrBrokenURL = errors.New("url is broken or invalid")

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
}

func NewProcessService(props *ImageProperty, headers http.Header) *ImageProcessService {
	return &ImageProcessService{
		InputProps:     props,
		OriginalHeader: headers,
		Client:         &http.Client{},
	}
}

func (s *ImageProcessService) Invoke() ([]byte, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}

	img, err := s.ProxyRequest()
	if err != nil {
		log.Println(err)
		return nil, err
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

func (s *ImageProcessService) Validate() error {
	schema := "http://"
	fullURL := schema + s.InputProps.targetURL

	validURL, err := url.ParseRequestURI(fullURL)
	if err != nil {
		log.Println(err)
		return err
	}

	s.InputProps.targetURL = validURL.String()

	return nil
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
	}

	return data, nil
}
