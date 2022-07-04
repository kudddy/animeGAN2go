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

type KeyboardButtonPollType struct {
	// Type is if quiz is passed, the user will be allowed to create only polls
	// in the quiz mode. If regular is passed, only regular polls will be
	// allowed. Otherwise, the user will be allowed to create a poll of any type.
	Type string `json:"type"`
}

type KeyboardButton struct {
	// Text of the button. If none of the optional fields are used,
	// it will be sent as a message when the button is pressed.
	Text string `json:"text"`
	// RequestContact if True, the user's phone number will be sent
	// as a contact when the button is pressed.
	// Available in private chats only.
	//
	// optional
	RequestContact bool `json:"request_contact,omitempty"`
	// RequestLocation if True, the user's current location will be sent when
	// the button is pressed.
	// Available in private chats only.
	//
	// optional
	RequestLocation bool `json:"request_location,omitempty"`
	// RequestPoll if True, the user will be asked to create a poll and send it
	// to the bot when the button is pressed. Available in private chats only
	//
	// optional
	RequestPoll *KeyboardButtonPollType `json:"request_poll,omitempty"`
}

// ReplyKeyboardMarkup represents a custom keyboard with reply options.
type ReplyKeyboardMarkup struct {
	// Keyboard is an array of button rows, each represented by an Array of KeyboardButton objects
	Keyboard [][]KeyboardButton `json:"keyboard"`
	// ResizeKeyboard requests clients to resize the keyboard vertically for optimal fit
	// (e.g., make the keyboard smaller if there are just two rows of buttons).
	// Defaults to false, in which case the custom keyboard
	// is always of the same height as the app's standard keyboard.
	//
	// optional
	ResizeKeyboard bool `json:"resize_keyboard,omitempty"`
	// OneTimeKeyboard requests clients to hide the keyboard as soon as it's been used.
	// The keyboard will still be available, but clients will automatically display
	// the usual letter-keyboard in the chat – the user can press a special button
	// in the input field to see the custom keyboard again.
	// Defaults to false.
	//
	// optional
	OneTimeKeyboard bool `json:"one_time_keyboard,omitempty"`
	// InputFieldPlaceholder is the placeholder to be shown in the input field when
	// the keyboard is active; 1-64 characters.
	//
	// optional
	InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
	// Selective use this parameter if you want to show the keyboard to specific users only.
	// Targets:
	//  1) users that are @mentioned in the text of the Message object;
	//  2) if the bot's message is a reply (has Message.ReplyToMessage not nil), sender of the original message.
	//
	// Example: A user requests to change the bot's language,
	// bot replies to the request with a keyboard to select the new language.
	// Other users in the group don't see the keyboard.
	//
	// optional
	Selective bool `json:"selective,omitempty"`
}

type LoginURL struct {
	// URL is an HTTP URL to be opened with user authorization data added to the
	// query string when the button is pressed. If the user refuses to provide
	// authorization data, the original URL without information about the user
	// will be opened. The data added is the same as described in Receiving
	// authorization data.
	//
	// NOTE: You must always check the hash of the received data to verify the
	// authentication and the integrity of the data as described in Checking
	// authorization.
	URL string `json:"url"`
	// ForwardText is the new text of the button in forwarded messages
	//
	// optional
	ForwardText string `json:"forward_text,omitempty"`
	// BotUsername is the username of a bot, which will be used for user
	// authorization. See Setting up a bot for more details. If not specified,
	// the current bot's username will be assumed. The url's domain must be the
	// same as the domain linked with the bot. See Linking your domain to the
	// bot for more details.
	//
	// optional
	BotUsername string `json:"bot_username,omitempty"`
	// RequestWriteAccess if true requests permission for your bot to send
	// messages to the user
	//
	// optional
	RequestWriteAccess bool `json:"request_write_access,omitempty"`
}

// CallbackGame is for starting a game in an inline keyboard button.
type CallbackGame struct{}

type InlineKeyboardButton struct {
	// Text label text on the button
	Text string `json:"text"`
	// URL HTTP or tg:// url to be opened when button is pressed.
	//
	// optional
	URL *string `json:"url,omitempty"`
	// LoginURL is an HTTP URL used to automatically authorize the user. Can be
	// used as a replacement for the Telegram Login Widget
	//
	// optional
	LoginURL *LoginURL `json:"login_url,omitempty"`
	// CallbackData data to be sent in a callback query to the bot when button is pressed, 1-64 bytes.
	//
	// optional
	CallbackData string `json:"callback_data,omitempty"`
	// SwitchInlineQuery if set, pressing the button will prompt the user to select one of their chats,
	// open that chat and insert the bot's username and the specified inline query in the input field.
	// Can be empty, in which case just the bot's username will be inserted.
	//
	// This offers an easy way for users to start using your bot
	// in inline mode when they are currently in a private chat with it.
	// Especially useful when combined with switch_pm… actions – in this case
	// the user will be automatically returned to the chat they switched from,
	// skipping the chat selection screen.
	//
	// optional
	SwitchInlineQuery *string `json:"switch_inline_query,omitempty"`
	// SwitchInlineQueryCurrentChat if set, pressing the button will insert the bot's username
	// and the specified inline query in the current chat's input field.
	// Can be empty, in which case only the bot's username will be inserted.
	//
	// This offers a quick way for the user to open your bot in inline mode
	// in the same chat – good for selecting something from multiple options.
	//
	// optional
	SwitchInlineQueryCurrentChat *string `json:"switch_inline_query_current_chat,omitempty"`
	// CallbackGame description of the game that will be launched when the user presses the button.
	//
	// optional
	CallbackGame *CallbackGame `json:"callback_game,omitempty"`
	// Pay specify True, to send a Pay button.
	//
	// NOTE: This type of button must always be the first button in the first row.
	//
	// optional
	Pay *bool `json:"pay,omitempty"`
}

type InlineKeyboardMarkup struct {
	// InlineKeyboard array of button rows, each represented by an Array of
	// InlineKeyboardButton objects
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type OutMessage struct {
	ChatId        int                  `json:"chat_id"`
	Text          string               `json:"text"`
	ReplayToMsgId *int                 `json:"reply_to_message_id,omitempty"`
	ReplyMarkup   InlineKeyboardMarkup `json:"reply_markup,omitempty"`
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
