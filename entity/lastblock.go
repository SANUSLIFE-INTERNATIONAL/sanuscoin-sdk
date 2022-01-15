package entity

import (
	"encoding/json"
)

type LastBlockEntity struct {
	Index int    `json:"index"`
	Hash  string `json:"hash"`
}

func (entity *LastBlockEntity) Key() []byte {
	return []byte("last_block")
}
func (entity *LastBlockEntity) Value() []byte {
	bytes, _ := json.Marshal(entity)
	return bytes
}
func (entity *LastBlockEntity) From(index int, hash string) *LastBlockEntity {
	return &LastBlockEntity{
		Index: index,
		Hash:  hash,
	}
}
