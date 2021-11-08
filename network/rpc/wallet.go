package rpc

import (
	"encoding/hex"

	"sanus/sanus-sdk/sanus/sdk"
)

type Wallet struct {
	wallet *sdk.BTCWallet
}

func NewWalletHandler(wallet *sdk.BTCWallet) *Wallet {
	return &Wallet{wallet: wallet}
}

type BoolResponse struct {
	Success bool `json:"success"`
}

type OpenWalletRequest struct {
	Password string `json:"password"`
}

func (wallet *Wallet) Open(r OpenWalletRequest, resp *BoolResponse) (err error) {
	if err := wallet.wallet.Open([]byte(r.Password)); err != nil {
		resp.Success = false
		return err
	}
	resp.Success = true
	return
}

type CreateWalletRequest struct {
	Password string `json:"password"`
	Seed     string `json:"seed"`
}

func (wallet *Wallet) Create(r CreateWalletRequest, resp *BoolResponse) (err error) {
	seedHex, err := hex.DecodeString(r.Seed)
	if err != nil {
		resp.Success = false
		return err
	}
	if err := wallet.wallet.Create([]byte(r.Password), []byte(r.Password), seedHex); err != nil {
		resp.Success = false
		return err
	}
	resp.Success = true
	return
}

type UnlockWalletRequest struct {
	Password string `json:"password"`
}

func (wallet *Wallet) Unlock(r UnlockWalletRequest, resp *BoolResponse) (err error) {
	if err := wallet.wallet.Unlock([]byte(r.Password)); err != nil {
		resp.Success = false
		return err
	}
	resp.Success = true
	return
}

func (wallet *Wallet) Lock(r interface{}, resp *BoolResponse) (err error) {
	wallet.wallet.Lock()
	resp.Success = true
	return
}

func (wallet *Wallet) Synced(r interface{}, resp *BoolResponse) (err error) {
	resp.Success = wallet.wallet.Synced()
	return
}
