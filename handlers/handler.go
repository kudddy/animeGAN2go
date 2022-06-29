package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const urlPathToSkill = "https://smartapp-code.sberdevices.ru/chatadapter/chatapi/webhook/sber_nlp2/cGnGPZWb:45c9c4e54edfcf2cfe505f84e3f338185a334e42"

func policyTlgSm(update UpdateType, sessionId string, messageId int) error {

	// convert message to sm
	reqToSm := generatePayloadForSm(update.Message.Text, sessionId, messageId)

	// send message to sm and get resp
	resp, err := sendReqToSm(urlPathToSkill, reqToSm)
	if err != nil {
		log.Printf("Someting wrong with request to SM with mid - %d", messageId)

	}
	// convert message to tlg format
	var textToUser string

	if resp.MessageName == "ANSWER_TO_USER" {
		textToUser = resp.Payload.PronounceText
	} else if resp.MessageName == "NOTHING_FOUND" {

		CacheSystem.ChangeBotStatus(string(rune(update.Message.User.Id)))

		textToUser = "Переадресую на оператора"
	}

	reqToTlg := OutMessage{
		Text:   textToUser,
		ChatId: update.Message.Chat.Id,
	}

	// send req to tlg

	err = sendReqToTlg(BuildUrl(PathSendMessage, BotsInfo["bot"]), reqToTlg)

	if err != nil {
		log.Printf("Someting wrong with request to tlg")
		log.Print(err)
		return err
	}

	return nil

}

func policyOperatorBot(update UpdateType) error {
	reqToTlg := OutMessage{
		Text:   update.Message.Text,
		ChatId: update.Message.Chat.Id,
	}
	// send req to tlg
	err := sendReqToTlg(BuildUrl(PathSendMessage, BotsInfo["operator"]), reqToTlg)
	if err != nil {
		log.Printf("Someting wrong with request to tlg")
		log.Print(err)
		return err
	}
	return nil
}

// Метод Handler. Данный метод будет обрабатывать HTTP запросы поступающие к функции
func Handler(w http.ResponseWriter, r *http.Request) {

	// get message from tlg
	if r.Body != nil {
		defer r.Body.Close()
	}

	decoder := json.NewDecoder(r.Body)

	var update UpdateType
	err := decoder.Decode(&update)
	if err != nil {
		panic(err)
	}

	// Логирование входящего запроса
	log.Printf("Request received: %s\nMethod: %s\nPATH: %s\nRAW_PATH: %s\nRAW_QUERY:%s", update.Message.Text, r.Method, r.URL.Path, r.URL.RawPath, r.URL.RawQuery)

	// check cache
	cache, check := CacheSystem.Get(string(rune(update.Message.User.Id)))
	// проверяем кэш
	if check {
		//если есть, то смотрим что лежит внутри
		// достаем сессию
		//достаем флаг редиректа на оператора, если флаг == true, отсылаем запрос в бот оператора
		if cache.botStatus {
			_ = policyTlgSm(update, cache.sessionId, cache.messageId)
		} else {
			_ = policyOperatorBot(update)
		}
	} else {
		// создаем новые сессионные данные
		session := "bot-" + time.Now().Format("20060102150405")
		CacheSystem.Put(string(rune(update.Message.User.Id)), sessionData{
			messageId: 0,
			sessionId: session,
			botStatus: true,
		})
		_ = policyTlgSm(update, session, 0)
	}

	var workerStatus RespByServ

	workerStatus.Ok = true

	js, err := json.Marshal(workerStatus)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	// get message from tlg
}
