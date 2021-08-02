package encdec

import (
	"reflect"
	"testing"

	"sanus/sanus-sdk/cc/utils"
)

func Test_transferPaymentEncode(t *testing.T) {
	type args struct {
		paymentObject *utils.PaymentData
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				paymentObject: &utils.PaymentData{
					Skip:    false,
					Range:   false,
					Percent: false,
					Output:  1,
					Amount:  15,
				},
			},
			want:    []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transferPaymentEncode(tt.args.paymentObject)
			if (err != nil) != tt.wantErr {
				t.Errorf("transferPaymentEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transferPaymentEncode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
