package sdk

import (
	"github.com/btcsuite/btcwallet/waddrmgr"
)

type NetworkData struct {
	SyncedTo waddrmgr.BlockStamp
}

func (w *BTCWallet) NetworkStatus() *NetworkData {
	return &NetworkData{SyncedTo: w.wlt.Manager.SyncedTo()}
}
