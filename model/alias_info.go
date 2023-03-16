package model

type AliasInfo struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  Result `json:"result"`
	Id      int    `json:"id"`
}

type Result struct {
	Memo      string   `json:"memo"`
	Addresses []string `json:"addresses"`
	Threshold string   `json:"threshold"`
}
