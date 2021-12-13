package handlers

import (
	"animeGAN2go/MessageTypes"
	"animeGAN2go/plugins/pg"
	"encoding/json"
	"net/http"
)

func DeleteJob(res http.ResponseWriter, req *http.Request) {
	// нужно обернуть для получения данных от основной горутины
	// проверяем валидный ли токен
	decoder := json.NewDecoder(req.Body)
	var t MessageTypes.ReqData
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}

	var workerStatus MessageTypes.CheckTokenResp

	workerStatus.MessageName = "DELETE_JOB"

	//запуск воркера
	//fmt.Println("запуск воркера")
	// нужна проверка не запущен ли воркер уже/узнать статус и только потом запускать
	// пишем в memcached

	pg.InsertCancelAction(int(t.ChatId))

	workerStatus.Desc = "OK, start working"

	js, err := json.Marshal(workerStatus)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(js)

}
