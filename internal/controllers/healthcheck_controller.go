package controllers

import (
	"log"
	"net/http"
)

func Check(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Println("error writer " + err.Error())
	}
}
