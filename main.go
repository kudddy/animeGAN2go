package main

import (
	"animeGAN2go/handlers"
	"github.com/gorilla/mux"
	"net/http"
)


func main() {


	//memcacheClient := memcache.New("127.0.0.1:11211")
	r := mux.NewRouter()

	r.HandleFunc("/start", handlers.StarJobAdd)

	http.Handle("/", r)
	http.ListenAndServe(":9000", nil)
}