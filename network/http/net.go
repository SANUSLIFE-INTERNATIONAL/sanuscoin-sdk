package http

import (
	"encoding/json"
	"net/http"
)

const (
	ProtocolVersion = "v1"
)

type appHandler func(http.ResponseWriter, *http.Request) *AppResponse

type AppResponse struct {
	Error    error       `json:"error"`
	Response interface{} `json:"response"`
	Code     int         `json:"code"`
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := fn(w, r)
	response, _ := json.MarshalIndent(resp, "", "")
	w.Write(response)
}