package main

import (
	"dialog-policy/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// запускаем горутину

	r := mux.NewRouter()

	http.Handle("/", r)
	r.HandleFunc("/{authToken}/bot", handlers.Handler)
	r.HandleFunc("/{authToken}/operator", handlers.Handler)
	r.HandleFunc("/{authToken}/update", handlers.Handler)

	http.ListenAndServe(":9001", nil)
}
