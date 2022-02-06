package controllers

import (
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"otus_go_final/config"
	imagecache "otus_go_final/internal/cache"
	"otus_go_final/internal/services"
	"strconv"
	"strings"
)

type BaseHandler struct {
	cache imagecache.Cache
	cfg   *config.Config
}

func NewBaseHandler(cfg *config.Config) *BaseHandler {
	return &BaseHandler{
		cfg:   cfg,
		cache: imagecache.NewCache(cfg.Capacity),
	}
}

func (h *BaseHandler) Index(w http.ResponseWriter, r *http.Request) {
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

	imageProp := services.NewImageProperty(widthInt, heightInt, target)

	service := services.NewProcessService(imageProp, r.Header)

	resized, err := service.Invoke()
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(resized)
}
