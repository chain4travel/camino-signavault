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

type IssueTxArgs struct {
	SignedTx  string `json:"signedTx" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type IssueTxResponse struct {
	TxID string `json:"txID"`
}
