package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	// запускаем горутину

	r := mux.NewRouter()

	// запускаем в фоне воркер который в фоне опрашивает очередь на предмет обновлений

	http.Handle("/", r)

	http.ListenAndServe(":9000", nil)
}
