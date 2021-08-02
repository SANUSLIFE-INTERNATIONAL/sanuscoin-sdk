package sffc

import (
	"reflect"
	"testing"
)

func Test_intToFloatArray(t *testing.T) {
	type args struct {
		number int
		n      interface{}
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "OK:1",
			args: args{
				number: 10,
				n:      2,
			},
			want: []int{1, 3},
		},
		{
			name: "OK:2",
			args: args{
				number: 99,
				n:      2,
			},
			want: []int{99, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intToFloatArray(tt.args.number, tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("intToFloatArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	type args struct {
		amount int
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "OK:10",
			args:    args{amount: 10},
			want:    []byte{byte(10)},
			wantErr: false,
		},
		{
			name:    "OK:99",
			args:    args{amount: 99},
			want:    []byte{38, 48},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.amount)
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
