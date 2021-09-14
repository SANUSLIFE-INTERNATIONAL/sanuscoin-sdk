package rpc

import (
	"sanus/sanus-sdk/sanus/sdk"
)

type Tx struct {
	wallet *sdk.BTCWallet
}

func NewTxHandler(wallet *sdk.BTCWallet) *Tx {
	return &Tx{wallet: wallet}
}
