package rpc

import (
	"sanus/sanus-sdk/sanus/sdk"
)

type Network struct {
	wallet *sdk.BTCWallet
}

func NewNetworkHandler(wallet *sdk.BTCWallet) *Network {
	return &Network{wallet: wallet}
}

func (ntw *Network) Status(r interface{}, resp *sdk.NetworkData) (err error) {
	resp = ntw.wallet.NetworkStatus()
	return
}
