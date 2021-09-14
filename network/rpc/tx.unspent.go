package rpc

import (
	"fmt"

	"sanus/sanus-sdk/sanus/daemon"

	"github.com/btcsuite/btcutil"
)

type UnspentTxRequest struct {
	Address string
}

type UnspentTxResponse struct {
	List []string `json:"list"`
}

func (tx *Tx) Unspent(r UnspentTxRequest, resp *UnspentTxResponse) (err error) {
	if r.Address == "" {
		return fmt.Errorf("address can't be empty")
	}
	address, err := btcutil.DecodeAddress(r.Address, daemon.ActiveNetParams.Params)
	if err != nil {
		return err
	}
	resp.List, err = tx.wallet.UnspentTx(address)
	return
}
