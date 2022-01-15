package asset

import (
	"encoding/json"

	"github.com/btcsuite/btcd/wire"
)

type CCVin struct {
	Input  *wire.TxIn     `json:"input"`
	Assets map[int]*Asset `json:"assets"`
}

func (a *CCVin) String() string {
	bytes, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (vin *CCVin) GetInput() *wire.TxIn {
	return vin.Input
}

func (vin *CCVin) GetAssets() map[int]*Asset {
	return vin.Assets
}

func (vin *CCVin) AppendAsset(asset *Asset) {
	vin.Assets[len(vin.Assets)] = asset
}

func (vin *CCVin) SetInput(input *wire.TxIn) {
	vin.Input = input
}
