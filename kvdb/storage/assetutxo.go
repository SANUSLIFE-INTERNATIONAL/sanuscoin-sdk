package storage

import (
	"sanus/sanus-sdk/entity"
	"sanus/sanus-sdk/kvdb/driver"
)

type AssetUtxoDB struct {
	db driver.Driver
}

func NewAssetUtxoDB(db driver.Driver) *AssetUtxoDB {
	return &AssetUtxoDB{db: db}
}

func (db *AssetUtxoDB) Update(raw *entity.AssetUtxoRaw) error {
	key := raw.TxId
	result, err := db.db.Get(key)
	if err != nil {
		//@TODO implement logging
		return err
	}
	var assetTxEntity = &entity.AssetUtxoEntity{}
	if err := assetTxEntity.From([]byte(key), result); err != nil {
		return err
	}
	assetTxEntity.Append(raw)
	return db.db.Set(key, assetTxEntity.Value())
}
