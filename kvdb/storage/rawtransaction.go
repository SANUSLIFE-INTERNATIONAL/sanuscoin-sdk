package storage

import (
	"sanus/sanus-sdk/entity"
	"sanus/sanus-sdk/kvdb/driver"
)

type RawTransactionDB struct {
	db driver.Driver
}

func NewRawTransactionDB(db driver.Driver) *RawTransactionDB {
	return &RawTransactionDB{db: db}
}

func (db *RawTransactionDB) Update(raw *entity.RawTransactionEntity) error {
	return db.db.Set(string(raw.Key()), raw.Value())
}

func (db *RawTransactionDB) GetByTxId(id string) (*entity.RawTransactionEntity, error) {
	data, err := db.db.Get(id)
	if err != nil {
		return nil, err
	}
	if data != nil {
		raw := &entity.RawTransactionEntity{}
		if err = raw.From([]byte(id), data); err != nil {
			return nil, err
		}
		return raw, nil
	}
	return &entity.RawTransactionEntity{}, nil
}
