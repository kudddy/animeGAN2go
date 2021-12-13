package MessageTypes

// ответ сервиса
type CheckTokenResp struct {
	MessageName string
	Status      bool
	Desc        string
	Token       string
}

/// ответ от телеграмма
type Result struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	FilePath     string `json:"file_path"`
}

type GetFilePath struct {
	Ok     bool   `json:"ok"`
	Result Result `json:"result"`
}

/// ответ от сервиса с моделью
type GetModelHash struct {
	Hash          string `json:"hash"`
	QueuePosition int    `json:"queue_position"`
}

// проверяем статус

type DoneData struct {
	AvgDurations []float32 `json:"avg_durations"`
	Data         []string  `json:"data"`
	Duration     []float32 `json:"duration"`
}
type CheckStatus struct {
	Data   DoneData `json:"data"`
	Status string   `json:"status"`
}

type CheckStatusQueen struct {
	Data   int    `json:"data"`
	Status string `json:"status"`
}

// получаем на вход
type ReqData struct {
	FileId    string `json:"file_id"`
	ChatId    int    `json:"chat_id"`
	UserID    int    `json:"user_id"`
	UserModel string `json:"user_model"`
}

// ответ от серверов телеграмм на отправку сообщения

type From struct {
	Id        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type Chat struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}
type ResultResp struct {
	MessageId int64  `json:"message_id"`
	Data      int64  `json:"date"`
	Text      string `json:"text"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
}
type RespDataTlg struct {
	Ok     bool       `json:"ok"`
	Result ResultResp `json:"result"`
}

type StartJobResp struct {
	MessageName string
	Status      bool
}

type CheckJobStatusResp struct {
	MessageName string
	Status      string
	Token       string
}
