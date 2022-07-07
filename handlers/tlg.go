package handlers

//import (
//	"bytes"
//	"encoding/json"
//	"log"
//	"net/http"
//)

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

func BuildUrl(param string, token string) string {
	return API_URL + token + param
}
