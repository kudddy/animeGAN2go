package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

const urlPathToSkill = "https://smartapp-code.sberdevices.ru/chatadapter/chatapi/webhook/sber_nlp2/cGnGPZWb:45c9c4e54edfcf2cfe505f84e3f338185a334e42"

func policyTlgSm(update UpdateType) error {

	session, _ := CacheSystem.Get(update.Message.User.Id)

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
		operators := CacheSystem.GetRandomAuthOperators()
		if len(operators) > 0 {
			// first in, first out
			var operatorBotId = operators[0]

			operatorSession, _ := CacheSystem.Get(operatorBotId)

			// and create new session data for bot

			//session, _ := CacheSystem.Get(string(rune(update.Message.User.Id)))

			log.Printf("save session parametrs for operator when id is %d, companion id is %d", operatorBotId, update.Message.User.Id)
			// TODO its not true
			CacheSystem.ChangeBotStatus(update.Message.User.Id)

			// create session id for bot
			CacheSystem.Put(operatorBotId, sessionData{
				messageId:       0,
				sessionId:       operatorSession.sessionId,
				botStatus:       false,
				companionUserId: update.Message.User.Id,
				auth:            operatorSession.auth,
				busy:            true,
			})

			log.Printf("save session parametrs for bot when id is %d, companion id is %d", update.Message.User.Id, operatorBotId)
			// update user session param as companionUserId
			CacheSystem.Put(update.Message.User.Id, sessionData{
				messageId:       session.messageId,
				sessionId:       session.sessionId,
				botStatus:       false,
				companionUserId: operatorBotId,
			})

			s, _ := CacheSystem.Get(update.Message.User.Id)

			log.Printf("session parameters from cache for %d is %d", update.Message.User.Id, s.companionUserId)

			d, _ := CacheSystem.Get(operatorBotId)

			log.Printf("session parameters from cache for %d is %d", operatorBotId, d.companionUserId)

			textToUser = "Переадресую на оператора❤️"

		} else {
			textToUser = "Сейчас все операторы заняты. Пока только бот🏃‍♀️"
		}

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

	session, _ := CacheSystem.Get(update.Message.User.Id)

	if path == "/operator" {

		if update.Message.Text == "Завершить чат" {

			reqToTlg := OutMessage{
				Text:   "Чат с оператором завершен!😎",
				ChatId: session.companionUserId,
			}

			// send req to tlg
			err := sendReqToTlg(BuildUrl(PathSendMessage, BotsInfo["bot"]), reqToTlg)

			if err != nil {
				log.Printf("Someting wrong with request to tlg")
				log.Print(err)
				return err
			}

			reqToTlg = OutMessage{
				Text:   "Сессия с клиентом завершена!😎",
				ChatId: update.Message.Chat.Id,
			}

			// send req to tlg
			err = sendReqToTlg(BuildUrl(PathSendMessage, BotsInfo["operator"]), reqToTlg)

			if err != nil {
				log.Printf("Someting wrong with request to tlg")
				log.Print(err)
				return err
			}

			// delete cache from

			CacheSystem.ChangeBusyStatus(update.Message.User.Id)

			CacheSystem.Delete(session.companionUserId)

			return nil

		}

		log.Printf("operator with id - %d send message to user with id - %d", update.Message.User.Id, session.companionUserId)

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

		reqToTlg := OutMessage{
			Text:   update.Message.Text,
			ChatId: session.companionUserId,
		}

		log.Printf("user with id - %d send message to operator with id - %d", update.Message.User.Id, session.companionUserId)

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

func mainPolicy(update UpdateType, path string) {

	// check cache
	cache, isOldSession := CacheSystem.Get(update.Message.User.Id)

	// if request from operator bot
	if path == "/operator" {
		// it is old session?
		if isOldSession {
			log.Printf("we in old session for operator")
			if cache.auth {
				log.Printf("and operator has logging in")
				// it is not bot mode? if yes we send text user
				if !cache.botStatus {
					log.Printf("bot mode is false for operator")
					_ = policyOperatorBot(update, path)
					// if it bot mode, i don't know why operator send text:D
				} else {
					log.Printf("bot mode is true for operator")

					reqToTlg := OutMessage{
						Text:   "Активный диалогов нет:) Отдыхайте!😍",
						ChatId: update.Message.Chat.Id,
					}
					// send req to tlg
					_ = sendReqToTlg(BuildUrl(PathSendMessage, BotsInfo["operator"]), reqToTlg)
				}
			} else {
				if update.Message.Text == "lolkaperduska" {
					reqToTlg := OutMessage{
						Text:   "Вы успешно авторизовались, молодец! Ожидайте диалога с клиентом!😍",
						ChatId: update.Message.Chat.Id,
					}
					// send req to tlg
					_ = sendReqToTlg(BuildUrl(PathSendMessage, BotsInfo["operator"]), reqToTlg)

					CacheSystem.ChangeAuthStatus(update.Message.User.Id)

				} else {

					reqToTlg := OutMessage{
						Text:   "Пароль неправильный, попытайтесь еще раз. Будьте старательны👹",
						ChatId: update.Message.Chat.Id,
					}
					// send req to tlg
					_ = sendReqToTlg(BuildUrl(PathSendMessage, BotsInfo["operator"]), reqToTlg)

				}
			}
		} else {
			log.Printf("for operator is a new session")

			// in this place i should generate session data for bot and i know that operator not auth
			// create new session data
			session := "bot-" + time.Now().Format("20060102150405")
			CacheSystem.Put(update.Message.User.Id, sessionData{
				messageId: 0,
				sessionId: session,
				botStatus: true,
				auth:      false,
			})

			reqToTlg := OutMessage{
				Text:   "Вы не авторизовались, войдите в систему чтобы начать обслиживать клиентов!😍",
				ChatId: update.Message.Chat.Id,
			}
			// send req to tlg
			_ = sendReqToTlg(BuildUrl(PathSendMessage, BotsInfo["operator"]), reqToTlg)

		}

	} else {
		if isOldSession {
			log.Printf("we in old session for user")
			if cache.botStatus {
				log.Printf("bot status is true for user")
				_ = policyTlgSm(update)
			} else {
				log.Printf("bot status is false for user")
				_ = policyOperatorBot(update, path)
			}
		} else {
			log.Printf("for user is new session")

			// create new session data
			session := "bot-" + time.Now().Format("20060102150405")
			CacheSystem.Put(update.Message.User.Id, sessionData{
				messageId: 0,
				sessionId: session,
				botStatus: true,
			})
			_ = policyTlgSm(update)

		}
	}

}

// Метод Handler. Данный метод будет обрабатывать HTTP запросы поступающие к функции
func Handler(w http.ResponseWriter, r *http.Request) {

	// get message from tlg
	if r.Body != nil {
		defer r.Body.Close()
	}

	decoder := json.NewDecoder(r.Body)

	log.Println(decoder)

	var update UpdateType
	err := decoder.Decode(&update)
	if err != nil {
		panic(err)
	}

	// Логирование входящего запроса
	log.Printf("Request received: %s\nMethod: %s\nPATH: %s\nRAW_PATH: %s\nRAW_QUERY:%s", update.Message.Text, r.Method, r.URL.Path, r.URL.RawPath, r.URL.RawQuery)

	mainPolicy(update, r.URL.Path)

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
