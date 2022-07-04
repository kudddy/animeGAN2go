package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const urlPathToSkill = "https://smartapp-code.sberdevices.ru/chatadapter/chatapi/webhook/sber_nlp2/cGnGPZWb:45c9c4e54edfcf2cfe505f84e3f338185a334e42"

func policyTlgSm(update UpdateType) error {

	session, _ := CacheSystem.Get(string(rune(update.Message.User.Id)))

	// convert message to sm
	reqToSm := generatePayloadForSm(update.Message.Text, session.sessionId, session.messageId)

	// send message to sm and get resp
	resp, err := sendReqToSm(urlPathToSkill, reqToSm)
	if err != nil {
		log.Printf("Someting wrong with request to SM with mid - %d", session.messageId)

	}
	// convert message to tlg format
	var textToUser string

	if resp.MessageName == "ANSWER_TO_USER" {
		textToUser = resp.Payload.PronounceText
	} else if resp.MessageName == "NOTHING_FOUND" {

		// in this place we should get from db user_id with max score

		var operatorBotId = 81432612

		// and create new session data for bot

		//session, _ := CacheSystem.Get(string(rune(update.Message.User.Id)))

		// create session id for bot
		CacheSystem.Put(string(rune(operatorBotId)), sessionData{
			messageId:       0,
			sessionId:       session.sessionId,
			botStatus:       session.botStatus,
			companionUserId: update.Message.User.Id,
		})

		// update user session param as companionUserId
		CacheSystem.Put(string(rune(update.Message.User.Id)), sessionData{
			messageId:       session.messageId,
			sessionId:       session.sessionId,
			botStatus:       session.botStatus,
			companionUserId: operatorBotId,
		})

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

func policyOperatorBot(update UpdateType, path string) error {

	if path == "/operator" {

		session, _ := CacheSystem.Get(string(rune(update.Message.User.Id)))

		reqToTlg := OutMessage{
			Text:   update.Message.Text,
			ChatId: session.companionUserId,
		}

		// send req to tlg
		err := sendReqToTlg(BuildUrl(PathSendMessage, BotsInfo["bot"]), reqToTlg)

		if err != nil {
			log.Printf("Someting wrong with request to tlg")
			log.Print(err)
			return err
		}

	} else {

		var buttons []InlineKeyboardButton

		buttons = append(buttons, InlineKeyboardButton{
			"Завершить чат",
			nil,
			nil,
			"close_chat",
			nil,
			nil,
			nil,
			nil,
		},
		)

		var arrayOfByttons [][]InlineKeyboardButton

		arrayOfByttons = append(arrayOfByttons, buttons)

		var inlineButtons = InlineKeyboardMarkup{
			InlineKeyboard: arrayOfByttons,
		}

		reqToTlg := OutMessage{
			Text:        update.Message.Text,
			ChatId:      update.Message.Chat.Id,
			ReplyMarkup: &inlineButtons,
		}

		// send req to tlg
		err := sendReqToTlg(BuildUrl(PathSendMessage, BotsInfo["operator"]), reqToTlg)
		if err != nil {
			log.Printf("Someting wrong with request to tlg")
			log.Print(err)
			return err
		}

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
			if r.URL.Path != "/operator" {
				_ = policyTlgSm(update)
			}

		} else {
			_ = policyOperatorBot(update, r.URL.Path)
		}
	} else {
		if r.URL.Path != "/operator" {

			// create new session data
			session := "bot-" + time.Now().Format("20060102150405")
			CacheSystem.Put(string(rune(update.Message.User.Id)), sessionData{
				messageId: 0,
				sessionId: session,
				botStatus: true,
			})
			_ = policyTlgSm(update)

		} else {
			reqToTlg := OutMessage{
				Text:   "Активный диалогов нет:) Отдыхайте!😍",
				ChatId: update.Message.Chat.Id,
			}
			// send req to tlg
			err = sendReqToTlg(BuildUrl(PathSendMessage, BotsInfo["bot"]), reqToTlg)
		}
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
