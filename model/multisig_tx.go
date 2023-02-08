package model

type MultisigTx struct {
	Id            int64              `json:"id"	 		`
	Alias         string             `json:"alias" 		`
	Threshold     int                `json:"threshold" 	`
	Signers       []MultisigTxSigner `json:"signers" 		`
	TransactionId string             `json:"transactionId" "`
	UnsignedTx    string             `json:"unsignedTx" 	`
}

type MultisigTxSigner struct {
	Id           int64
	MultisigTxId int64
	Address      string `json:"address"`
	Signature    string `json:"signature"`
}
