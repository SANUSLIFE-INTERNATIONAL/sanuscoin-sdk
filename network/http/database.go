package http

import (
	"net/http"
	"strconv"
)

type RawTransactionCond struct {
	TxId string `json:"txId"`
}

func (server *HTTPServer) RawTransaction(w http.ResponseWriter, r *http.Request) *AppResponse {
	rValues := r.URL.Query()
	txId := rValues.Get("tx")
	if txId == "" {
		return &AppResponse{
			Response: "nil",
			Code:     404,
		}
	}
	db := server.db.RawTransaction()
	result, err := db.GetByTxId(txId)
	if err != nil {
		return &AppResponse{Response: err.Error()}
	}
	return &AppResponse{Response: result.Data()}
}

func (server *HTTPServer) Utxo(w http.ResponseWriter, r *http.Request) *AppResponse {
	rValues := r.URL.Query()
	txId := rValues.Get("tx")
	if txId == "" {
		return &AppResponse{
			Response: "nil",
			Code:     404,
		}
	}
	index := rValues.Get("index")
	if index == "" {
		index = "0"
	}
	indexInt, _ := strconv.Atoi(index)

	db := server.db.Utxo()
	result, err := db.GetByTxIdAndIndex(txId, indexInt)
	if err != nil {
		return &AppResponse{Response: err.Error()}
	}
	return &AppResponse{Response: result}
}
