package http

import (
	"fmt"
	"net/http"
)

func (server *HTTPServer) NetworkStatus(w http.ResponseWriter, r *http.Request) *AppResponse {
	return &AppResponse{Error: nil, Response: server.wallet.NetworkStatus()}
}

func (server *HTTPServer) TestMethod(w http.ResponseWriter, r *http.Request) *AppResponse {
	go func() {
		if err := server.wallet.Scan(); err != nil {
			fmt.Println("Error caused when trying to scan", err)
		}
	}()
	return &AppResponse{Response: "ok"}
}
