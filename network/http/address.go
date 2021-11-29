package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type NewAddressRequest struct {
	Account string `json:"account"`
}

func (server *HTTPServer) NewAddress(w http.ResponseWriter, r *http.Request) *AppResponse {
	address, err := server.wallet.NewAddress()
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
	PublicKey string `json:"publicKey"`
}

func (server *HTTPServer) ImportAddress(w http.ResponseWriter, r *http.Request) *AppResponse {
	var request ImportAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return &AppResponse{Error: err}
	}
	if request.PublicKey == "" {
		return &AppResponse{Error: fmt.Errorf("private key can't be empty")}
	}
	address, err := server.wallet.ImportAddress(request.PublicKey)
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

func (server *HTTPServer) ListAddresses(w http.ResponseWriter, r *http.Request) *AppResponse {
	list, err := server.wallet.List()
	return &AppResponse{Error: err, Response: list, Code: 200}
}
