package dto

type MultisigTxArgs struct {
	Alias      string `json:"alias"`
	UnsignedTx string `json:"unsignedTx"`
	Signature  string `json:"signature"`
}

type SignTxArgs struct {
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
}

type CompleteTxArgs struct {
	TransactionId string `json:"transactionId"`
}
