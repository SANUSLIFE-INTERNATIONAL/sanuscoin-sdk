package entity

import (
	"encoding/json"
)

type AssetUtxoRaw struct {
	AssetId string
	TxId    string
	Index   int
}

type AssetUtxoEntity struct {
	key  []byte
	data []*AssetUtxoRaw
}

func (entity *AssetUtxoEntity) Append(raw *AssetUtxoRaw) {
	entity.data = append(entity.data, raw)
}

func (entity *AssetUtxoEntity) Key() []byte {
	return entity.key
}

func (entity *AssetUtxoEntity) From(key, value []byte) error {
	entity.key = key
	if len(value) == 0 {
		return nil
	}
	return json.Unmarshal(value, &entity.data)

}

func (entity *AssetUtxoEntity) Value() []byte {
	bytes, _ := json.Marshal(entity.data)
	return bytes
}
