package controllers

import (
	"fmt"
	"log"
	"net/http"
	"otus_go_final/internal/services"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

func Index(w http.ResponseWriter, r *http.Request) {
	var err error

	width := chi.URLParam(r, "width")
	height := chi.URLParam(r, "height")

	fmt.Println(width, height)

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

	resized, err := service.Invoke()
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(resized)
}
