package asset

import (
	"encoding/json"
	"fmt"

	"sanus/sanus-sdk/cc/issuance"
	"sanus/sanus-sdk/cc/transfer"
	"sanus/sanus-sdk/cc/utils"

	"github.com/btcsuite/btcd/wire"
)

type CCTransaction struct {
	Tx         *wire.MsgTx           `json:"tx"`
	Issuance   *issuance.ColoredData `json:"issuance"`
	Transfer   *transfer.ColoredData `json:"transfer"`
	Type       string                `json:"type"`
	Input      map[int]*CCVin        `json:"input"`
	Output     map[int]*CCVout       `json:"output"`
	Overflow   bool                  `json:"overflow"`
	registered bool
}

func (tx *CCTransaction) Key() []byte {
	return []byte(tx.Tx.TxHash().String())
}

func (tx *CCTransaction) IsRegistered() bool {
	return tx.registered
}

func (tx *CCTransaction) Registered() {
	tx.registered = true
}

func (tx *CCTransaction) Value() []byte {
	bytes, _ := json.Marshal(tx)
	return bytes
}

func (tx *CCTransaction) AppendInput(input *CCVin) {
	tx.Input[len(tx.Input)] = input
}

func (tx *CCTransaction) AppendOutput(output *CCVout) {
	tx.Output[len(tx.Output)] = output
}

func (tx *CCTransaction) GetAssetOutput() (map[int]map[int]*Asset, error) {
	if tx.Type == "" {
		return nil, fmt.Errorf("isn't cc Tx")
	}
	var payments []*utils.PaymentData
	var assets = make(map[int]map[int]*Asset)
	if tx.Type == "issuance" {
		id, err := assetId(tx)
		if err != nil {
			return nil, fmt.Errorf("error while decoding issuance Tx %v", err)
		}
		if id == "La8e7WhGAEfiT9JGTmyPJopZhkRMwiEPz4uBEG" {
			tx.Input[0].AppendAsset(&Asset{
				AssetId:           id,
				Amount:            tx.Issuance.Amount,
				IssueTxid:         tx.Tx.TxHash().String(),
				Divisibility:      tx.Issuance.Divisibility,
				LockStatus:        tx.Issuance.LockStatus,
				AggregationPolicy: tx.Issuance.AggregationPolicy,
			})
		}
		payments = tx.Issuance.Payments
	}
	if tx.Type == "Transfer" {
		payments = tx.Transfer.Payments
	}

	var overflow = !isTransfer(assets, payments, tx)

	if overflow {
		// Transfer failed. Transfer all GetAssets in inputs to last output, aggregate those possible
		transferToLastOutput(assets, tx.Input, len(tx.Output)-1)
	}

	tx.Overflow = overflow
	return assets, nil
}
