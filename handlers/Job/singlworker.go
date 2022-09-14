package Job

import (
	"animeGAN2go/bot"
	"animeGAN2go/ganserv"
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
		userId := res["user_id"]

		chatIdInt, _ := strconv.Atoi(chatId)

		userIdInt, _ := strconv.Atoi(userId)

		if len(res) > 0 {
			fmt.Println("Получаем сообщение, обрабатываем и отсылаем")

			toQueen := make(map[string]string)

			toQueen["chat_id"] = chatId

			bot.SendMessage(chatIdInt, "Приступаем к трансформации изображения! Осталось совсем чуть-чуть!")

			for position, fileID := range res {

				fmt.Println("Key:", position, "=>", "Element:", fileID)
				fmt.Println("do something")

				fmt.Println("Получаем файл id")
				filePath := bot.GetFilePath(fileID)

				fmt.Println("получаем изображение")
				image := bot.GetImage(filePath)

				// Пока заглушка
				fmt.Println("Отправляем изображение в модель")
				d := ganserv.SendImageToModel(image, "version 2 (🔺 robustness,🔻 stylization)")

				var dataFromTlg MessageTypes.RespDataTlg
				if plugins.IsZeroOfUnderlyingType(d) {
					text := "Упс, с датацентром что то не так, повторите попытку чуть позже. Мы уже занимаемся решением этой проблемы!"

					dataFromTlg = bot.SendMessage(710828013, text)
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
										dataFromTlg = bot.SendMessage(710828013, text)
									} else {
										bot.EditMessage(710828013, text, int(dataFromTlg.Result.MessageId))
									}

									i++
								}
							} else {
								if data.Status == "COMPLETE" {
									//fmt.Println("Отправляем пользователю сообщение с фотографией")
									// Отправляем пользователю сообщение с фотографией
									imageString := strings.Split(data.Data.Data[0], ",")[1]
									// кидаем результат преобразования только в один чат
									f := bot.SendPhoto(710828013, imageString)
									toQueen[position] = f
									// TODO сделать запись в базу

									break
								}
							}
						} else {
							text := "Что то пошло не так:( Попробуйте загрузить другое фото!"
							dataFromTlg = bot.SendMessage(710828013, text)
							break
						}
					}
				}

			}

			// тут код обработчика

			fmt.Println("отсылаем")

			rds.Send("transformer_to_creator", toQueen)

			pg.InsertCancelAction(userIdInt)
		} else {

			fmt.Println("nothing found")

		}

	}

}
