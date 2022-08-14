package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// API errors
//TODO create generate auth token method
const (
	// TypeBot types of methods to services
	TypeBot           = "bot"
	TypeOperator      = "operator"
	TypeUpdateProject = "update"
)

func (update *UpdateType) policyTlgSm(projectId string) error {
	botParams, _ := BotsParams.GetData(projectId)
	CacheSystemUser, _ := CacheUser.GetData(projectId)
	CacheSystemOperator, _ := CacheOperator.GetData(projectId)

	// get user session
	session, _ := CacheSystemUser.Get(update.Message.User.Id)

	// convert message to sm
	reqToSm := update.generatePayloadForSm(session)

	// send message to sm and get resp
	resp, err := sendReqToSm(botParams.webhook, reqToSm)
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

		_, err = sendReqToTlg(BuildUrl(PathSendMessage, botParams.bot), reqToTlg)

		if err != nil {
			log.Printf("Someting wrong with request to tlg")
			log.Print(err)
			return err
		}

		return nil

	} else if resp.MessageName == "NOTHING_FOUND" {

		// in this place we should get from db user_id with max score
		operators := CacheSystemOperator.GetRandomAuthOperators()
		if len(operators) > 0 {
			// first in, first out
			var operatorBotId = operators[0]

			operatorSession, _ := CacheSystemOperator.Get(operatorBotId)

			// and create new session data for bot

			//session, _ := CacheSystem.Get(string(rune(update.Message.User.Id)))

			log.Printf("save session parametrs for operator when id is %d, companion id is %d", operatorBotId, update.Message.User.Id)
			// TODO its not true
			//CacheSystem.ChangeBotStatus(update.Message.User.Id)
			CacheSystemUser.ChangeBotStatus(update.Message.User.Id)

			CacheSystemOperator.Put(operatorBotId, sessionData{
				messageId:       0,
				sessionId:       operatorSession.sessionId,
				botStatus:       false,
				companionUserId: update.Message.User.Id,
				auth:            operatorSession.auth,
				busy:            true,
			})

			log.Printf("save session parametrs for bot when id is %d, companion id is %d", update.Message.User.Id, operatorBotId)
			// update user session param as companionUserId
			//CacheSystem.Put(update.Message.User.Id, sessionData{
			//	messageId:       session.messageId,
			//	sessionId:       session.sessionId,
			//	botStatus:       false,
			//	companionUserId: operatorBotId,
			//})

			CacheSystemUser.Put(update.Message.User.Id, sessionData{
				messageId:       session.messageId,
				sessionId:       session.sessionId,
				botStatus:       false,
				companionUserId: operatorBotId,
			})

			s, _ := CacheSystemUser.Get(update.Message.User.Id)

			log.Printf("session parameters from cache for %d is %d", update.Message.User.Id, s.companionUserId)

			d, _ := CacheSystemOperator.Get(operatorBotId)

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

	_, err = sendReqToTlg(BuildUrl(PathSendMessage, botParams.bot), reqToTlg)

	if err != nil {
		log.Printf("Someting wrong with request to tlg")
		log.Print(err)
		return err
	}
	return nil
}

func (update *UpdateType) policyOperator(projectId string) error {
	CacheSystemUser, _ := CacheUser.GetData(projectId)
	CacheSystemOperator, _ := CacheOperator.GetData(projectId)
	botParams, _ := BotsParams.GetData(projectId)

	session, isOldSession := CacheSystemOperator.Get(update.Message.User.Id)

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
					_, err := sendReqToTlg(BuildUrl(PathSendMessage, botParams.bot), reqToTlg)

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
					_, err = sendReqToTlg(BuildUrl(PathSendMessage, botParams.operator), reqToTlg)

					if err != nil {
						log.Printf("Someting wrong with request to tlg")
						log.Print(err)
						return err
					}

					// delete cache from

					CacheSystemOperator.ChangeBusyStatus(update.Message.User.Id)

					CacheSystemUser.Delete(session.companionUserId)

					return nil

				}

				log.Printf("operator with id - %d send message to user with id - %d", update.Message.User.Id, session.companionUserId)

				reqToTlg := OutMessage{
					Text:   update.Message.Text,
					ChatId: session.companionUserId,
				}

				// send req to tlg
				_, err := sendReqToTlg(BuildUrl(PathSendMessage, botParams.bot), reqToTlg)

				if err != nil {
					log.Printf("Someting wrong with request to tlg")
					log.Print(err)
					return err
				}

				// _ = policyOperatorBot(update, botType)

				// if it is bot mode, i don't know why operator send text:D
			} else {
				log.Printf("bot mode is true for operator")

				reqToTlg := OutMessage{
					Text:   "Активный диалогов нет:) Отдыхайте!😍",
					ChatId: update.Message.Chat.Id,
				}
				// send req to tlg
				_, _ = sendReqToTlg(BuildUrl(PathSendMessage, botParams.operator), reqToTlg)
			}
		} else {
			if update.Message.Text == "pass" {
				reqToTlg := OutMessage{
					Text:   "Вы успешно авторизовались, молодец! Ожидайте диалога с клиентом!😍",
					ChatId: update.Message.Chat.Id,
				}
				// send req to tlg
				_, _ = sendReqToTlg(BuildUrl(PathSendMessage, botParams.operator), reqToTlg)

				CacheSystemOperator.ChangeAuthStatus(update.Message.User.Id)

			} else {

				reqToTlg := OutMessage{
					Text:   "Пароль неправильный, попытайтесь еще раз. Будьте старательны👹",
					ChatId: update.Message.Chat.Id,
				}
				// send req to tlg
				_, _ = sendReqToTlg(BuildUrl(PathSendMessage, botParams.operator), reqToTlg)

			}
		}
	} else {
		log.Printf("for operator is a new session")

		// in this place i should generate session data for bot and i know that operator not auth
		// create new session data
		session := "bot-" + time.Now().Format("2017-09-07 17:06:04.000000")
		CacheSystemOperator.Put(update.Message.User.Id, sessionData{
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
		_, _ = sendReqToTlg(BuildUrl(PathSendMessage, botParams.operator), reqToTlg)

	}
	return nil
}

func (update *UpdateType) policyUser(projectId string) {
	CacheSystemUser, _ := CacheUser.GetData(projectId)
	botParams, _ := BotsParams.GetData(projectId)

	session, isOldSession := CacheSystemUser.Get(update.Message.User.Id)

	if isOldSession {
		log.Printf("we in old session for user")
		if session.botStatus {

			// change session status
			CacheSystemUser.ChangeSessionStatus(update.Message.User.Id)

			log.Printf("bot status is true for user")
			_ = update.policyTlgSm(projectId)
		} else {
			log.Printf("bot status is false for user")
			// TODO здесь нужно обработать запрос

			reqToTlg := OutMessage{
				Text:   update.Message.Text,
				ChatId: session.companionUserId,
			}

			log.Printf("user with id - %d send message to operator with id - %d", update.Message.User.Id, session.companionUserId)

			// send req to tlg
			_, err := sendReqToTlg(BuildUrl(PathSendMessage, botParams.operator), reqToTlg)
			if err != nil {
				log.Printf("Someting wrong with request to tlg")
				log.Print(err)
			}
		}
	} else {
		log.Printf("for user is new session")
		// create new session data
		session := "bot-" + time.Now().Format("2017-09-07 17:06:04.000000")
		CacheSystemUser.Put(update.Message.User.Id, sessionData{
			messageId:  0,
			sessionId:  session,
			botStatus:  true,
			newSession: true,
		})
		_ = update.policyTlgSm(projectId)

	}

}

func (update *UpdateType) mainPolicy(botType string, projectId string) (status bool, desc string) {

	//params, ok := BotsParams.GetData(projectId)

	// if request from operator bot
	if botType == TypeOperator {
		// TODO check errors
		_ = update.policyOperator(projectId)
		return true, "message from operator success processed"
	} else if botType == TypeBot {
		// TODO check errors
		update.policyUser(projectId)
		return true, "message from user success processed"
	} else {
		return false, "error, url is wrong"
	}

}

func (update *UpdateBotsParams) updateBotsParams(projectId string) (bool, string) {

	//TODO in this place we should add request to tlg and registration webhook for all tokens

	serviceHostBot := APIEndpoint + projectId + "/bot"

	data, err := sendReqToTlg(BuildUrl(PathSetWebhook, update.Bot)+"?url="+serviceHostBot, OutMessage{})

	if err != nil {
		log.Printf("someting wrong with resp to tlg where we update webhook")
		return false, data.Description
	}
	serviceHostOperator := APIEndpoint + projectId + "/operator"

	data, err = sendReqToTlg(BuildUrl(PathSetWebhook, update.Operator)+"?url="+serviceHostOperator, OutMessage{})

	if err != nil {
		log.Printf("someting wrong with resp to tlg where we update webhook")
		return false, data.Description
	}

	BotsParams.AddData(projectId, botsInfo{
		update.Bot,
		update.Operator,
		update.Webhook,
	})

	return true, data.Description
}

// Handler Метод Handler. Данный метод будет обрабатывать HTTP запросы поступающие к функции
func Handler(w http.ResponseWriter, r *http.Request) {

	// get message from tlg
	if r.Body != nil {
		defer r.Body.Close()
	}

	decoder := json.NewDecoder(r.Body)

	params := strings.Split(r.URL.Path, "/")
	var status bool
	var desc string
	if len(params) != 3 {
		status = false
		desc = "bad request"
	} else {
		projectId := params[1]

		// check auth token
		ok := contains(AuthTokens, projectId)

		if ok {
			method := params[2]

			if method == TypeBot || method == TypeOperator {

				var update UpdateType
				err := decoder.Decode(&update)
				// Логирование входящего запроса
				log.Printf("Request received: %s\nMethod: %s\nPATH: %s\nRAW_PATH: %s\nRAW_QUERY:%s", update.Message.Text, r.Method, r.URL.Path, r.URL.RawPath, r.URL.RawQuery)
				if err != nil {
					panic(err)
				}
				status, desc = update.mainPolicy(method, projectId)
			} else if method == TypeUpdateProject {

				var update UpdateBotsParams

				err := decoder.Decode(&update)
				if err != nil {
					panic(err)
				}

				status, desc = update.updateBotsParams(projectId)

			} else {
				status = false
				desc = "bad method format"
			}
		} else {
			status = false
			desc = "not valid auth token"
		}

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
