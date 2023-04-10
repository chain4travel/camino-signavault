/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package dto

type MultisigTxArgs struct {
	Alias        string `json:"alias" binding:"required"`
	UnsignedTx   string `json:"unsignedTx" binding:"required"`
	Signature    string `json:"signature" binding:"required"`
	OutputOwners string `json:"outputOwners" binding:"required"`
	Metadata     string `json:"metadata"`
	Expiration   int64  `json:"expiration"`
	ChainId      string `json:"chainId" binding:"required"`
}

type SignTxArgs struct {
	Signature string `json:"signature" binding:"required"`
}

type IssueTxArgs struct {
	SignedTx  string `json:"signedTx" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type IssueTxResponse struct {
	TxID string `json:"txID" binding:"required"`
}
