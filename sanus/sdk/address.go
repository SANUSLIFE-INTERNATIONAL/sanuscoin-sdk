package sdk

import (
	"fmt"

	"sanus/sanus-sdk/cc/transfer"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/waddrmgr"
	btcWallet "github.com/btcsuite/btcwallet/wallet"
)

// NewAddress method generates a new BIP44 address
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

// ImportAddress method imports a new address based on privat key
func (w *BTCWallet) ImportAddress(privateKey string) (btcutil.Address, error) {
	wif, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, err
	}
	blockStamp := waddrmgr.BlockStamp{
		Height: 0,
		Hash:   *w.wlt.ChainParams().GenesisHash,
	}
	addr, err := w.wlt.ImportPrivateKey(waddrmgr.KeyScopeBIP0044, wif, &blockStamp, true)
	if err != nil {
		return nil, err
	}
	decodedAddr, err := btcutil.DecodeAddress(addr, w.netParams)
	if err != nil {
		return nil, err
	}
	w.rescan(decodedAddr)
	return decodedAddr, nil
}

// get current address for default (number 0) account,
// a new address will be generated if not exists yet
func (w *BTCWallet) defaultAddress() (address btcutil.Address, err error) {
	return w.wlt.CurrentAddress(0, waddrmgr.KeyScopeBIP0044)
}

func (w *BTCWallet) SNCBalance(address btcutil.Address) (int64, error) {
	txs, err := w.wlt.ListUnspent(3, 9999999, map[string]struct{}{
		address.EncodeAddress(): {},
	})

	if err != nil {
		return 0, err
	}

	var balance int64 = 0
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
				pkScriptData, err := transfer.Decode(pkScript)
				fmt.Printf("%#v", pkScriptData)
				fmt.Printf("%#v", pkScriptData.Payments[0])
				if err != nil {
					w.Errorf("Error caused when trying to fetch data from PkScript | %v", err)
				}

				for _, p := range pkScriptData.Payments {
					balance += int64(p.Amount)
				}
			}
		}
	}
	return balance, err
}

func (w *BTCWallet) BTCBalance(address btcutil.Address) (float64, error) {
	txs, err := w.wlt.ListUnspent(3, 9999999, map[string]struct{}{
		address.EncodeAddress(): {},
	})
	if err != nil {
		return 0, nil
	}
	var amount float64
	for _, tx := range txs {
		amount += tx.Amount
	}
	return amount, nil
}
