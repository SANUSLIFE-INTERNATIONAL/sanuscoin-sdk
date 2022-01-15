package entity

import (
	"encoding/json"

	"sanus/sanus-sdk/cc/asset"
)

type RawTransactionEntity struct {
	key  []byte
	data *asset.CCTransaction
}

func NewRawTransactionEntity(tx *asset.CCTransaction) *RawTransactionEntity {
	return &RawTransactionEntity{
		key:  tx.Key(),
		data: tx,
	}
}

func (entity *RawTransactionEntity) Data() *asset.CCTransaction {
	return entity.data
}

func (entity *RawTransactionEntity) Key() []byte {
	return entity.key
}
func (entity *RawTransactionEntity) Value() []byte {
	bytes, _ := json.Marshal(entity.data)
	return bytes
}

func (entity *RawTransactionEntity) From(key, value []byte) error {
	entity.key = key
	return json.Unmarshal(value, &entity.data)
}
