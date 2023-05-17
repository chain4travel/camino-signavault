/*
 * Copyright (C) 2023, Chain4Travel AG. All rights reserved.
 * See the file LICENSE for licensing terms.
 */

package txs

import (
	"github.com/chain4travel/camino-signavault/caminogo/vms/secp256k1fx"
)

// UnsignedTx is an unsigned transaction
type UnsignedTx interface {
	secp256k1fx.UnsignedTx
}
