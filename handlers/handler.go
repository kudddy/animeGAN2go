package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

const urlPathToSkill = "https://smartapp-code.sberdevices.ru/chatadapter/chatapi/webhook/sber_nlp2/ZMgoqvmH:abf05f2ca8543405adad9b5bce52b548496dc2b8"

func policyTlgSm(update UpdateType, params map[string]string) error {

	session, _ := CacheSystem.Get(update.Message.User.Id)

	// convert message to sm
	reqToSm := generatePayloadForSm(update, session)

	// send message to sm and get resp
	resp, err := sendReqToSm(params["sm-webhook"], reqToSm)
	if err != nil {
		log.Printf("Someting wrong with request to SM with mid - %d", session.messageId)
	}
	// convert message to tlg format
	var textToUser string
	var extraText string
	var buts []Buttons

	if resp.MessageName == "ANSWER_TO_USER" {

		textToUser, extraText, buts = resp.processRespFromSm()

		var reqToTlg OutMessage

		if len(buts) > 0 {

			var buttons []InlineKeyboardButton

			var inlineKey = InlineKeyboardButton{
				Text: buts[0].text,
				URL:  buts[0].url,
			}
			buttons = append(buttons, inlineKey)

			var arrayOfButtons [][]InlineKeyboardButton

			arrayOfButtons = append(arrayOfButtons, buttons)

			var inlineButtons = InlineKeyboardMarkup{
				InlineKeyboard: arrayOfButtons,
			}

			reqToTlg = OutMessage{
				Text:        textToUser + "\n" + extraText,
				ChatId:      update.Message.Chat.Id,
				ReplyMarkup: &inlineButtons,
			}
		} else {
			reqToTlg = OutMessage{
				Text:   textToUser + "\n" + extraText,
				ChatId: update.Message.Chat.Id,
			}
		}

		err = sendReqToTlg(BuildUrl(PathSendMessage, params["bot"]), reqToTlg)

		if err != nil {
			log.Printf("Someting wrong with request to tlg")
			log.Print(err)
			return err
		}

		return nil

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

	err = sendReqToTlg(BuildUrl(PathSendMessage, params["bot"]), reqToTlg)

	if err != nil {
		log.Printf("Someting wrong with request to tlg")
		log.Print(err)
		return err
	}
	return nil
}

func policyOperator(update UpdateType, params map[string]string) error {

	session, isOldSession := CacheSystem.Get(update.Message.User.Id)

	if isOldSession {
		log.Printf("we in old session for operator")
		if session.auth {
			log.Printf("and operator has logging in")
			// it is not bot mode? if yes we send text user
			if !session.botStatus {
				log.Printf("bot mode is false for operator")

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

				// _ = policyOperatorBot(update, botType)

				// if it bot mode, i don't know why operator send text:D
			} else {
				log.Printf("bot mode is true for operator")

				reqToTlg := OutMessage{
					Text:   "Активный диалогов нет:) Отдыхайте!😍",
					ChatId: update.Message.Chat.Id,
				}
				// send req to tlg
				_ = sendReqToTlg(BuildUrl(PathSendMessage, params["operator"]), reqToTlg)
			}
		} else {
			if update.Message.Text == "pass" {
				reqToTlg := OutMessage{
					Text:   "Вы успешно авторизовались, молодец! Ожидайте диалога с клиентом!😍",
					ChatId: update.Message.Chat.Id,
				}
				// send req to tlg
				_ = sendReqToTlg(BuildUrl(PathSendMessage, params["operator"]), reqToTlg)

				CacheSystem.ChangeAuthStatus(update.Message.User.Id)

			} else {

				reqToTlg := OutMessage{
					Text:   "Пароль неправильный, попытайтесь еще раз. Будьте старательны👹",
					ChatId: update.Message.Chat.Id,
				}
				// send req to tlg
				_ = sendReqToTlg(BuildUrl(PathSendMessage, params["operator"]), reqToTlg)

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
			Text:   "Вы не авторизовались, войдите в систему чтобы начать обслуживать клиентов!😍",
			ChatId: update.Message.Chat.Id,
		}
		// send req to tlg
		_ = sendReqToTlg(BuildUrl(PathSendMessage, params["operator"]), reqToTlg)

	}
	return nil
}

func policyUser(update UpdateType, params map[string]string) {

	session, isOldSession := CacheSystem.Get(update.Message.User.Id)

	if isOldSession {
		log.Printf("we in old session for user")
		if session.botStatus {

			// change session status
			CacheSystem.ChangeSessionStatus(update.Message.User.Id)

			log.Printf("bot status is true for user")
			_ = policyTlgSm(update, params)
		} else {
			log.Printf("bot status is false for user")
			// TODO здесь нужно обработать запрос

			reqToTlg := OutMessage{
				Text:   update.Message.Text,
				ChatId: session.companionUserId,
			}

			log.Printf("user with id - %d send message to operator with id - %d", update.Message.User.Id, session.companionUserId)

			// send req to tlg
			err := sendReqToTlg(BuildUrl(PathSendMessage, params["operator"]), reqToTlg)
			if err != nil {
				log.Printf("Someting wrong with request to tlg")
				log.Print(err)
			}
		}
	} else {
		log.Printf("for user is new session")
		// create new session data
		session := "bot-" + time.Now().Format("20060102150405")
		CacheSystem.Put(update.Message.User.Id, sessionData{
			messageId:  0,
			sessionId:  session,
			botStatus:  true,
			newSession: true,
		})
		_ = policyTlgSm(update, params)

	}

}

func mainPolicy(update UpdateType, botType string, projectId string) (status bool, desc string) {

	params, ok := BotsParams.GetData(projectId)

	if !ok {
		return false, "project not found"
	}

	//// check cache
	//cache, isOldSession := CacheSystem.Get(update.Message.User.Id)

	// if request from operator bot
	if botType == "operator" {
		// TODO check errors
		_ = policyOperator(update, params)
		return true, "message from operator success processed"

	} else if botType == "bot" {
		policyUser(update, params)
		return true, "message from user success processed"
	} else {
		return false, "error, url is wrong"
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

	params := strings.Split(r.URL.Path, "/")

	var status bool
	var desc string

	if len(params) != 3 {

		status = false
		desc = "bad request"

	} else {

		projectId := params[1]
		botType := params[2]

		status, desc = mainPolicy(update, botType, projectId)

	}

	js, err := json.Marshal(RespByServ{
		Ok:   status,
		Desc: desc,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}
