package main

import (
	"animeGAN2go/Job"
	"animeGAN2go/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	// запускаем горутину

	go Job.StartSingleWorker()

	r := mux.NewRouter()

	// запускаем в фоне воркер который в фоне опрашивает очередь на предмет обновлений

	r.HandleFunc("/start", handlers.StarJobAdd)

	r.HandleFunc("/delete", handlers.DeleteJob)

	http.Handle("/", r)

	http.ListenAndServe(":9000", nil)
}
