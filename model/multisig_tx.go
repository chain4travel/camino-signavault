package model

type MultisigTx struct {
	Id            string            `json:"id"`
	UnsignedTx    string            `json:"unsignedTx"`
	Alias         string            `json:"alias"`
	Threshold     int8              `json:"threshold"`
	TransactionId string            `json:"transactionId"`
	OutputOwners  string            `json:"outputOwners"`
	Owners        []MultisigTxOwner `json:"owners"`
}

type MultisigTxOwner struct {
	MultisigTxId string `json:"-"`
	Address      string `json:"address"`
	Signature    string `json:"signature"`
}
