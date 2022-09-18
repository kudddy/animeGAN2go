package structure

type SendDataToPush struct {
	FnIndex     int      `json:"fn_index"`
	Data        []string `json:"data"`
	Action      string   `json:"action"`
	SessionHash string   `json:"session_hash"`
}

type SendDataStatus struct {
	Hash string `json:"hash"`
}
