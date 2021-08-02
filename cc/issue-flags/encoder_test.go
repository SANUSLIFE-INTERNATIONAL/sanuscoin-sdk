package issue_flags

import (
	"reflect"
	"testing"
)

func TestFlags_Encode(t *testing.T) {
	type fields struct {
		Divisibility      int
		LockStatus        bool
		AggregationPolicy string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "OK",
			fields: fields{
				Divisibility:      0,
				LockStatus:        false,
				AggregationPolicy: "aggregatable",
			},
			want:    []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Flags{
				Divisibility:      tt.fields.Divisibility,
				LockStatus:        tt.fields.LockStatus,
				AggregationPolicy: tt.fields.AggregationPolicy,
			}
			got, err := f.Encode()
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

func TestFlags_Encode1(t *testing.T) {
	type fields struct {
		Divisibility      int
		LockStatus        bool
		AggregationPolicy string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "aggregatable:ok",
			fields: fields{
				Divisibility:      2,
				LockStatus:        false,
				AggregationPolicy: "aggregatable",
			},
			want:    []byte{64},
			wantErr: false,
		},
		{
			name: "unknown: Error",
			fields: fields{
				Divisibility:      2,
				LockStatus:        false,
				AggregationPolicy: "unknown",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Flags{
				Divisibility:      tt.fields.Divisibility,
				LockStatus:        tt.fields.LockStatus,
				AggregationPolicy: tt.fields.AggregationPolicy,
			}
			got, err := f.Encode()
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
		consume func(int) []byte
	}
	tests := []struct {
		name string
		args args
		want *Flags
	}{
		{
			name: "aggregatable: OK",
			args: args{consume: consumer([]byte{64})},
			want: &Flags{
				Divisibility:      2,
				LockStatus:        false,
				AggregationPolicy: "aggregatable",
			},
		},
		{
			name: "unknown: Error",
			args: args{consume: consumer([]byte{})},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decode(tt.args.consume); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func consumer(buff []byte) func(int) []byte {
	var curr = 0
	return func(lgt int) []byte {
		if len(buff) < lgt+curr {
			return []byte{}
		}
		bytes := buff[curr : curr+lgt]
		curr = curr + lgt
		return bytes
	}
}
