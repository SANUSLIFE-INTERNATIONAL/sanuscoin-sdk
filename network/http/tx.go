package http

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"sanus/sanus-sdk/sanus/daemon"

	"github.com/btcsuite/btcutil"
)

type UnspentTxRequest struct {
	Address string
}

func (server *HTTPServer) UnspentTX(w http.ResponseWriter, r *http.Request) *AppResponse {
	var request UnspentTxRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &AppResponse{Error: err}
	}
	if request.Address == "" {
		return &AppResponse{Error: fmt.Errorf("address can't be empty"), Code: 400}
	}
	address, err := btcutil.DecodeAddress(request.Address, daemon.ActiveNetParams.Params)
	if err != nil {
		return &AppResponse{Error: err, Code: 400}
	}
	txs, err := server.wallet.UnspentTx(address)
	if err != nil {
		return &AppResponse{Error: err, Code: 400}
	}
	return &AppResponse{Response: txs, Code: 200}
}

type SendTxRequest struct {
	Address  string  `json:"address"`
	Amount   float64 `json:"amount"`
	PkScript string  `json:"pk_script"`
}

func (server *HTTPServer) SendTx(w http.ResponseWriter, r *http.Request) *AppResponse {
	var request SendTxRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &AppResponse{Error: err}
	}
	address, err := btcutil.DecodeAddress(request.Address, server.wallet.GetNetParams())
	if err != nil {
		return &AppResponse{Error: err}
	}
	if request.Amount <= 0 {
		return &AppResponse{
			Error: fmt.Errorf("amount can't be negative or equal 0"),
			Code:  400,
		}
	}
	amountReal, err := btcutil.NewAmount(request.Amount)
	if err != nil {
		return &AppResponse{Error: err}
	}
	var pkScriptByte []byte = nil
	pkScript := request.PkScript
	if pkScript != "" {
		pkScriptByte, err = hex.DecodeString(pkScript)
		if err != nil {
			return &AppResponse{Error: err}
		}
	}
	hash, err := server.wallet.SendTx(address, amountReal, pkScriptByte)
	if err != nil {
		return &AppResponse{Error: fmt.Errorf("error caused when trying to send tx %v", err), Code: 400}
	}
	return &AppResponse{Response: fmt.Sprintf("Hash:%v", hash)}
}
