package main

import (
	"dialog-policy/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	// запускаем горутину

	r := mux.NewRouter()

	http.Handle("/", r)
	r.HandleFunc("/start", handlers.Handler)

	http.ListenAndServe(":9000", nil)
}
