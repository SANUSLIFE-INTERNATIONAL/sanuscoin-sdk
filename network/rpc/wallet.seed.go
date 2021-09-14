package rpc

import (
	"encoding/hex"

	"sanus/sanus-sdk/sanus/sdk"
)

type GenerateSeedRequest struct {
	Mnemonic string `json:"mnemonic"`
	Public   string `json:"public"`
}

type GenerateSeedResponse struct {
	Hash string `json:"hash"`
}

func (*Wallet) Seed(r GenerateSeedRequest, res *GenerateSeedResponse) (err error) {
	seed := sdk.NewSeed(r.Mnemonic, r.Public)
	res.Hash = hex.EncodeToString(seed)
	return
}
