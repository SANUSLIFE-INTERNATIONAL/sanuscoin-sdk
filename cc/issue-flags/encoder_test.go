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
