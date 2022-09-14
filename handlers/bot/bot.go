package handlers

import (
	"animeGAN2go/MessageTypes"
	"animeGAN2go/plugins"
	"animeGAN2go/structure"
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
)

func SendPhoto(chatId int, image string) string {
	// TODO отрефакторить это
	bot, err := tgbotapi.NewBotAPI(plugins.Token)
	if err != nil {
		log.Panic(err)
	}
	sEnc, _ := b64.StdEncoding.DecodeString(image)

	photoFileBytes := tgbotapi.FileBytes{
		Name:  "picture",
		Bytes: sEnc,
	}

	message, _ := bot.Send(tgbotapi.NewPhotoUpload(int64(chatId), photoFileBytes))

	f := *(message.Photo)

	var largerNumber, temp int
	// TODO найти самый большой элемент и передать его
	for position, element := range f {
		if element.FileSize > temp {
			temp = element.FileSize
			largerNumber = position
		}
	}
	return f[largerNumber].FileID
}

func EditMessage(chatId int, text string, messageId int) {
	url := "https://api.telegram.org/bot" + plugins.Token + "/editMessageText"

	d := &structure.EditDataToTlg{ChatId: chatId, Text: text, MessageId: messageId}

	jsonString, err := json.Marshal(d)

	if err != nil {
		log.Fatal(err)
	}

	r := bytes.NewReader(jsonString)

	contentType := "application/json"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, r)

	req.Header.Set("Content-Type", contentType)

	client.Do(req)
}

func GetImage(path string) string {
	url := "https://api.telegram.org/file/bot" + plugins.Token + "/" + path

	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		sEnc := b64.StdEncoding.EncodeToString([]byte(bodyBytes))
		return sEnc
	}
	return "-1"
}

func GetFilePath(token string) string {
	// отправляем запрос на получение файла
	url := "https://api.telegram.org/bot" + plugins.Token + "/getFile?file_id=" + token
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var data MessageTypes.GetFilePath

	if resp.StatusCode == http.StatusOK {

		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&data)
		if err != nil {
			//log.Fatal(err)
			log.Println("отправялем сообщение пользователю что что то не так")
		}
		return data.Result.FilePath
	}
	return "-1"
}

func SendMessage(chatId int, text string) MessageTypes.RespDataTlg {

	var respData MessageTypes.RespDataTlg

	url := "https://api.telegram.org/bot" + plugins.Token + "/sendMessage"

	d := &structure.SendDataToTlg{ChatId: chatId, Text: text}

	jsonString, err := json.Marshal(d)

	if err != nil {
		log.Fatal(err)
	}

	r := bytes.NewReader(jsonString)

	contentType := "application/json"
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, r)

	req.Header.Set("Content-Type", contentType)

	res, err := client.Do(req)

	bodyBytes, err := ioutil.ReadAll(res.Body)

	errUnm := json.Unmarshal(bodyBytes, &respData)

	if errUnm != nil {
		log.Fatal(errUnm)
	}

	if res.StatusCode == http.StatusOK {
		fmt.Println("ok, delivered")
	} else {
		fmt.Println("not ok, someting wrong")
	}
	return respData

}
