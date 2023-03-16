package model

type TxInfo struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Tx       string `json:"tx"`
		Encoding string `json:"encoding"`
	} `json:"result"`
	Id int `json:"id"`
}
