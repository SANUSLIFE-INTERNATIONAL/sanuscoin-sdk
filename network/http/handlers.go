package http

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"sanus/sanus-sdk/cc/issuance"
	"sanus/sanus-sdk/cc/transfer"
	"sanus/sanus-sdk/misc/random"
	"sanus/sanus-sdk/sanus/daemon"
	"sanus/sanus-sdk/sanus/sdk"

	"github.com/btcsuite/btcd/txscript"
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

func (server *HTTPServer) Synced(w http.ResponseWriter, r *http.Request) *AppResponse {
	return &AppResponse{
		Response: server.wallet.Synced(),
		Code:     200,
	}
}

func (server *HTTPServer) NewAddress(w http.ResponseWriter, r *http.Request) *AppResponse {
	var request struct {
		Account string
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &AppResponse{Error: err}
	}

	if request.Account == "" {
		request.Account = random.RandStringRunes(16)
	}
	address, err := server.wallet.NewAddress(request.Account)
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
	var request struct {
		Address string
	}
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

func (server *HTTPServer) Script(w http.ResponseWriter, r *http.Request) *AppResponse {
	var script []byte
	var err error
	if r.URL.Query().Get("type") == "issuance" {
		var request = struct {
			data *issuance.ColoredData
		}{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return &AppResponse{Error: err}
		}
		if script, err = request.data.Encode(80); err != nil {
			return &AppResponse{Error: err}
		}
	} else {
		var request = struct {
			data *transfer.ColoredData
		}{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return &AppResponse{Error: err}
		}
		if script, err = request.data.Encode(80); err != nil {
			return &AppResponse{Error: err}
		}
	}

	script = append([]byte{txscript.OP_RETURN, byte(len(script))}, script...)
	scriptStr := hex.EncodeToString(script)
	return &AppResponse{Response: fmt.Sprintf("Script:%v", scriptStr)}
}

func (server *HTTPServer) SendTx(w http.ResponseWriter, r *http.Request) *AppResponse {
	addr := r.FormValue("addr")
	address, err := btcutil.DecodeAddress(addr, server.wallet.GetNetParams())
	if err != nil {
		return &AppResponse{Error: err}
	}
	amount := r.FormValue("amount")
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return &AppResponse{Error: err}
	}
	amountReal, err := btcutil.NewAmount(amountFloat)
	if err != nil {
		return &AppResponse{Error: err}
	}
	pkScript := r.FormValue("script")
	pkScriptByte, err := hex.DecodeString(pkScript)
	if err != nil {
		return &AppResponse{Error: err}
	}
	hash, err := server.wallet.SendTx(address, amountReal, pkScriptByte)
	if err != nil {
		return &AppResponse{Error: err}
	}
	return &AppResponse{Response: fmt.Sprintf("Hash:%v", hash)}
}

func (server *HTTPServer) Balance(w http.ResponseWriter, r *http.Request) *AppResponse {
	address := r.FormValue("address")
	coin := r.FormValue("coin")
	if coin == "" {
		coin = "btc"
	}

	var addr, err = btcutil.DecodeAddress(address, daemon.ActiveNetParams.Params)
	if err != nil {
		return &AppResponse{
			Response: fmt.Sprintf("Error caused"),
			Error:    err,
		}
	}
	var balance int
	switch coin {
	case "btc":
		balance, err = server.wallet.BTCBalance(addr)
	default:
		balance, err = server.wallet.SNCBalance(addr)
	}
	return &AppResponse{
		Response: fmt.Sprintf("Balance:%v", balance),
		Error:    err,
	}
}

func (server *HTTPServer) NetworkStatus(w http.ResponseWriter, r *http.Request) *AppResponse {
	return &AppResponse{Error: nil, Response: server.wallet.NetworkStatus()}
}

func (server *HTTPServer) BTCBalance(w http.ResponseWriter, r *http.Request) *AppResponse {
	return nil
}
