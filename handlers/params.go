package handlers

// params for erkc adapter

var APIEndpoint = ""

var MessengerEntryPoint = map[string]string{
	"psi": "messenger-t.sberbank.ru:7764",
	"ift": "messenger-ift.sberbank.ru:7764",
}

var AuthServiceHost = map[string]map[string]string{
	"psi": {
		"sbl": "https://i-see-you.ru/sbl/psi/session/init",
		"sbd": "https://i-see-you.ru/devices/psi/session/init",
	},
	"ift": {
		"npf": "https://i-see-you.ru/npf/ift/session/init",
		"smm": "https://i-see-you.ru/smm/ift/session/init",
	},
}
