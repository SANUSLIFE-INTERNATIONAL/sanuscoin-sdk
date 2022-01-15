package entity

import (
	"encoding/json"
)

type AssetTransactionRaw struct {
	AssetId string `json:"assetId"`
	TxId    string `json:"txId"`
	Type    string `json:"type"`
}

type AssetTransactionEntity struct {
	key  []byte
	data []*AssetTransactionRaw
}

func (entity *AssetTransactionEntity) Append(raw *AssetTransactionRaw) {
	entity.data = append(entity.data, raw)
}

func (entity *AssetTransactionEntity) Key() []byte {
	return entity.key
}

func (entity *AssetTransactionEntity) From(key, value []byte) error {
	entity.key = key
	if len(value) == 0 {
		return nil
	}
	return json.Unmarshal(value, &entity.data)

}

func (entity *AssetTransactionEntity) Value() []byte {
	bytes, _ := json.Marshal(entity.data)
	return bytes
}
