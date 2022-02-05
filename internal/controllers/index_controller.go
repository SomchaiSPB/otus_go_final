package controllers

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"otus_go_final/internal/errors"
	"otus_go_final/internal/services"
	"strconv"
	"strings"
)

var (
	ErrNotFound = &errors.ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}
)

func Index(w http.ResponseWriter, r *http.Request) {
	var err error

	width := chi.URLParam(r, "width")
	height := chi.URLParam(r, "height")

	url := strings.Split(r.URL.String(), "/")
	target := strings.Join(url[4:], "/")

	widthInt, err := strconv.Atoi(width)
	if err != nil {
		log.Println("error converting width " + err.Error())
		w.Write([]byte(err.Error()))

		return
	}

	heightInt, err := strconv.Atoi(height)
	if err != nil {
		log.Println("error converting height " + err.Error())
		w.Write([]byte(err.Error()))

		return
	}

	image := services.NewImageProperty(widthInt, heightInt, target)

	service := services.NewProcessService(image, r.Header)

	err = service.Invoke()

	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}

	w.Write([]byte("request received"))
}
