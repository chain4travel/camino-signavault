package model

type MultisigTx struct {
	Alias         string  `json:"alias"`
	Threshold     int     `json:"threshold"`
	Signers       []Owner `json:"signers"`
	TransactionId string  `json:"transactionId"`
	UnsignedTx    string  `json:"unsignedTx"`
}

type Owner struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
}
