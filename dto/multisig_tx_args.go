package dto

type MultisigTxArgs struct {
	Alias        string `json:"alias" binding:"required"`
	UnsignedTx   string `json:"unsignedTx" binding:"required"`
	Signature    string `json:"signature" binding:"required"`
	OutputOwners string `json:"outputOwners" binding:"required"`
}

type SignTxArgs struct {
	Signature string `json:"signature" binding:"required"`
}

type CompleteTxArgs struct {
	TransactionId string `json:"transactionId"`
	Signature     string `json:"signature"`
	Timestamp     string `json:"timestamp"`
}
