package sdk

import (
	"github.com/btcsuite/btcutil"
)

func (w *BTCWallet) UnspentTx(addr btcutil.Address) ([]string, error) {
	list, err := w.wlt.ListUnspent(3, 99999, map[string]struct{}{
		addr.EncodeAddress(): {},
	})
	if err != nil {
		return nil, err
	}
	var txs = make([]string, len(list), len(list))
	for k, tx := range list {
		txs[k] = tx.TxID
	}
	return txs, err
}
