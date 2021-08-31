package http

import (
	"net/http"
)

func (server *HTTPServer) NetworkStatus(w http.ResponseWriter, r *http.Request) *AppResponse {
	return &AppResponse{Error: nil, Response: server.wallet.NetworkStatus()}
}
