package sdk

import (
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/waddrmgr"
)

func (w *BTCWallet) NewAddress(account string) (btcutil.Address, error) {
	idx, err := w.wlt.NextAccount(waddrmgr.KeyScopeBIP0044, account)
	if err != nil {
		return nil, err
	}
	addr, err := w.wlt.NewAddress(idx, waddrmgr.KeyScopeBIP0044)
	if err != nil {
		return nil, err
	}
	return addr, nil
}
