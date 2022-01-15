package entity

import (
	"encoding/json"
	"fmt"

	"sanus/sanus-sdk/cc/asset"
)

type UtxoRaw struct {
	Assets map[int]*asset.Asset `json:"assets"`
	TxId   string               `json:"txId"`
	Index  int                  `json:"index"`
}

type UtxoEntity struct {
	key  []byte
	data []*UtxoRaw
}

func (entity *UtxoEntity) Key() []byte {
	return entity.key
}

func (entity *UtxoEntity) Value() []byte {
	bytes, err := json.Marshal(entity.data)
	if err != nil {
		fmt.Println("error caused when trying to marshal utxo data", err)
	}
	return bytes
}

func (entity *UtxoEntity) Append(raw *UtxoRaw) {
	entity.data = append(entity.data, raw)
}

func (entity *UtxoEntity) Data() []*UtxoRaw {
	return entity.data
}

func (entity *UtxoEntity) PutData(index int, value *UtxoRaw) {
	entity.data[index] = value
}

func (entity *UtxoEntity) From(key, value []byte) error {
	entity.key = key
	if len(value) == 0 {
		return nil
	}
	return json.Unmarshal(value, &entity.data)
}
