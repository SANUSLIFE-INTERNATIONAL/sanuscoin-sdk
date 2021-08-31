package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/btcsuite/btcutil"
)

type BalanceRequest struct {
	Address string `json:"address"`
	Coin    string `json:"coin"`
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
	var balance interface{}
	switch request.Coin {
	case "btc":
		balance, err = server.wallet.BTCBalance(addr)
	case "snc":
		balance, err = server.wallet.SNCBalance(addr)
	default:
		balance, err = 0, fmt.Errorf("invalid coin type")
	}
	return &AppResponse{
		Response: fmt.Sprintf("Balance:%v", balance),
		Error:    err,
	}
}
