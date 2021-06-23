package sdk

import (
	"github.com/btcsuite/btcutil"
)

func (w *BTCWallet) UnspentTx(addr btcutil.Address) ([]string, error) {
	if err := w.rescan(addr); err != nil {
		return nil, err
	}
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

func (w *BTCWallet) rescan(addr btcutil.Address) error {
	nData := w.NetworkStatus()
	currentBlock := nData.SyncedTo.Height
	blockHash, err := w.wlt.ChainClient().GetBlockHash(int64(currentBlock - 500))
	if err != nil {
		return err
	}
	addrs := []btcutil.Address{addr}
	return w.wlt.ChainClient().Rescan(blockHash, addrs, nil)
}
