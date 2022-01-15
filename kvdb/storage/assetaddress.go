package storage

import (
	"sanus/sanus-sdk/entity"
	"sanus/sanus-sdk/kvdb/driver"
)

type AssetAddressDB struct {
	db driver.Driver
}

func NewAssetAddressDB(db driver.Driver) *AssetAddressDB {
	return &AssetAddressDB{db: db}
}

func (db *AssetAddressDB) Update(address, assetId string) error {
	result, err := db.db.Get(address)
	if err != nil {
		//@TODO implement logging
		return err
	}
	var assetTxEntity = &entity.AssetAddressEntity{}
	if err := assetTxEntity.From([]byte(address), result); err != nil {
		return err
	}
	assetTxEntity.Append(assetId)
	return db.db.Set(address, assetTxEntity.Value())
}
