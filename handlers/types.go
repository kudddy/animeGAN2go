package handlers

/**
| ============== Types ============== |
*/
type Chat struct {
	Id    int    `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
}

type User struct {
	Id        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	UserName  string `json:"username"`
}

type Message struct {
	Id   int    `json:"message_id"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
	User User   `json:"from"`
}

type UpdateType struct {
	UpdateId      int     `json:"update_id"`
	Message       Message `json:"message"`
	EditedMessage Message `json:"edited_message"`
}

type OutMessage struct {
	ChatId        int    `json:"chat_id"`
	Text          string `json:"text"`
	ReplayToMsgId int    `json:"reply_to_message_id,omitempty"`
}

/**
| ============== Types for serv answer ============== |
*/

type RespByServ struct {
	Ok bool `json:"ok"`
}

/**
| ============== Types for send data to sm ============== |
*/

type Uuid struct {
	UserId      string
	Sub         string
	UserChannel string
}

type grammemInfo struct {
	Aspect       string `json:"aspect"`
	Mood         string `json:"mood"`
	Number       string `json:"number"`
	Person       string `json:"person"`
	Tense        string `json:"tense"`
	Transitivity string `json:"transitivity"`
	Verbform     string `json:"verbform"`
	Voice        string `json:"voice"`
	RawGramInfo  string `json:"raw_gram_info"`
	PartOfSpeech string `json:"part_of_speech"`
}

type tokenizedElements struct {
	Text             string
	RawText          string      `json:"raw_text"`
	GrammemInfo      grammemInfo `json:"grammem_info"`
	Lemma            string      `json:"lemma"`
	IsStopWord       bool        `json:"is_stop_word"`
	ListOfDependents []int       `json:"list_of_dependents"`
	DependencyType   string      `json:"dependency_type"`
	Head             int
}

type message struct {
	OriginalText                    string              `json:"original_text"`
	NormalizedText                  string              `json:"normalized_text"`
	TokenizedElementsList           []tokenizedElements `json:"tokenized_elements_list"`
	OriginalMessageName             string              `json:"original_message_name"`
	HumanNormalizedText             string              `json:"human_normalized_text"`
	HumanNormalizedTextWithAnaphora string              `json:"human_normalized_text_with_anaphora"`
}

type character struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Appeal string `json:"appeal"`
}

type appInfo struct {
	ProjectId       string `json:"projectId"`
	ApplicationId   string `json:"applicationId"`
	AppversionId    string `json:"appversionId"`
	FrontendType    string `json:"frontendType"`
	AgeLimit        int    `json:"ageLimit"`
	AffiliationType string `json:"affiliationType"`
}
type payload struct {
	Intent         string    `json:"intent"`
	OriginalIntent string    `json:"original_intent"`
	Msg            message   `json:"message"`
	NewSession     bool      `json:"new_session"`
	Character      character `json:"character"`
	ApplicationId  string    `json:"applicationId"`
	AppversionId   string    `json:"appversionId"`
	ProjectName    string    `json:"projectName"`
	AppInfo        appInfo   `json:"app_info"`
}

type ReqToSmType struct {
	MessageId   int     `json:"messageId"`
	SessionId   string  `json:"sessionId"`
	MessageName string  `json:"messageName"`
	Payload     payload `json:"payload"`
	Uuid        Uuid    `json:"uuid"`
}

/**
| ============== Types for answer from SmartMarket ============== |
*/

type strategies struct {
	LastCall string `json:"last_call"`
}

type actions struct {
	Type     string `json:"left"`
	Text     string `json:"text"`
	DeepLink string `json:"deep_link"`
}

type padding struct {
	Left   string `json:"left"`
	Top    string `json:"top"`
	Right  string `json:"right"`
	Bottom string `json:"bottom"`
}

type margins struct {
	Left   string `json:"left"`
	Top    string `json:"top"`
	Right  string `json:"right"`
	Bottom string `json:"bottom"`
}

