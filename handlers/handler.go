package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)


const urlPathToSkill = "https://smartapp-code.sberdevices.ru/chatadapter/chatapi/webhook/sber_nlp2/cGnGPZWb:45c9c4e54edfcf2cfe505f84e3f338185a334e42"

func generatePayloadForSm(text string) ReqToSmType {

	messageID := 45345345
	sessionId := "avaya-lolo-fdsfsdfaf-fsdfasdfsdf"
	messageName := "MESSAGE_TO_SKILL"

	userUuid := Uuid{
		UserId: "9485D45E-466E-4852-B5DA-1A27DFF5EFC8",
		Sub: "1hkmItxUo6BDBmNvGM7inj4kNvWIRyQOaUzWdlqxYafPUqNZ/fTLMJ8M4idi1y467byHIwH8zAnbqt6glUevV0d8+tppO2Ysr1Ryn5PPj7nkk+7kTtDC1MnJvZVaJP3uzHxG5PPxvQpIbtQccKxegw==",
		UserChannel: "SBOL",
	}

	appInfo := appInfo{
		ProjectId: "12f20e40-efc6-4ff5-9179-f5c51f7197b3",
		ApplicationId: "7aa5ae84-c668-4e24-94d8-e35cf053e7a1",
		AppversionId: "bbddbed8-a8c6-483f-99b5-516dbae4ea70",
		FrontendType: "DIALOG",
		AgeLimit: 18,
		AffiliationType: "ECOSYSTEM",

	}


	message := message{
		OriginalText: text,
		NormalizedText: text,
		OriginalMessageName: "MESSAGE_FROM_USER",
		HumanNormalizedText: text,
		HumanNormalizedTextWithAnaphora: text,
	}

	payload := payload{
		Intent: "sberauto_main",
		OriginalIntent: "food",
		NewSession: false,
		ApplicationId: "7aa5ae84-c668-4e24-94d8-e35cf053e7a1",
		AppversionId: "bbddbed8-a8c6-483f-99b5-516dbae4ea70",
		ProjectName: "СберАвто. Подбор автомобиля",
		AppInfo: appInfo,
		Msg: message,

	}

	reqToSmType := ReqToSmType{
		MessageId: messageID,
		SessionId: sessionId,
		MessageName: messageName,
		Payload: payload,
		Uuid: userUuid,

	}

	return reqToSmType

}

//update.Message.Text
//
//outPhoto := OutPhoto{
//ChatId:  chatId,
//Photo:   srcUrl,
//Caption: caption,
//// "entities":[{"offset":10,"length":4,"type":"hashtag"},{"offset":15,"length":48,"type":"url"}]}
//}

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
	log.Printf("Request received: %s\nMethod: %s", update.Message.Text, r.Method)

	// convert message to sm
	reqToSm := generatePayloadForSm(update.Message.Text)


	// send message to sm and get resp

	resp, _ := sendReqToSm(urlPathToSkill, reqToSm)

	// convert message to tlg format

	var textToUser string

	if resp.MessageName == "ANSWER_TO_USER"{
		textToUser = resp.Payload.PronounceText
	} else if resp.MessageName == "NOTHING_FOUND" {
		textToUser = "Переадресую на оператора"
	}

	reqToTlg := OutMessage{
		Text: textToUser,
		ChatId: 81432612,

	}


	// send req to tlg

	err = sendReqToTlg(BuildUrl(PathSendMessage), reqToTlg)

	if err != nil {
		log.Printf("Someting wrong with request to tlg")
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