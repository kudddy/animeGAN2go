package handlers

import (
	"encoding/json"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// нужно обернуть для получения данных от основной горутины
	// проверяем валидный ли токен
	decoder := json.NewDecoder(r.Body)
	var t ReqData
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	var workerStatus CheckTokenResp

	workerStatus.MessageName = "STARTJOBADD"

	workerStatus.Desc = "OK, start working"
	go StartWorker(t)

	js, err := json.Marshal(workerStatus)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
