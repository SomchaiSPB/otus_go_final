package controllers

import (
	"log"
	"net/http"
)

func Check(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))

	if err != nil {
		log.Println("error writer " + err.Error())
	}
}
