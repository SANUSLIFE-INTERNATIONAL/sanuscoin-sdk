package http

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"sanus/sanus-sdk/sanus/sdk"
)

type CreateWalletRequest struct {
	Public  string `json:"public"`
	Private string `json:"private"`
	Seed    string `json:"seed"`
}

func (server *HTTPServer) CreateWallet(w http.ResponseWriter, r *http.Request) *AppResponse {
	var request CreateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &AppResponse{Error: err}
	}
	seedHex, err := hex.DecodeString(request.Seed)
	if err != nil {
		return &AppResponse{
			Error: err,
			Code:  400,
		}
	}
	if err := server.wallet.Create([]byte(request.Public), []byte(request.Private), seedHex); err != nil {
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

type SeedRequest struct {
	Mnemonic string `json:"mnemonic"`
	Public   string `json:"public"`
}

func (server *HTTPServer) Seed(w http.ResponseWriter, r *http.Request) *AppResponse {
	var request SeedRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &AppResponse{Error: err}
	}
	server.Infof("Generating seed phrase Mnemonic <%v> PublicPass <%v>", request.Mnemonic, request.Public)
	seed := sdk.NewSeed(request.Mnemonic, request.Public)
	return &AppResponse{
		Response: hex.EncodeToString(seed),
		Code:     200,
	}
}

type OpenWalletRequest struct {
	Public string `json:"public"`
}

func (server *HTTPServer) OpenWallet(w http.ResponseWriter, r *http.Request) *AppResponse {
	var request OpenWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &AppResponse{Error: err}
	}
	if err := server.wallet.Open([]byte(request.Public)); err != nil {
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

type UnlockWalletRequest struct {
	Private string `json:"private"`
}

func (server *HTTPServer) Unlock(w http.ResponseWriter, r *http.Request) *AppResponse {
	var request UnlockWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &AppResponse{Error: err}
	}
	if err := server.wallet.Unlock([]byte(request.Private)); err != nil {
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
