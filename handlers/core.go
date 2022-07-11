package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func sendReqToSm(urlPath string, outData ReqToSmType) (RespFromSmType, error) {

	var data RespFromSmType

	body, err := json.Marshal(outData)
	if err != nil {
		log.Fatal(err)
		return data, err
	}
	bodyBytes := bytes.NewBuffer(body)
	response, err := http.Post(urlPath, "application/json", bodyBytes)

	if err != nil {
		return data, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Все ок, код положительный")
		decoder := json.NewDecoder(response.Body)

		err = decoder.Decode(&data)
		return data, nil

	} else {
		// TODO this logic not true
		return data, nil
	}
}

func sendReqToTlg(urlPath string, outData OutMessage) error {

	var data RespFromSmType

	body, err := json.Marshal(outData)
	if err != nil {
		log.Fatal(err)
		return err
	}
	bodyBytes := bytes.NewBuffer(body)

	response, err := http.Post(urlPath, "application/json", bodyBytes)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("Все ок, код положительный")
		decoder := json.NewDecoder(response.Body)

		err = decoder.Decode(&data)

		return nil
	} else {
		//  TODO it is bug
		return nil
	}
}
