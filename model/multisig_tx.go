/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package model

import "time"

type MultisigTx struct {
	Id            string            `json:"id" binding:"required"`
	UnsignedTx    string            `json:"unsignedTx" binding:"required"`
	Alias         string            `json:"alias" binding:"required"`
	Threshold     int8              `json:"threshold" binding:"required"`
	TransactionId string            `json:"transactionId"`
	OutputOwners  string            `json:"outputOwners" binding:"required"`
	Metadata      string            `json:"metadata"`
	Owners        []MultisigTxOwner `json:"owners" binding:"required"`
	Timestamp     time.Time         `json:"timestamp"`
}

type MultisigTxOwner struct {
	MultisigTxId string `json:"-" binding:"required"`
	Address      string `json:"address" binding:"required"`
	Signature    string `json:"signature"`
}
