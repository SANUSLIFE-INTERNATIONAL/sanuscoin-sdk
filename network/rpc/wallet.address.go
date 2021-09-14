package rpc

import (
	"fmt"

	"sanus/sanus-sdk/misc/random"

	"github.com/btcsuite/btcutil"
)

type NewAddressRequest struct {
	Account string `json:"account"`
}

type NewAddressResponse struct {
	Address string `json:"address"`
}

func (wallet *Wallet) NewAddress(r NewAddressRequest, resp *NewAddressResponse) (err error) {
	if r.Account == "" {
		r.Account = random.RandStringRunes(16)
	}
	address, err := wallet.wallet.NewAddress(r.Account)
	if err != nil {
		return fmt.Errorf("error caused when trying to generate new address %v", err)

	}
	resp.Address = address.EncodeAddress()
	return
}

type BalanceRequest struct {
	Address string `json:"address"`
	Coin    string `json:"coin"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}

func (wallet *Wallet) Balance(r BalanceRequest, resp *BalanceResponse) (err error) {
	addr, err := btcutil.DecodeAddress(r.Address, wallet.wallet.GetNetParams())
	if err != nil {
		return err
	}
	switch r.Coin {
	case "btc":
		balance, err := wallet.wallet.BTCBalance(addr)
		resp.Balance = balance
		return err

	case "snc":
		balance, err := wallet.wallet.SNCBalance(addr)
		resp.Balance = float64(balance)
		return err
	default:
		return fmt.Errorf("invalid coin type")
	}
}
