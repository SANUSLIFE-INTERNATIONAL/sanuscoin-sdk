package http

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"sanus/sanus-sdk/misc/random"
	"sanus/sanus-sdk/sanus/daemon"
	"sanus/sanus-sdk/sanus/sdk"

	"github.com/btcsuite/btcutil"
)

func (server *HTTPServer) Seed(w http.ResponseWriter, r *http.Request) *AppResponse {
	mnemonic := r.FormValue("mnemonic")
	publicPass := r.FormValue("public")
	server.Infof("Generating seed phrase Mnemonic <%v> PublicPass <%v>", mnemonic, publicPass)
	seed := sdk.NewSeed(mnemonic, publicPass)
	return &AppResponse{
		Response: hex.EncodeToString(seed),
		Code:     200,
	}
}

func (server *HTTPServer) CreateWallet(w http.ResponseWriter, r *http.Request) *AppResponse {
	publicPass := r.FormValue("public")
	privatePass := r.FormValue("private")
	seed := r.FormValue("seed")

	seedHex, err := hex.DecodeString(seed)
	if err != nil {
		return &AppResponse{
			Error: err,
			Code:  400,
		}
	}

	if err := server.wallet.Create([]byte(publicPass), []byte(privatePass), seedHex); err != nil {
		server.Infof("error caused when trying to create wallet | %v", err)
		return &AppResponse{
			Error: err,
			Code:  400,
		}
	}
	return &AppResponse{
		Response: "Success",
		Code:     200,
	}
}

func (server *HTTPServer) OpenWallet(w http.ResponseWriter, r *http.Request) *AppResponse {
	public := r.FormValue("public")
	if err := server.wallet.Open([]byte(public)); err != nil {
		server.Errorf("error caused when trying to open wallet | %v", err)
		return &AppResponse{
			Error:    err,
			Response: "Failed",
			Code:     400,
		}
	}
	return &AppResponse{
		Response: "Success",
		Code:     200,
	}
}

func (server *HTTPServer) Unlock(w http.ResponseWriter, r *http.Request) *AppResponse {
	private := r.FormValue("private")
	if err := server.wallet.Unlock([]byte(private)); err != nil {
		return &AppResponse{
			Error:    err,
			Response: "Failed",
			Code:     400,
		}
	}
	return &AppResponse{
		Response: "Success",
		Code:     200,
	}
}

func (server *HTTPServer) Lock(w http.ResponseWriter, r *http.Request) *AppResponse {
	server.wallet.Lock()
	return &AppResponse{
		Response: "Success",
		Code:     200,
	}
}

func (server *HTTPServer) TestDo(w http.ResponseWriter, r *http.Request) *AppResponse {
	address := r.FormValue("address")
	addr, _ := btcutil.DecodeAddress(address, daemon.ActiveNetParams.Params)
	balance, err := server.wallet.SNCBalance(addr)
	return &AppResponse{
		Response: fmt.Sprintf("Balance:%v", balance),
		Error:    err,
	}
}

func (server *HTTPServer) Synced(w http.ResponseWriter, r *http.Request) *AppResponse {
	return &AppResponse{
		Response: server.wallet.Synced(),
		Code:     200,
	}
}

func (server *HTTPServer) NewAddress(w http.ResponseWriter, r *http.Request) *AppResponse {
	account := r.FormValue("account")
	if account == "" {
		account = random.RandStringRunes(16)
	}
	address, err := server.wallet.NewAddress(account)
	if err != nil {
		return &AppResponse{
			Error: err,
			Code:  200,
		}
	}
	return &AppResponse{
		Error:    err,
		Response: address.EncodeAddress(),
		Code:     200,
	}
}

func (server *HTTPServer) UnspentTX(w http.ResponseWriter, r *http.Request) *AppResponse {
	addr := r.FormValue("address")
	if addr == "" {
		return &AppResponse{Error: fmt.Errorf("address can't be empty"), Code: 400}
	}
	address, err := btcutil.DecodeAddress(addr, daemon.ActiveNetParams.Params)
	if err != nil {
		return &AppResponse{Error: err, Code: 400}
	}
	txs, err := server.wallet.UnspentTx(address)
	if err != nil {
		return &AppResponse{Error: err, Code: 400}
	}
	return &AppResponse{Response: txs, Code: 200}
}

func (server *HTTPServer) SNCBalance(w http.ResponseWriter, r *http.Request) *AppResponse {
	return nil
}

func (server *HTTPServer) NetworkStatus(w http.ResponseWriter, r *http.Request) *AppResponse {
	return &AppResponse{Error: nil, Response: server.wallet.NetworkStatus()}
}

func (server *HTTPServer) BTCBalance(w http.ResponseWriter, r *http.Request) *AppResponse {
	return nil
}
