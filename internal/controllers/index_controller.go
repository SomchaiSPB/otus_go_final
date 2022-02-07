package controllers

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"otus_go_final/config"
	"otus_go_final/internal/services"
	"strconv"
	"strings"

	"github.com/go-chi/chi"

	imagecache "otus_go_final/internal/cache"
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

	key := asSha256(widthInt, heightInt, target)

	result, ok := h.cache.Get(imagecache.Key(key))

	if ok {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(result.([]byte))
		return
	}

	imageProp := services.NewImageProperty(widthInt, heightInt, target)

	service := services.NewProcessService(imageProp, r.Header)

	resized, err := service.Invoke()
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}

	h.cache.Set(imagecache.Key(key), resized)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(resized)
}

func asSha256(width, height int, url string) string {
	h := sha256.New()

	h.Write([]byte(fmt.Sprintf("%v", width)))
	h.Write([]byte(fmt.Sprintf("%v", height)))
	h.Write([]byte(fmt.Sprintf("%v", url)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
