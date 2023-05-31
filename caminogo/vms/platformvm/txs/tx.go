/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package txs

// Tx is a signed transaction
type Tx struct {
	// The body of this transaction
	Unsigned UnsignedTx `serialize:"true" json:"unsignedTx"`
	id       [32]byte
	bytes    []byte
}
