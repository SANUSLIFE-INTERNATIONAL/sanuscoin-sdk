package storage

import (
	"sanus/sanus-sdk/cc/asset"
	"sanus/sanus-sdk/entity"
	"sanus/sanus-sdk/kvdb/driver"
)

type UtxoDB struct {
	db driver.Driver
}

func NewUtxoDB(db driver.Driver) *UtxoDB {
	return &UtxoDB{
		db: db,
	}
}

func (db *UtxoDB) Update(raw *entity.UtxoRaw) error {
	key := raw.TxId
	result, err := db.db.Get(key)
	if err != nil {
		return err
	}
	utxoEntity := &entity.UtxoEntity{}
	if err = utxoEntity.From([]byte(key), result); err != nil {
		return err
	}
	utxoEntity.Append(raw)
	return db.db.Set(key, utxoEntity.Value())
}

func (db *UtxoDB) GetByTxIdAndIndex(tx string, index int) (map[int]*asset.Asset, error) {
	result, err := db.db.Get(tx)
	if err != nil {
		return map[int]*asset.Asset{}, err
	}
	var utxoEntity = &entity.UtxoEntity{}
	if err = utxoEntity.From([]byte(tx), result); err != nil {
		return map[int]*asset.Asset{}, err
	}
	var assets = map[int]*asset.Asset{}
	for _, assetData := range utxoEntity.Data() {
		if assetData.Index == index {
			assets = assetData.Assets
		}
	}
	return assets, nil

}

func (db *UtxoDB) ToEmptyByIndex(raw *entity.UtxoRaw) error {
	key := raw.TxId
	result, err := db.db.Get(key)
	if err != nil {
		return err
	}
	utxoEntity := &entity.UtxoEntity{}
	if err = utxoEntity.From([]byte(key), result); err != nil {
		return err
	}
	for index, utxo := range utxoEntity.Data() {
		if utxo.Index == raw.Index {
			utxo.Assets = map[int]*asset.Asset{}
			utxoEntity.PutData(index, utxo)
		}
	}
	return db.db.Set(key, utxoEntity.Value())
}
