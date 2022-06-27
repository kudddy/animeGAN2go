package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

const (
	PathSetWebhook            = "/setWebhook"
	PathSendMessage           = "/sendMessage"
	PathSendPhoto             = "/sendPhoto"
	PathSendVideo             = "/sendVideo"
	PathSendMediaGroup        = "/sendMediaGroup"
	PathSetMyCommands         = "/setMyCommands"
	PathGetChatAdministrators = "/getChatAdministrators"
)

const API_URL = "https://api.telegram.org/bot"

func BuildUrl(param string) string {
	return API_URL + BotsInfo["bot"] + param
}

func sendJson(urlPath string, outData interface{}) error {
	body, err := json.Marshal(outData)
	if err != nil {
		log.Fatal(err)
		return err
	}
	bodyBytes := bytes.NewBuffer(body)
	log.Println(bodyBytes)
	_, err = http.Post(BuildUrl(urlPath), "application/json", bodyBytes)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func SendMessage(msg Message, text string) error {
	if text == "" {
		text = ">>" + msg.Text
	}
	outData := OutMessage{
		ChatId:        msg.Chat.Id,
		Text:          text,
		ReplayToMsgId: msg.Id,
	}
	return sendJson(PathSendMessage, outData)
}
