package asset

import (
	"encoding/json"

	"github.com/btcsuite/btcd/wire"
)

type CCVout struct {
	Out    *wire.TxOut    `json:"out"`
	Assets map[int]*Asset `json:"assets"`
}

func (a *CCVout) String() string {
	bytes, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (vout *CCVout) Output() *wire.TxOut {
	return vout.Out
}

func (vout *CCVout) SetOutput(out *wire.TxOut) {
	vout.Out = out
}

func (vout *CCVout) GetAssets() map[int]*Asset {
	return vout.Assets
}

func (vout *CCVout) AppendAsset(asset *Asset) {
	vout.Assets[len(vout.Assets)] = asset
}

func (vout *CCVout) SetAssetsArray(assets map[int]*Asset) {
	vout.Assets = assets
}
