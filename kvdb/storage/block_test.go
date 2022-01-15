package storage

import (
	"reflect"
	"testing"

	"sanus/sanus-sdk/entity"
	"sanus/sanus-sdk/kvdb/driver"
)

func TestLastBlockDB_Update(t *testing.T) {
	coreDB := New()
	type fields struct {
		db driver.Driver
	}
	type args struct {
		raw *entity.LastBlockEntity
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Save new entity",
			fields: fields{db: coreDB.LastBlockDB().db},
			args: args{
				raw: &entity.LastBlockEntity{
					Index: 10,
					Hash:  "abc",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &LastBlockDB{
				db: tt.fields.db,
			}
			if err := db.Update(tt.args.raw); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLastBlockDB_GetLastIndex(t *testing.T) {
	var coreDb = New().LastBlockDB().db
	type fields struct {
		db driver.Driver
	}
	tests := []struct {
		name    string
		fields  fields
		want    *entity.LastBlockEntity
		wantErr bool
	}{
		{
			name:    "Fetch the last block",
			fields:  fields{db: coreDb},
			want:    &entity.LastBlockEntity{Index: 10, Hash: "abc"},
			wantErr: false,
		},
		{
			name:    "Fetch the last block",
			fields:  fields{db: coreDb},
			want:    &entity.LastBlockEntity{Index: 10, Hash: "abc"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &LastBlockDB{
				db: tt.fields.db,
			}
			got, err := db.GetLastIndex()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLastIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLastIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}
