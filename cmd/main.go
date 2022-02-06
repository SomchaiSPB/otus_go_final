package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"otus_go_final/config"
	"otus_go_final/internal/controllers"
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

	err := http.ListenAndServe(":"+port, r)
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
	} else {
		err := envconfig.Process("previewer", &cfg)
		if err != nil {
			panic("error load config file")
		}
	}
}
