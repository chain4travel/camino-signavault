// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"github.com/chain4travel/camino-signavault/caminogo/ids"
)

// BaseTx contains fields common to many transaction types. It should be
// embedded in transaction implementations.
type BaseTx struct {
	NetworkID     uint32 `serialize:"true" json:"networkID"` // ID of the network this chain lives on
	BlockchainID  ids.ID `serialize:"true" json:"blockchainID"`
	unsignedBytes []byte // Unsigned byte representation of this data
}

func (tx *BaseTx) Initialize(unsignedBytes []byte) {
	tx.unsignedBytes = unsignedBytes
}

func (tx *BaseTx) Bytes() []byte {
	return tx.unsignedBytes
}
