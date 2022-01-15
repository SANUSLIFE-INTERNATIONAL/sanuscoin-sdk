package entity

import (
	"encoding/json"
)

type AssetAddressEntity struct {
	key  []byte
	data []string
}

func (entity *AssetAddressEntity) Key() []byte {
	return entity.key
}

func (entity *AssetAddressEntity) Value() []byte {
	bytes, _ := json.Marshal(entity.data)
	return bytes
}

func (entity *AssetAddressEntity) From(key, value []byte) error {
	entity.key = key
	return json.Unmarshal(value, &entity.data)
}

func (entity *AssetAddressEntity) Append(raw string) {
	entity.data = append(entity.data, raw)
}
