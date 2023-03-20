package dto

type MultisigTxArgs struct {
	Alias      string `json:"alias"`
	UnsignedTx string `json:"unsignedTx"`
	Signature  string `json:"signature"`
}

type SignTxArgs struct {
	Signature string `json:"signature"`
}

type CompleteTxArgs struct {
	TransactionId string `json:"transactionId"`
	Signature     string `json:"signature"`
	Timestamp     string `json:"timestamp"`
}
