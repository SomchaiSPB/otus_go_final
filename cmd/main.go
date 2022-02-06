package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"otus_go_final/config"
	"otus_go_final/internal/controllers"
	"strconv"
)

var port string
var cfg *config.Config

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	handler := controllers.NewBaseHandler(cfg)

	r.Get("/check", controllers.Check)
	r.Get("/fill/{width}/{height}/{target}*", handler.Index)

	r.NotFoundHandler()

	err := http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Println("server error " + err.Error())
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error loading .env file. Using default values")
		cfg = &config.Config{
			Port:     "4000",
			Capacity: 100,
		}
		return
	}

	capacity, err := strconv.Atoi(os.Getenv("PREVIEWER_CAPACITY"))

	if err != nil {
		panic("error converting port to int " + err.Error())
	}

	cfg = &config.Config{
		Port:     os.Getenv("PREVIEWER_PORT"),
		Capacity: capacity,
	}
}
