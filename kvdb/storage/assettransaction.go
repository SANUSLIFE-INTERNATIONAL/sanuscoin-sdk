package storage

import (
	"sanus/sanus-sdk/entity"
	"sanus/sanus-sdk/kvdb/driver"
)

type AssetTxDB struct {
	db driver.Driver
}

func (db *AssetTxDB) Driver() driver.Driver {
	return db.db
}

func NewAssetTxDB(db driver.Driver) *AssetTxDB {
	return &AssetTxDB{
		db: db,
	}
}

func (db *AssetTxDB) Update(raw *entity.AssetTransactionRaw) error {
	key := raw.TxId
	result, err := db.db.Get(key)
	if err != nil {
		//@TODO implement logging
		return err
	}
	var assetTxEntity = &entity.AssetTransactionEntity{}
	if err := assetTxEntity.From([]byte(key), result); err != nil {
		return err
	}
	assetTxEntity.Append(raw)
	return db.db.Set(key, assetTxEntity.Value())
}
