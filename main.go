package main

import (
	"dialog-policy/handlers"
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	// запускаем горутину

	flag.Parse()

	r := mux.NewRouter()

	fmt.Println("app start")

	http.Handle("/", r)
	fmt.Println("start handle /")
	r.HandleFunc("/{authToken}/bot", handlers.Handler)
	fmt.Println("start handle /bot")
	r.HandleFunc("/{authToken}/operator", handlers.Handler)
	fmt.Println("start handle /operator")
	r.HandleFunc("/{authToken}/update", handlers.Handler)

	http.ListenAndServe(":9001", nil)
}
