package sdk

import (
	"fmt"

	"sanus/sanus-sdk/cc"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/waddrmgr"
	btcWallet "github.com/btcsuite/btcwallet/wallet"
)

func (w *BTCWallet) NewAddress(account string) (btcutil.Address, error) {
	idx, err := w.wlt.NextAccount(waddrmgr.KeyScopeBIP0044, account)
	if err != nil {
		return nil, err
	}
	addr, err := w.wlt.NewChangeAddress(idx, waddrmgr.KeyScopeBIP0044)
	if err != nil {
		return nil, err
	}
	return addr, nil
}

func (w *BTCWallet) SNCBalance(address btcutil.Address) (int, error) {
	txs, err := w.wlt.ListUnspent(3, 9999999, map[string]struct{}{
		address.EncodeAddress(): {},
	})

	if err != nil {
		return 0, err
	}

	balance := 0
	for _, tx := range txs {
		h, err := chainhash.NewHashFromStr(tx.TxID)
		if err != nil {
			return 0, err
		}
		txDetail, err := btcWallet.UnstableAPI(w.wlt).TxDetails(h)
		if err != nil {
			return 0, err
		}

		for _, out := range txDetail.MsgTx.TxOut {
			pkScript := out.PkScript
			if pkScript[0] == txscript.OP_RETURN {
				pkScriptData, err := cc.Decode(pkScript)
				if err != nil {
					w.Errorf("Error caused when trying to fetch data from PkScript | %v", err)
				}
				fmt.Printf("%+v", pkScriptData)
				for _, p := range pkScriptData.Payments {
					balance += p.Amount
				}
			}
		}
	}
	return balance, err
}
