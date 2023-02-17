package handlers

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"

	"golang.org/x/net/websocket"
)

var debugMode = flag.String("d", "off", "debug mode")

func initWebsocketClient(userId int, projectId string, chatId int) {
	botParams, _ := BotsParams.GetData(projectId)
	CacheSystemUser, _ := CacheUser.GetData(projectId)

	s, _ := CacheSystemUser.Get(userId)

	// flag.Parse()

	fmt.Println("get token for messenger")

	urlPath, exist := AuthServiceHost[botParams.standType][botParams.botType]

	if !exist {
		fmt.Printf(
			"stand or partner with params: %s, %s is not exist",
			botParams.standType,
			botParams.botType,
		)

		// check exist params for
		return

	}

	messengerHost, exs := MessengerEntryPoint[botParams.standType]

	if !exs {
		fmt.Printf("stand - %s is not exists")
		return
	}

	// TODO check success connection, if proplems retry
	res := reqPostMessenger(urlPath)

	sessionId := res["session_id"]
	URL := fmt.Sprintf("wss://%s/api/ws/prelogin/%s", messengerHost, sessionId)
	ORIGIN := fmt.Sprintf("https://%s/", messengerHost)
	config, _ := websocket.NewConfig(URL, ORIGIN)
	fistIncomingMessage := true

	config.TlsConfig = &tls.Config{InsecureSkipVerify: true}

	fmt.Println("Starting Client")

	ws, err := websocket.DialConfig(config)

	if err != nil {
		fmt.Printf("Dial failed: %s\n", err.Error())
		return
	}

	go readClientMessages(ws, s.incomingMessages)

	for {
		select {
		case response, ok := <-s.messageToSend:

			if !ok {
				fmt.Println("chat old, close ws connection and stop goroutine")
				ws.Close()
				return
			}

			err = websocket.JSON.Send(ws, response)
			if err != nil {
				fmt.Printf("Send failed: %s\n", err.Error())
				// TODO in this plase we should reconnect
				return
			}
		case message := <-s.incomingMessages:
			if *debugMode == "on" {
				fmt.Println(`Message Received:`, message)
			}

			data := Messenger{}
			json.Unmarshal([]byte(message), &data)
			if fistIncomingMessage {

				CacheSystemUser.ChangeCompanionId(userId, int(data.Data.Conversation.Id))

				fistIncomingMessage = false

			}

			if data.Method == "receive_text_message" && (data.Data.Author.Type == "BOT" || data.Data.Author.Type == "OPERATOR") {
				fmt.Println("Answer from bot or operator:", data.Data.Text)

				if data.Data.Text == "Оператор завершил чат" {
					// Oper close chat, delete cash and close connection
					CacheSystemUser.Delete(userId)

					ws.Close()
					return

				}

				// send text to user
				reqToTlg := OutMessage{
					Text:   data.Data.Text,
					ChatId: chatId,
				}

				sendReqToTlg(BuildUrl(PathSendMessage, botParams.bot), reqToTlg)

			}

		}
	}
}

func readClientMessages(ws *websocket.Conn, incomingMessages chan string) {
	for {
		var message string
		// err := websocket.JSON.Receive(ws, &message)
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			// TODo check erros from messenger
			//"method":"ack","error":{"status_code":500,"reason":"Internal error","url":null},"on_method":"send_text_message"
			// TODO in this plase we should reconnect
			fmt.Printf("Error::: %s\n", err.Error())
			return
		}
		incomingMessages <- message
	}
}
