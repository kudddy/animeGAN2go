package Job

import (
	"animeGAN2go/MessageTypes"
	"animeGAN2go/bot"
	"animeGAN2go/ganserv"
	"animeGAN2go/plugins"
	"animeGAN2go/rds"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// задача получить из очереди сообщения с файлами, обработать их и передать пути к файлам к следующему воркеру
//

func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func StartSingleWorker() {
	fmt.Println("get file ids from queen")

	for {
		fmt.Println("tic")
		time.Sleep(3 * time.Second)
		fmt.Println("tock")

		var res = rds.Receive("parser_to_transformer")

		chatId := res["chat_id"]

		chatIdInt, _ := strconv.Atoi(chatId)

		if len(res) > 0 {
			fmt.Println("Получаем сообщение, обрабатываем и отсылаем")

			toQueen := make(map[string]string)

			toQueen["chat_id"] = chatId

			for position, fileID := range res {

				fmt.Println("Key:", position, "=>", "Element:", fileID)
				fmt.Println("do something")

				fmt.Println("Получаем файл id")
				filePath := bot.GetFilePath(fileID)

				fmt.Println("получаем изображение")
				image := bot.GetImage(filePath)

				// Пока заглушка
				fmt.Println("Отправляем изображение в модель")
				d := ganserv.SendImageToModel(image, "version 1 (🔺 stylization, 🔻 robustness)")

				var dataFromTlg MessageTypes.RespDataTlg
				if plugins.IsZeroOfUnderlyingType(d) {
					text := "Упс, с датацентром что то не так, повторите попытку чуть позже. Мы уже занимаемся решением этой проблемы!"

					dataFromTlg = bot.SendMessage(chatIdInt, text)
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
										dataFromTlg = bot.SendMessage(chatIdInt, text)
									} else {
										bot.EditMessage(chatIdInt, text, int(dataFromTlg.Result.MessageId))
									}

									i++
								}
							} else {
								if data.Status == "COMPLETE" {
									//fmt.Println("Отправляем пользователю сообщение с фотографией")
									// Отправляем пользователю сообщение с фотографией
									imageString := strings.Split(data.Data.Data[0], ",")[1]
									f := bot.SendPhoto(chatIdInt, imageString)
									toQueen[position] = f
									// TODO сделать запись в базу
									//pg.InsertCancelAction(t.UserID)
									break
								}
							}
						} else {
							text := "Что то пошло не так:( Попробуйте загрузить другое фото!"
							dataFromTlg = bot.SendMessage(chatIdInt, text)
							break
						}
					}
				}

			}

			// тут код обработчика

			fmt.Println("отсылаем")

			rds.Send("transformer_to_creator", toQueen)
		} else {

			fmt.Println("nothing found")

		}

	}

}
