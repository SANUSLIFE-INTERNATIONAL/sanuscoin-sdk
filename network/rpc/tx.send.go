package rpc

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcutil"
)

type SendTxRequest struct {
	To       string  `json:"to"`
	From     string  `json:"from"`
	Amount   float64 `json:"amount"`
	PkScript string  `json:"pk_script"`
}

type SendTxResponse struct {
	Hash string `json:"hash"`
}

func (tx *Tx) Send(r SendTxRequest, resp *SendTxResponse) (err error) {
	addressTo, err := btcutil.DecodeAddress(r.To, tx.wallet.GetNetParams())
	if err != nil {
		return err
	}
	addressFrom, err := btcutil.DecodeAddress(r.From, tx.wallet.GetNetParams())
	if err != nil {
		return err
	}
	if r.Amount <= 0 {
		return fmt.Errorf("amount can't be negative or equal 0")

	}
	amountReal, err := btcutil.NewAmount(r.Amount)
	if err != nil {
		return err
	}
	var pkScriptByte []byte = nil
	pkScript := r.PkScript
	if pkScript != "" {
		pkScriptByte, err = hex.DecodeString(pkScript)
		if err != nil {
			return err
		}
	}
	resp.Hash, err = tx.wallet.SendTx(addressTo, addressFrom, amountReal, pkScriptByte)
	if err != nil {
		return fmt.Errorf("error caused when trying to send tx %v", err)
	}
	return
}
