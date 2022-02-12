package controllers

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"otus_go_final/config"
	imagecache "otus_go_final/internal/cache"
	"otus_go_final/internal/services"
)

const (
	HTTP  = "http://"
	HTTPS = "https://"
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
	target = HTTP + target

	widthInt, err := strconv.Atoi(width)
	if err != nil {
		log.Println("error converting width " + err.Error())
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	heightInt, err := strconv.Atoi(height)
	if err != nil {
		log.Println("error converting height " + err.Error())
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))

		return
	}

	key := asSha256(widthInt, heightInt, target)

	result, hit := h.cache.Get(imagecache.Key(key))

	if hit {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write(result.([]byte))

		return
	}

	imageProp := services.NewImageProperty(widthInt, heightInt, target)

	service := services.NewProcessService(imageProp, r.Header)

	resized, err := service.Invoke()
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(service.ResponseCode)
		w.Write([]byte(err.Error()))

		return
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
