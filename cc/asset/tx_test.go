package asset

import (
	"testing"

	"sanus/sanus-sdk/cc/issuance"
	"sanus/sanus-sdk/cc/transfer"
	"sanus/sanus-sdk/cc/utils"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

func TestCCTransaction_GetAssetOutput(t *testing.T) {
	type fields struct {
		Tx         *wire.MsgTx
		Issuance   *issuance.ColoredData
		Transfer   *transfer.ColoredData
		Type       string
		Input      map[int]*CCVin
		Output     map[int]*CCVout
		Overflow   bool
		registered bool
	}

	firstTxIn := wire.NewTxIn(wire.NewOutPoint(&chainhash.Hash{}, 0), nil, nil)
	secondTxIn := wire.NewTxIn(wire.NewOutPoint(&chainhash.Hash{}, 1), nil, nil)

	firstTxOut := wire.NewTxOut(10, nil)
	secondTxOut := wire.NewTxOut(100, nil)

	var cctx = fields{
		Tx:       nil,
		Issuance: nil,
		Transfer: &transfer.ColoredData{
			Payments: []*utils.PaymentData{
				{
					Range:   false,
					Skip:    false,
					Percent: false,
					Input:   0,
					Output:  0,
					Amount:  10000000,
				},
			},
		},
		Type: "Transfer",
		Input: map[int]*CCVin{
			0: {
				Input: firstTxIn,
				Assets: map[int]*Asset{
					0: {
						AssetId:           "La8e7WhGAEfiT9JGTmyPJopZhkRMwiEPz4uBEG",
						Amount:            1111111110000000,
						IssueTxid:         "d72207fed0f86f73cf9a2ccdc1e9a02dcfea60a46762cbb97f2cbea60d72057b",
						Divisibility:      7,
						LockStatus:        true,
						AggregationPolicy: "aggregatable",
					},
				},
			},
			1: {
				Input:  secondTxIn,
				Assets: map[int]*Asset{},
			},
		},
		Output: map[int]*CCVout{
			0: {
				Out:    firstTxOut,
				Assets: map[int]*Asset{},
			},
			1: {
				Out:    secondTxOut,
				Assets: map[int]*Asset{},
			},
			2: {
				Out:    firstTxOut,
				Assets: map[int]*Asset{},
			},
			3: {
				Out:    secondTxOut,
				Assets: map[int]*Asset{},
			},
		},
		Overflow: false,
	}
	tests := []struct {
		name       string
		fields     fields
		wantLength int
		wantErr    bool
	}{
		{
			name:       "2 assets",
			fields:     cctx,
			wantLength: 2,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &CCTransaction{
				Tx:         tt.fields.Tx,
				Issuance:   tt.fields.Issuance,
				Transfer:   tt.fields.Transfer,
				Type:       tt.fields.Type,
				Input:      tt.fields.Input,
				Output:     tt.fields.Output,
				Overflow:   tt.fields.Overflow,
				registered: tt.fields.registered,
			}

			for x := 0; x < 100; x++ {
				got, err := tx.GetAssetOutput()
				if (err != nil) != tt.wantErr {
					t.Errorf("GetAssetOutput() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if len(got) != tt.wantLength {
					t.Errorf("GetAssetOutput() got = %v, want %v", got, tt.wantLength)
				}
			}

		})
	}
}
