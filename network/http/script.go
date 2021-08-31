package http

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"sanus/sanus-sdk/cc/issuance"
	"sanus/sanus-sdk/cc/transfer"

	"github.com/btcsuite/btcd/txscript"
)

type ScriptRequest struct {
	Data *issuance.ColoredData `json:"data"`
}

func (server *HTTPServer) Script(w http.ResponseWriter, r *http.Request) *AppResponse {
	var script []byte
	var err error
	if r.URL.Query().Get("type") == "issuance" {
		var request ScriptRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return &AppResponse{Error: err}
		}
		if script, err = request.Data.Encode(80); err != nil {
			return &AppResponse{Error: err}
		}
	} else {
		var request = struct {
			Data *transfer.ColoredData `json:"data"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return &AppResponse{Error: err}
		}
		if script, err = request.Data.Encode(80); err != nil {
			return &AppResponse{Error: err}
		}
	}

	script = append([]byte{txscript.OP_RETURN, byte(len(script))}, script...)
	scriptStr := hex.EncodeToString(script)
	return &AppResponse{Response: fmt.Sprintf("Script:%v", scriptStr)}
}