type content struct {
	Url         string  `json:"url"`
	Hash        string  `json:"hash"`
	Width       string  `json:"width"`
	AspectRatio int     `json:"aspect_ratio"`
	Text        string  `json:"text"`
	Typeface    string  `json:"typeface"`
	TextColor   string  `json:"text_color"`
	MaxLines    int     `json:"max_lines"`
	Style       string  `json:"default"`
	Actions     actions `json:"actions"`
	Margins     margins `json:"margins"`
}

type cell struct {
	Type    string  `json:"type"`
	Content content `json:"content"`
	Padding padding `json:"padding"`
}

type card struct {
	Type  string `json:"type"`
	Cells []cell `json:"cells"`
}

type payloadForSm struct {
	OriginalIntent   string        `json:"original_intent"`
	IntentMeta       interface{}   `json:"intent_meta"`
	SelectedItem     interface{}   `json:"selected_item"`
	Strategies       strategies    `json:"strategies"`
	Asr              interface{}   `json:"asr"`
	ReverseGeocoding interface{}   `json:"reverseGeocoding"`
	BackInfo         []interface{} `json:"backInfo"`
	ApplicationId    string        `json:"applicationId"`
	AppversionId     string        `json:"appversionId"`
	ProjectName      string        `json:"projectName"`
	AppInfo          appInfo       `json:"app_info"`
	PronounceText    string        `json:"pronounceText"`
	Emotion          interface{}   `json:"emotion"`
	Items            []card        `json:"items"`
	AutoListening    bool          `json:"auto_listening"`
}

type RespFromSmType struct {
	MessageId   int          `json:"messageId"`
	SessionId   string       `json:"sessionId"`
	MessageName string       `json:"messageName"`
	Uuid        Uuid         `json:"uuid"`
	Payload     payloadForSm `json:"payload"`
}

func generatePayloadForSm(text string, sessionId string, messageId int) ReqToSmType {

	messageName := "MESSAGE_TO_SKILL"

	userUuid := Uuid{
		UserId:      "9485D45E-466E-4852-B5DA-1A27DFF5EFC8",
		Sub:         "1hkmItxUo6BDBmNvGM7inj4kNvWIRyQOaUzWdlqxYafPUqNZ/fTLMJ8M4idi1y467byHIwH8zAnbqt6glUevV0d8+tppO2Ysr1Ryn5PPj7nkk+7kTtDC1MnJvZVaJP3uzHxG5PPxvQpIbtQccKxegw==",
		UserChannel: "SBOL",
	}

	appInfo := appInfo{
		ProjectId:       "12f20e40-efc6-4ff5-9179-f5c51f7197b3",
		ApplicationId:   "7aa5ae84-c668-4e24-94d8-e35cf053e7a1",
		AppversionId:    "bbddbed8-a8c6-483f-99b5-516dbae4ea70",
		FrontendType:    "DIALOG",
		AgeLimit:        18,
		AffiliationType: "ECOSYSTEM",
	}

	message := message{
		OriginalText:                    text,
		NormalizedText:                  text,
		OriginalMessageName:             "MESSAGE_FROM_USER",
		HumanNormalizedText:             text,
		HumanNormalizedTextWithAnaphora: text,
	}

	payload := payload{
		Intent:         "sberauto_main",
		OriginalIntent: "food",
		NewSession:     false,
		ApplicationId:  "7aa5ae84-c668-4e24-94d8-e35cf053e7a1",
		AppversionId:   "bbddbed8-a8c6-483f-99b5-516dbae4ea70",
		ProjectName:    "СберАвто. Подбор автомобиля",
		AppInfo:        appInfo,
		Msg:            message,
	}

	reqToSmType := ReqToSmType{
		MessageId:   messageId,
		SessionId:   sessionId,
		MessageName: messageName,
		Payload:     payload,
		Uuid:        userUuid,
	}

	return reqToSmType

}
