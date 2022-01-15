package storage

import (
	"encoding/json"

	"sanus/sanus-sdk/entity"
	"sanus/sanus-sdk/kvdb/driver"
)

type LastBlockDB struct {
	db driver.Driver
}

func NewLastBlockDB(db driver.Driver) *LastBlockDB {
	return &LastBlockDB{db: db}
}

func (db *LastBlockDB) Update(raw *entity.LastBlockEntity) error {
	return db.db.Set(string(raw.Key()), raw.Value())
}

func (db *LastBlockDB) GetLastIndex() (*entity.LastBlockEntity, error) {
	data, err := db.db.Get("last_block")
	if err != nil {
		return nil, err
	}
	var lastBlockEntity entity.LastBlockEntity
	if err = json.Unmarshal(data, &lastBlockEntity); err != nil {
		return nil, err
	}
	return &lastBlockEntity, nil
}
