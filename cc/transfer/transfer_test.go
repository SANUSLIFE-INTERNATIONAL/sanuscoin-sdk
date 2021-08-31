package transfer

import (
	"reflect"
	"testing"

	"sanus/sanus-sdk/cc/utils"
)

func TestColoredData_Encode(t *testing.T) {
	type fields struct {
		Type        string
		Protocol    int64
		Version     int64
		TorrentHash []byte
		Sha2        []byte
		NoRules     bool
		MultiSig    []MultiSigData
		Payments    []*utils.PaymentData
	}
	type args struct {
		byteSize int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantHash []byte
		wantErr  bool
	}{
		{
			name: "OK",
			fields: fields{
				Type:        "transfer",
				Protocol:    17219,
				Version:     0x02,
				TorrentHash: nil,
				Sha2:        nil,
				NoRules:     false,
				MultiSig:    nil,
				Payments: []*utils.PaymentData{
					{
						Skip:    false,
						Range:   false,
						Percent: false,
						Output:  0,
						Amount:  800000,
					},
				},
			},
			args:    args{byteSize: 80},
			wantErr: false,
			wantHash: []byte{
				67, 67, 2, 21, 0, 32, 133,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cd := &ColoredData{
				Type:        tt.fields.Type,
				Protocol:    tt.fields.Protocol,
				Version:     tt.fields.Version,
				TorrentHash: tt.fields.TorrentHash,
				Sha2:        tt.fields.Sha2,
				NoRules:     tt.fields.NoRules,
				MultiSig:    tt.fields.MultiSig,
				Payments:    tt.fields.Payments,
			}
			gotHash, err := cd.Encode(tt.args.byteSize)

			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHash, tt.wantHash) {
				t.Errorf("Encode() gotHash = %v, want %v", gotHash, tt.wantHash)
			}
		})
	}
}
