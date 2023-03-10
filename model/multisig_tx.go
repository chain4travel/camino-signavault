package model

type MultisigTx struct {
	Id            int64              `json:"id"`
	UnsignedTx    string             `json:"unsignedTx"`
	Alias         string             `json:"alias"`
	Threshold     int8               `json:"threshold"`
	TransactionId string             `json:"transactionId"`
	Owners        []MultisigTxOwner  `json:"owners"`
	Signers       []MultisigTxSigner `json:"signers"`
}

type MultisigTxOwner struct {
	Id           int64  `json:"-"`
	MultisigTxId int64  `json:"-"`
	Address      string `json:"address"`
}

type MultisigTxSigner struct {
	MultisigTxOwner
	Signature string `json:"signature"`
}
