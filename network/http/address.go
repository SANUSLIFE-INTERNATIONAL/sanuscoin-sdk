package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sanus/sanus-sdk/misc/random"
)

type NewAddressRequest struct {
	Account string `json:"account"`
}

func (server *HTTPServer) NewAddress(w http.ResponseWriter, r *http.Request) *AppResponse {
	var request NewAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &AppResponse{Error: err}
	}
	if request.Account == "" {
		request.Account = random.RandStringRunes(16)
	}
	address, err := server.wallet.NewAddress(request.Account)
	if err != nil {
		return &AppResponse{
			Error: fmt.Errorf("error caused when trying to generate new address %v", err),
			Code:  400,
		}
	}
	return &AppResponse{
		Error:    err,
		Response: address.EncodeAddress(),
		Code:     200,
	}
}

type ImportAddressRequest struct {
	PrivateKey string `json:"privateKey"`
}

func (server *HTTPServer) ImportAddress(w http.ResponseWriter, r *http.Request) *AppResponse {
	var request ImportAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &AppResponse{Error: err}
	}
	if request.PrivateKey == "" {
		return &AppResponse{Error: fmt.Errorf("private key can't be empty")}
	}
	address, err := server.wallet.ImportAddress(request.PrivateKey)
	if err != nil {
		return &AppResponse{
			Error: fmt.Errorf("error caused when trying to generate new address %v", err),
			Code:  400,
		}
	}
	return &AppResponse{
		Error:    err,
		Response: address.EncodeAddress(),
		Code:     200,
	}
}
