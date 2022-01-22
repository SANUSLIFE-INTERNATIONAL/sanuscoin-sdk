package rpc

import (
	"fmt"

	"github.com/btcsuite/btcutil"
)

type NewAddressRequest struct {
	Account string `json:"account"`
}

type NewAddressResponse struct {
	Address string `json:"address"`
}

func (wallet *Wallet) NewAddress(r NewAddressRequest, resp *NewAddressResponse) (err error) {
	address, err := wallet.wallet.NewAddress()
	if err != nil {
		return fmt.Errorf("error caused when trying to generate new address %v", err)

	}
	resp.Address = address.EncodeAddress()
	return
}

type ImportAddressRequest struct {
	PublicKey string `json:"publicKey"`
}

func (wallet *Wallet) ImportAddress(r ImportAddressRequest, resp *NewAddressResponse) (err error) {
	address, err := wallet.wallet.ImportAddress(r.PublicKey)
	if err != nil {
		return fmt.Errorf("error caused when trying to import address by private key %v", err)
	}
	resp.Address = address.EncodeAddress()
	return
}

func (wallet *Wallet) List(r ImportAddressRequest, resp *NewAddressResponse) (err error) {
	address, err := wallet.wallet.ImportAddress(r.PublicKey)
	if err != nil {
		return fmt.Errorf("error caused when trying to import address by private key %v", err)
	}
	resp.Address = address.EncodeAddress()
	return
}

type BalanceRequest struct {
	Address string `json:"address"`
}

type BalanceResponse struct {
	BTC float64 `json:"balance"`
	SNC int     `json:"snc"`
}

func (wallet *Wallet) Balance(r BalanceRequest, resp *BalanceResponse) (err error) {
	addr, err := btcutil.DecodeAddress(r.Address, wallet.wallet.GetNetParams())
	if err != nil {
		return err
	}
	resp.BTC, resp.SNC, err = wallet.wallet.Balance(addr)
	return err
}
