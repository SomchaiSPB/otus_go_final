package services

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
)

var (
	ErrBrokenUrl = errors.New("url is broken or invalid")
)

type ImageProperty struct {
	width     int
	height    int
	targetUrl string
}

func NewImageProperty(width int, height int, target string) *ImageProperty {
	return &ImageProperty{
		width:     width,
		height:    height,
		targetUrl: target,
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

func (s *ImageProcessService) Invoke() error {
	err := s.Validate()

	if err != nil {
		return err
	}

	response, err := s.ProxyRequest()

	if err != nil {
		log.Println(err)
		return err
	}

	var buf bytes.Buffer

	n, err := response.Read(buf.Bytes())

	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(buf)

	return nil
}

func (s *ImageProcessService) Validate() error {
	schema := "https://"
	fullUrl := schema + s.InputProps.targetUrl

	validUrl, err := url.ParseRequestURI(fullUrl)

	if err != nil {
		log.Println(err)
		return err
	}

	s.InputProps.targetUrl = validUrl.String()

	return nil
}

func (s *ImageProcessService) ProxyRequest() (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, s.InputProps.targetUrl, nil)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for key, val := range s.OriginalHeader {
		req.Header.Set(key, val[0])
	}

	res, err := s.Client.Do(req)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, ErrBrokenUrl
	}

	return res.Body, nil
}
