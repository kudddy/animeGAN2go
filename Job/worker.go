package Job

import (
	"animeGAN2go/MessageTypes"
	"animeGAN2go/plugins"
	"animeGAN2go/plugins/pg"
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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

type SendDataToPush struct {
	Data   []string `json:"data"`
	Action string   `json:"action"`
}

func SendImageToModel(sEncPhoto string, userModel string) MessageTypes.GetModelHash {

	var arr []string
	arr = append(arr, "data:image/jpeg;base64,"+sEncPhoto)

	fmt.Println(userModel)

	arr = append(arr, userModel)

	d := &SendDataToPush{Data: arr, Action: "predict"}

	jsonString, err := json.Marshal(d)
	if err != nil {
		fmt.Println(err)
	}

	r := bytes.NewReader(jsonString)

	url := "https://hf.space/gradioiframe/akhaliq/AnimeGANv2/api/queue/push/"
	//url:= "http://0.0.0.0:8080/push/"

	var client http.Client
	contentType := "application/json"
	resp, err := client.Post(url, contentType, r)

	var data MessageTypes.GetModelHash
	if err != nil {
		fmt.Println("ошибка при отправлке запроса в модель")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Все ок, код положительный")
		decoder := json.NewDecoder(resp.Body)

		err = decoder.Decode(&data)
		return data
	}
	fmt.Println("Странный код запроса")
	return data
}

type SendDataStatus struct {
	Hash string `json:"hash"`
}

func GetQueenNumber(hash string) (MessageTypes.CheckStatus, MessageTypes.CheckStatusQueen, bool, bool) {

	d := &SendDataStatus{Hash: hash}

	jsonString, err := json.Marshal(d)

	r := bytes.NewReader(jsonString)

	url := "https://hf.space/gradioiframe/akhaliq/AnimeGANv2/api/queue/status/"

	//url := "http://0.0.0.0:8080/status/"

	contentType := "application/json"

	var client http.Client
	resp, err := client.Post(url, contentType, r)
	// для готово результата модели
	var data MessageTypes.CheckStatus
	// для очереди
	var dataQueen MessageTypes.CheckStatusQueen

	var queen bool

	var globalError bool

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {

		errDec := json.Unmarshal(bodyBytes, &data)

		if errDec != nil {
			errUnm := json.Unmarshal(bodyBytes, &dataQueen)

			if errUnm != nil {
				globalError = true
				return data, dataQueen, queen, globalError
			}
			globalError = false
			queen = true
			return data, dataQueen, queen, globalError
		}
		queen = false
		globalError = false
		return data, dataQueen, queen, globalError
	}
	globalError = false
	return data, dataQueen, queen, globalError

}

type SendDataToTlg struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

func SendMessage(chatId int, text string) MessageTypes.RespDataTlg {

	var respData MessageTypes.RespDataTlg

	url := "https://api.telegram.org/bot" + plugins.Token + "/sendMessage"

	d := &SendDataToTlg{ChatId: chatId, Text: text}

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

type EditDataToTlg struct {
	Text      string `json:"text"`
	ChatId    int    `json:"chat_id"`
	MessageId int    `json:"message_id"`
}

func EditMessage(chatId int, text string, messageId int) {
	url := "https://api.telegram.org/bot" + plugins.Token + "/editMessageText"
	d := &EditDataToTlg{ChatId: chatId, Text: text, MessageId: messageId}

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

func SendPhoto(chatId int, image string) {
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
	_, err = bot.Send(tgbotapi.NewPhotoUpload(int64(chatId), photoFileBytes))

}

func StartWorker(t MessageTypes.ReqData) {
	// бесконечный цикл
	// получили путь к файлу

	fmt.Println("Получаем файл id")
	filePath := GetFilePath(t.FileId)
	// получаем путь
	fmt.Println("получаем изображение")
	image := GetImage(filePath)
	// получили изображение в base64

	fmt.Println("Отправляем изображение в модель")
	d := SendImageToModel(image, t.UserModel)

	var dataFromTlg MessageTypes.RespDataTlg
	if plugins.IsZeroOfUnderlyingType(d) {
		text := "Упс, с датацентром что то не так, повторите попытку чуть позже. Мы уже занимаемся решением этой проблемы!"
		dataFromTlg = SendMessage(t.ChatId, text)
	} else {

		i := 0
		for {

			time.Sleep(1 * time.Second)

			fmt.Println("в цикле")

			data, dataQueen, queen, err := GetQueenNumber(d.Hash)
			if !err {
				if queen {
					if dataQueen.Status == "QUEUED" {
						text := fmt.Sprintf("Ваша очередь: %s", strconv.Itoa(dataQueen.Data))
						if i == 0 {
							dataFromTlg = SendMessage(t.ChatId, text)
						} else {
							EditMessage(t.ChatId, text, int(dataFromTlg.Result.MessageId))
						}

						i++
					}
				} else {
					if data.Status == "COMPLETE" {
						//fmt.Println("Отправляем пользователю сообщение с фотографией")
						// Отправляем пользователю сообщение с фотографией
						imageString := strings.Split(data.Data.Data[0], ",")[1]
						SendPhoto(t.ChatId, imageString)
						pg.InsertCancelAction(t.UserID)

						break
					}
				}
			} else {
				text := "Что то пошло не так:( Попробуйте загрузить другое фото!"
				dataFromTlg = SendMessage(t.ChatId, text)
				break
			}

		}

	}
	// отправляем файл на машину с преобразователем
	// отправляем сообщение в push

}
