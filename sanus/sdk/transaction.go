package sdk

import (
	"fmt"

	"sanus/sanus-sdk/cc/transfer"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/coinset"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/btcsuite/btcutil/txsort"
	"github.com/btcsuite/btcwallet/wallet/txrules"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
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

func (w *BTCWallet) SendTx(addressTo, addressFrom btcutil.Address, amountTarget btcutil.Amount, pkScript []byte) (string, error) {
	tx, err := w.buildTx(addressTo, addressFrom, amountTarget, 1, pkScript)
	if err != nil {
		return "", err
	}
	hash, err := w.wlt.ChainClient().SendRawTransaction(tx, false)
	if err != nil {
		return "", err
	}
	return hash.String(), nil
}

func (w *BTCWallet) buildTx(
	addressTo btcutil.Address,
	addressFrom btcutil.Address,
	amountTarget btcutil.Amount,
	feeLevel FeeLevel,
	pkScript []byte,
) (*wire.MsgTx, error) {
	amountIn, txIns, keysByAddrs, prevScripts, err := w.fetchUnspent(amountTarget, addressFrom)
	if err != nil {
		return nil, err
	}

	msgTx := wire.NewMsgTx(1)
	// Build the target output
	script, err := txscript.PayToAddrScript(addressTo)
	if err != nil {
		return nil, err
	}
	msgTx.AddTxOut(wire.NewTxOut(int64(amountTarget), script))

	txOutFee := msgTx.TxOut
	if pkScript != nil && len(pkScript) > 0 {
		msgTx.AddTxOut(wire.NewTxOut(0, pkScript))
		txOutFee = append(txOutFee, wire.NewTxOut(0, pkScript))
	}

	amountChange := amountIn - amountTarget
	if amountChange > 0 {
		// Build the change output
		script, err = txscript.PayToAddrScript(addressFrom)
		if err != nil {
			return nil, err
		}
		txOutFee = append(txOutFee, wire.NewTxOut(int64(amountChange), script))
	}

	targetFee := txrules.FeeForSerializeSize(
		w.estimateFee(feeLevel),
		txsizes.EstimateSerializeSize(len(txIns), txOutFee, true))

	// Check for dust output
	if txrules.IsDustAmount(amountTarget-targetFee, len(script), txrules.DefaultRelayFeePerKb) {
		return nil, ErrorDustAmount
	}

	amountChange = amountIn - amountTarget - targetFee
	if amountChange < 0 {
		return nil, fmt.Errorf("minimal estimated fee for this payment: %s", targetFee.String())
	} else if amountChange > 0 && txrules.IsDustAmount(amountChange, len(script), txrules.DefaultRelayFeePerKb) {
		return nil, ErrorDustAmountChange
	}

	keyClosure := txscript.KeyClosure(func(address btcutil.Address) (*btcec.PrivateKey, bool, error) {
		wif, found := keysByAddrs[address.EncodeAddress()]
		if !found {
			return nil, false, ErrKeyNotFound
		}

		return wif.PrivKey, wif.CompressPubKey, nil
	})

	scriptClosure := txscript.ScriptClosure(func(addr btcutil.Address) ([]byte, error) {
		return []byte{}, nil
	})

	// push TxIns into TX message
	msgTx.TxIn = txIns

	// push TxOuts of change amount into TX message
	if amountChange > 0 {
		msgTx.AddTxOut(wire.NewTxOut(int64(amountChange), script))
	}

	txsort.InPlaceSort(msgTx)

	for idx, txIn := range msgTx.TxIn {
		outScript, found := prevScripts[txIn.PreviousOutPoint]
		if !found {
			return nil, ErrPrevOutScript
		}

		script, err := txscript.SignTxOutput(
			w.netParams,
			msgTx,
			idx,
			outScript,
			txscript.SigHashAll,
			keyClosure,
			scriptClosure,
			txIn.SignatureScript)
		if err != nil {
			return nil, ErrFailedToSignTx
		}

		txIn.SignatureScript = script
	}

	return msgTx, nil
}

func (w *BTCWallet) fetchUnspent(target btcutil.Amount, source btcutil.Address) (
	amountIn btcutil.Amount,
	txIns []*wire.TxIn,
	keysByAddrs map[string]*btcutil.WIF,
	prevScripts map[wire.OutPoint][]byte,
	err error) {

	activeNet := w.GetNetParams()

	coinSet, err := w.genCoinSet(source, target, 0)
	if err != nil {
		return amountIn, txIns, keysByAddrs, prevScripts, err
	}

	coins := make([]coinset.Coin, 0, len(coinSet))
	for coin := range coinSet {
		coins = append(coins, coin)
	}
	list := coinset.NewCoinSet(coins)

	target = btcutil.Amount(btcutil.MaxSatoshi)

	keysByAddrs = make(map[string]*btcutil.WIF)
	prevScripts = make(map[wire.OutPoint][]byte)

	for _, coin := range list.Coins() {
		outpoint := wire.NewOutPoint(coin.Hash(), coin.Index())
		prevScripts[*outpoint] = coin.PkScript()

		txIn := wire.NewTxIn(outpoint, []byte{}, [][]byte{})
		txIn.Sequence = 0 // Opt-in RBF so we can bump fees

		amountIn += coin.Value()
		txIns = append(txIns, txIn)

		addr, err := coinSet[coin].Address(activeNet)
		if err != nil {
			continue
		}

		privateKey, err := coinSet[coin].ECPrivKey()
		if err != nil {
			continue
		}

		wif, err := btcutil.NewWIF(privateKey, activeNet, true)
		if err != nil {
			continue
		}

		keysByAddrs[addr.EncodeAddress()] = wif
	}

	return amountIn, txIns, keysByAddrs, prevScripts, nil
}

func (w *BTCWallet) genCoinSet(source btcutil.Address, targetBTC btcutil.Amount, targetSNC int) (map[coinset.Coin]*hdkeychain.ExtendedKey, error) {
	coinSet := make(map[coinset.Coin]*hdkeychain.ExtendedKey)

	var (
		pBTC = btcutil.Amount(0)
		pSNC = 0
	)

	unspent, err := w.wlt.ListUnspent(3, 9999999, map[string]struct{}{
		source.String(): {},
	})
	if err != nil {
		return coinSet, err
	}

	params := w.GetNetParams()

	for _, u := range unspent {
		if !u.Spendable {
			continue
		}

		amount, err := btcutil.NewAmount(u.Amount)
		if err != nil {
			return coinSet, err
		}

		address, err := btcutil.DecodeAddress(u.Address, params)
		if err != nil {
			return coinSet, err
		}

		privateKeyAddr, err := w.wlt.PrivKeyForAddress(address)
		if err != nil {
			return coinSet, err
		}

		scriptPubKey, err := txscript.PayToAddrScript(address)
		if err != nil {
			return coinSet, err
		}

		txHash := &chainhash.Hash{}
		if err := chainhash.Decode(txHash, u.TxID); err != nil {
			return coinSet, err
		}

		if scriptPubKey[0] == txscript.OP_RETURN {
			if isSatisfiedSNC(pSNC, targetSNC) {
				continue
			}

			pkScriptData, err := transfer.Decode(scriptPubKey)
			if err != nil {
				w.Logger.Errorf("error caused when trying to decode OP_RETURN script")
				continue
			}
			for _, p := range pkScriptData.Payments {
				pSNC = pSNC + p.Amount
			}
		} else {
			if isSatisfiedBTC(pBTC, targetBTC) {
				continue
			}
			pBTC = pBTC + amount
		}

		coin := &coinBase{
			TxHash:       txHash,
			TxIndex:      u.Vout,
			TxValue:      amount,
			TxNumConfs:   u.Confirmations,
			ScriptPubKey: scriptPubKey,
		}

		coinSet[coin] = hdkeychain.NewExtendedKey(
			params.HDPrivateKeyID[:],
			privateKeyAddr.Serialize(),
			make([]byte, 32),
			[]byte{0x00, 0x00, 0x00, 0x00},
			0,
			0,
			true)
	}

	return coinSet, nil
}

func isSatisfiedSNC(input, target int) bool {
	return input >= target
}

func isSatisfiedBTC(input, target btcutil.Amount) bool {
	return input >= target
}

func (w *BTCWallet) rescan(addr btcutil.Address) error {
	go func() {
		w.Logger.Infof("%v addr rescan has been started \n", addr)
		addrs := []btcutil.Address{addr}
		if err := w.wlt.ChainClient().Rescan(w.wlt.ChainParams().GenesisHash, addrs, nil); err != nil {
			w.Logger.Errorf("error caused when trying to rescan %v | %v \n", addr, err)
		}
		w.Logger.Errorf("rescan has been finished for %v\n", addr)
	}()
	return nil
}
