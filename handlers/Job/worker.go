package Job

import (
	"animeGAN2go/MessageTypes"
	"animeGAN2go/bot"
	"animeGAN2go/ganserv"
	"animeGAN2go/plugins"
	"animeGAN2go/plugins/pg"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func StartWorker(t MessageTypes.ReqData) {
	// бесконечный цикл
	// получили путь к файлу

	fmt.Println("Получаем файл id")
	filePath := bot.GetFilePath(t.FileId)
	// получаем путь
	fmt.Println("получаем изображение")
	image := bot.GetImage(filePath)
	// получили изображение в base64

	fmt.Println("Отправляем изображение в модель")
	d := ganserv.SendImageToModel(image, t.UserModel)

	var dataFromTlg MessageTypes.RespDataTlg
	if plugins.IsZeroOfUnderlyingType(d) {
		text := "Упс, с датацентром что то не так, повторите попытку чуть позже. Мы уже занимаемся решением этой проблемы!"
		dataFromTlg = bot.SendMessage(t.ChatId, text)
	} else {

		i := 0
		for {

			time.Sleep(1 * time.Second)

			fmt.Println("в цикле")

			data, dataQueen, queen, err := ganserv.GetQueenNumber(d.Hash)
			if !err {
				if queen {
					if dataQueen.Status == "QUEUED" {
						text := fmt.Sprintf("Ваша очередь: %s", strconv.Itoa(dataQueen.Data))
						if i == 0 {
							dataFromTlg = bot.SendMessage(t.ChatId, text)
						} else {
							bot.EditMessage(t.ChatId, text, int(dataFromTlg.Result.MessageId))
						}

						i++
					}
				} else {
					if data.Status == "COMPLETE" {
						//fmt.Println("Отправляем пользователю сообщение с фотографией")
						// Отправляем пользователю сообщение с фотографией
						imageString := strings.Split(data.Data.Data[0], ",")[1]
						bot.SendPhoto(t.ChatId, imageString)
						pg.InsertCancelAction(t.UserID)

						break
					}
				}
			} else {
				text := "Что то пошло не так:( Попробуйте загрузить другое фото!"
				dataFromTlg = bot.SendMessage(t.ChatId, text)
				break
			}

		}

	}

}
