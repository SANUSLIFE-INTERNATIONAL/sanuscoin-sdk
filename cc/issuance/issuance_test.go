package issuance

import (
	"encoding/hex"
	"reflect"
	"testing"

	"sanus/sanus-sdk/cc/transfer"
	"sanus/sanus-sdk/cc/utils"
)

var (
	sha2, _        = hex.DecodeString("03ffdf3d6790a21c5fc97a62fe1abc5f66922d7dee3725261ce02e86f078d190")
	torrentHash, _ = hex.DecodeString("46b7e0d000d69330ac1caa48c6559763828762e1")

	validColoredData = &ColoredData{
		Amount:            15,
		Protocol:          0x4343,
		Version:           0x02,
		Divisibility:      2,
		LockStatus:        true,
		AggregationPolicy: "aggregatable",
		MultiSig:          []transfer.MultiSigData{},
		NoRules:           false,
		Payments: []*utils.PaymentData{
			{
				Skip:    false,
				Range:   false,
				Percent: false,
				Output:  1,
				Amount:  15,
			},
		},
	}

	validColoredDataBytes = []byte{67, 67, 2, 6, 15, 1, 15, 80}
)

func TestColoredData_Encode(t *testing.T) {
	type fields struct {
		Amount            int
		Protocol          int64
		Version           int64
		Divisibility      int
		LockStatus        bool
		AggregationPolicy string
		MultiSig          []transfer.MultiSigData
		NoRules           bool
		Sha2              []byte
		TorrentHash       []byte
		Payments          []*utils.PaymentData
	}
	type args struct {
		byteSize int
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "OK [without sha2 and torrentHash]",
			fields: fields{
				Amount:            15,
				Protocol:          0x4343,
				Version:           0x02,
				Divisibility:      2,
				LockStatus:        true,
				AggregationPolicy: "aggregatable",
				NoRules:           false,
				Sha2:              []byte{},
				TorrentHash:       []byte{},
				Payments: []*utils.PaymentData{
					{
						Skip:    false,
						Range:   false,
						Percent: false,
						Output:  1,
						Amount:  15,
					},
				},
			},
			args: args{byteSize: 80},
			want: []byte{
				67, 67, 2, 6, 15, 1, 15, 80,
			},
			wantErr: false,
		},
		{
			name: "OK [with sha2 and torrentHash]",
			fields: fields{
				Amount:            15,
				Protocol:          0x4343,
				Version:           0x02,
				Divisibility:      2,
				LockStatus:        true,
				AggregationPolicy: "aggregatable",
				NoRules:           false,
				Sha2:              sha2,
				TorrentHash:       torrentHash,
				Payments: []*utils.PaymentData{
					{
						Skip:    false,
						Range:   false,
						Percent: false,
						Output:  1,
						Amount:  15,
					},
				},
			},
			args: args{byteSize: 80},
			want: []byte{
				67, 67, 2, 1, 70, 183, 224, 208, 0, 214, 147, 48,
				172, 28, 170, 72, 198, 85, 151, 99, 130, 135, 98,
				225, 3, 255, 223, 61, 103, 144, 162, 28, 95, 201,
				122, 98, 254, 26, 188, 95, 102, 146, 45, 125, 238, 55,
				37, 38, 28, 224, 46, 134, 240, 120, 209, 144, 15, 1,
				15, 80,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cd := &ColoredData{
				Amount:            tt.fields.Amount,
				Protocol:          tt.fields.Protocol,
				Version:           tt.fields.Version,
				Divisibility:      tt.fields.Divisibility,
				LockStatus:        tt.fields.LockStatus,
				AggregationPolicy: tt.fields.AggregationPolicy,
				MultiSig:          tt.fields.MultiSig,
				NoRules:           tt.fields.NoRules,
				Sha2:              tt.fields.Sha2,
				TorrentHash:       tt.fields.TorrentHash,
				Payments:          tt.fields.Payments,
			}
			got, err := cd.Encode(tt.args.byteSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type args struct {
		hexByte []byte
	}
	tests := []struct {
		name     string
		args     args
		wantData *ColoredData
		wantErr  bool
	}{
		{
			name:     "OK",
			args:     args{hexByte: validColoredDataBytes},
			wantData: validColoredData,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := Decode(tt.args.hexByte)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("Decode() \n gotData = %#v, \n want %#v", gotData, tt.wantData)
			}
		})
	}
}
