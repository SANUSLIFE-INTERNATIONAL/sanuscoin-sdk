package storage

import (
	"fmt"
	"strconv"
	"testing"

	"sanus/sanus-sdk/cc/asset"
	"sanus/sanus-sdk/entity"
)

var testDb = New()

func TestAssetUtxoDB_Update(t *testing.T) {

	testutxoDB := testDb.Utxo()
	for x := 0; x < 500; x++ {
		id := strconv.Itoa(x)
		var utxo = &entity.UtxoRaw{
			Assets: map[int]*asset.Asset{
				0: {
					Amount: 100,
				},
			},
			TxId:  id,
			Index: x,
		}
		testutxoDB.Update(utxo)
	}
}

func TestAssetUtxoDB_(t *testing.T) {
	var testutxoDB = testDb.Utxo()
	var count = 0
	for x := 0; x < 500; x++ {
		id := strconv.Itoa(x)

		raw, _ := testutxoDB.GetByTxIdAndIndex(id, x)
		if len(raw) > 0 {
			for _, v := range raw {
				if v.Amount > -1 {
					count++
				}
			}
		}
	}

	fmt.Println(count)
}
