package storage

import (
	"path/filepath"

	"sanus/sanus-sdk/config"
	"sanus/sanus-sdk/kvdb/badger"
)

const (
	assetAddressDBName     = "assetaddress"
	assetTransactionDBName = "assettransaction"
	utxoDBName             = "utxo"
	assetUtxoDBName        = "assetutxo"
	rawTransactionDBName   = "rawtransaction"
	lastBlockDBName        = "lastblock"
)

type DB struct {
	assetTransactionDB *AssetTxDB
	assetAddressDB     *AssetAddressDB
	utxoDB             *UtxoDB
	assetUtxoDB        *AssetUtxoDB
	rawTransactionDB   *RawTransactionDB
	lastBlockDB        *LastBlockDB
}

func New() *DB {
	return &DB{
		assetAddressDB:     assetAddressDB(),
		assetTransactionDB: assetTransactionDB(),
		utxoDB:             utxoDB(),
		assetUtxoDB:        assetUtxoDB(),
		rawTransactionDB:   rawTransactionDB(),
		lastBlockDB:        lastBlockDB(),
	}
}

func (db *DB) AssetAddress() *AssetAddressDB {
	return db.assetAddressDB
}

func (db *DB) LastBlockDB() *LastBlockDB {
	return db.lastBlockDB
}

func (db *DB) AssetUtxo() *AssetUtxoDB {
	return db.assetUtxoDB
}

func (db *DB) AssetTransaction() *AssetTxDB {
	return db.assetTransactionDB
}

func (db *DB) Utxo() *UtxoDB {
	return db.utxoDB
}

func (db *DB) RawTransaction() *RawTransactionDB {
	return db.rawTransactionDB
}

func assetAddressDB() *AssetAddressDB {
	store := badger.DataStore()
	if err := store.Open(filepath.Join(config.AppDataPath(), assetAddressDBName)); err != nil {
		return nil
	}
	return NewAssetAddressDB(store)
}

func assetUtxoDB() *AssetUtxoDB {
	store := badger.DataStore()
	if err := store.Open(filepath.Join(config.AppDataPath(), assetUtxoDBName)); err != nil {
		return nil
	}
	return NewAssetUtxoDB(store)
}

func utxoDB() *UtxoDB {
	store := badger.DataStore()
	if err := store.Open(filepath.Join(config.AppDataPath(), utxoDBName)); err != nil {
		return nil
	}
	return NewUtxoDB(store)
}

func rawTransactionDB() *RawTransactionDB {
	store := badger.DataStore()
	if err := store.Open(filepath.Join(config.AppDataPath(), rawTransactionDBName)); err != nil {
		return nil
	}
	return NewRawTransactionDB(store)
}

func lastBlockDB() *LastBlockDB {
	store := badger.DataStore()
	if err := store.Open(filepath.Join(config.AppDataPath(), lastBlockDBName)); err != nil {
		return nil
	}
	return NewLastBlockDB(store)
}

func assetTransactionDB() *AssetTxDB {
	store := badger.DataStore()
	if err := store.Open(filepath.Join(config.AppDataPath(), assetTransactionDBName)); err != nil {
		return nil
	}
	return NewAssetTxDB(store)
}
