package handlers

type EditDataToTlg struct {
	Text      string `json:"text"`
	ChatId    int    `json:"chat_id"`
	MessageId int    `json:"message_id"`
}

type SendDataToTlg struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}
