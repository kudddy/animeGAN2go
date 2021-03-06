package handlers

import (
	"animeGAN2go/Job"
	"animeGAN2go/MessageTypes"
	"encoding/json"
	"net/http"
)

func StarJobAdd(res http.ResponseWriter, req *http.Request) {
	// нужно обернуть для получения данных от основной горутины
	// проверяем валидный ли токен
	decoder := json.NewDecoder(req.Body)
	var t MessageTypes.ReqData
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	var workerStatus MessageTypes.CheckTokenResp

	workerStatus.MessageName = "STARTJOBADD"

	workerStatus.Desc = "OK, start working"
	go Job.StartWorker(t)

	js, err := json.Marshal(workerStatus)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(js)

}
