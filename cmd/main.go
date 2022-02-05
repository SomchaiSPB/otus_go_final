package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"otus_go_final/internal/controllers"
)

var (
	port string
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Recoverer)

	r.Get("/check", controllers.Check)
	r.Get("/fill/{width}/{height}/{target}*", controllers.Index)

	r.NotFoundHandler()

	err := http.ListenAndServe(":"+port, r)

	if err != nil {
		log.Println("server error " + err.Error())
	}
}

func init() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("error loading .env file. Using default values")
		port = "4000"
	} else {
		port = os.Getenv("SERVER_PORT")
	}
}
