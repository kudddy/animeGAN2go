package ganserv

import (
	"animeGAN2go/MessageTypes"
	"animeGAN2go/structure"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func GetQueenNumber(hash string) (MessageTypes.CheckStatus, MessageTypes.CheckStatusQueen, bool, bool) {

	d := &structure.SendDataStatus{Hash: hash}

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

func SendImageToModel(sEncPhoto string, userModel string) MessageTypes.GetModelHash {

	var arr []string
	arr = append(arr, "data:image/jpeg;base64,"+sEncPhoto)

	fmt.Println(userModel)

	arr = append(arr, userModel)

	d := &structure.SendDataToPush{Data: arr, Action: "predict", FnIndex: 0, SessionHash: "q0tcnstipjf"}

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
