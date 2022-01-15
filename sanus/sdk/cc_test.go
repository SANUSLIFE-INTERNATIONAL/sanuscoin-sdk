package sdk

import (
	"encoding/hex"
	"testing"
)

func Test_bs58check(t *testing.T) {
	var addrHex, _ = hex.DecodeString("003c176e659bea0f29a3e9bf7880c112b1b31b4dc826268187")
	type args struct {
		payload []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "check address",
			args: args{payload: addrHex},
			want: "1cr9M1wYcuv3gn659v4s5qJihwVUkdxC4k9nBAD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bs58check(tt.args.payload); got != tt.want {
				t.Errorf("bs58check() = %v, want %v", got, tt.want)
			}
		})
	}
}
