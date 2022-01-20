package http

import (
	"encoding/json"
	"net/http"

	"github.com/btcsuite/btcutil"
)

type BalanceRequest struct {
	Address string `json:"address"`
}

type BalanceResponse struct {
	SNC int     `json:"snc"`
	BTC float64 `json:"btc"`
}

func (server *HTTPServer) Balance(w http.ResponseWriter, r *http.Request) *AppResponse {
	var request BalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &AppResponse{Error: err, Code: 400}
	}
	var addr, err = btcutil.DecodeAddress(request.Address, server.wallet.GetNetParams())
	if err != nil {
		return &AppResponse{
			Code:  400,
			Error: err,
		}
	}
	var balance = &BalanceResponse{}
	balance.BTC, balance.SNC, err = server.wallet.Balance(addr)
	return &AppResponse{
		Response: balance,
		Error:    err,
	}
}
