package structure

type SendDataToPush struct {
	Data   []string `json:"data"`
	Action string   `json:"action"`
}

type SendDataStatus struct {
	Hash string `json:"hash"`
}
