package Job

import (
	"animeGAN2go/bot"
	"animeGAN2go/ganserv"
	"animeGAN2go/plugins"
	"animeGAN2go/plugins/pg"
	"animeGAN2go/rds"
	"animeGAN2go/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// задача получить из очереди сообщения с файлами, обработать их и передать пути к файлам к следующему воркеру
//

func StartSingleWorker() {
	fmt.Println("get file ids from queen")

	for {

		time.Sleep(1 * time.Second)

		var res = rds.Receive("parser_to_transformer")

		if len(res) > 0 {

			chatId := res["chat_id"]
			userId := res["user_id"]
			userModel := res["user_model"]

			fmt.Printf("user is is %s\n", userId)

			chatIdInt, _ := strconv.Atoi(chatId)

			userIdInt, _ := strconv.Atoi(userId)

			fmt.Printf("user id after convert is %d\n", userIdInt)

			fmt.Println("Получаем сообщение, обрабатываем и отсылаем")

			toQueen := make(map[string]string)

			toQueen["chat_id"] = chatId

			var bar utils.Bar
			bar.NewOptionWithGraph(0, int64(len(res)), "#")

			progresResp := bot.SendMessage(chatIdInt, "Прогресс: ")

			j := 0

			for position, fileID := range res {

				fmt.Println("Key:", position, "=>", "Element:", fileID)
				fmt.Println("do something")

				fmt.Println("Получаем файл id")
				filePath := bot.GetFilePath(fileID)

				fmt.Println("получаем изображение")
				image := bot.GetImage(filePath)

				// Пока заглушка
				fmt.Println("Отправляем изображение в модель")
				d := ganserv.SendImageToModel(image, userModel)

				//var dataFromTlg MessageTypes.RespDataTlg
				if plugins.IsZeroOfUnderlyingType(d) {
					//var dataFromTlg MessageTypes.RespDataTlg
					text := "Упс, с датацентром что то не так, повторите попытку чуть позже. Мы уже занимаемся решением этой проблемы!"

					bot.SendMessage(chatIdInt, text)
					j++
				} else {
					i := 0
					for {
						// TODO make faster
						time.Sleep(500 * time.Millisecond)

						data, dataQueen, queen, err := ganserv.GetQueenNumber(d.Hash)
						if !err {
							if queen {
								if dataQueen.Status == "QUEUED" {
									//text := fmt.Sprintf("Ваша очередь: %s", strconv.Itoa(dataQueen.Data))
									if i == 0 {
										//dataFromTlg = bot.SendMessage(chatIdInt, text)
									} else {
										//bot.EditMessage(chatIdInt, text, int(dataFromTlg.Result.MessageId))
									}

									i++
								}
							} else {
								if data.Status == "COMPLETE" {
									// bot send message with status
									bot.EditMessage(int(progresResp.Result.Chat.Id),
										bar.Play(int64(j)),
										int(progresResp.Result.MessageId))

									// Отправляем пользователю сообщение с фотографией
									imageString := strings.Split(data.Data.Data[0], ",")[1]
									f := bot.SendPhoto(710828013, imageString)
									toQueen[position] = f
									j++
									break
								}
							}
						} else {
							//text := "Что то пошло не так:( Попробуйте загрузить другое фото!"
							//dataFromTlg = bot.SendMessage(chatIdInt, text)
							j++
							break
						}
					}
				}
			}

			rds.Send("transformer_to_creator", toQueen)

			bar.Finish()

			// delete from base info about busy worker
			// TODO this line should be in anime-gan-worker-creator
			pg.InsertCancelAction(userIdInt)
		} else {

			fmt.Println("nothing found")

		}

	}

}
