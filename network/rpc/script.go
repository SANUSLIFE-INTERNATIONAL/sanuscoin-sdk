package rpc

import (
	"encoding/hex"

	"sanus/sanus-sdk/cc/issuance"
	"sanus/sanus-sdk/cc/transfer"
	"sanus/sanus-sdk/sanus/sdk"

	"github.com/btcsuite/btcd/txscript"
)

type Script struct {
	wallet *sdk.BTCWallet
}

func NewScriptHandler(wallet *sdk.BTCWallet) *Script {
	return &Script{wallet: wallet}
}

type ScriptResponse struct {
	Hash string `json:"hash"`
}

func (s *Script) Issuance(r issuance.ColoredData, resp *ScriptResponse) (err error) {
	var script []byte
	if script, err = r.Encode(80); err != nil {
		return err
	}
	script = append([]byte{txscript.OP_RETURN, byte(len(script))}, script...)
	resp.Hash = hex.EncodeToString(script)
	return
}

func (s *Script) Transfer(r transfer.ColoredData, resp *ScriptResponse) (err error) {
	var script []byte
	if script, err = r.Encode(80); err != nil {
		return err
	}
	script = append([]byte{txscript.OP_RETURN, byte(len(script))}, script...)
	resp.Hash = hex.EncodeToString(script)
	return
}
