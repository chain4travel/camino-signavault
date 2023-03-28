/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package model

type MultisigTx struct {
	Id            string            `json:"id"`
	UnsignedTx    string            `json:"unsignedTx"`
	Alias         string            `json:"alias"`
	Threshold     int8              `json:"threshold"`
	TransactionId string            `json:"transactionId"`
	OutputOwners  string            `json:"outputOwners"`
	Metadata      string            `json:"metadata"`
	Owners        []MultisigTxOwner `json:"owners"`
}

type MultisigTxOwner struct {
	MultisigTxId string `json:"-"`
	Address      string `json:"address"`
	Signature    string `json:"signature"`
}
