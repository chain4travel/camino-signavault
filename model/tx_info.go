/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package model

type TxInfo struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Tx       string `json:"tx"`
		Encoding string `json:"encoding"`
	} `json:"result"`
	Id int `json:"id"`
}
