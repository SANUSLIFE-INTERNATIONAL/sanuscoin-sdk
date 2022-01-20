package asset

import (
	"encoding/json"
	"math"

	"sanus/sanus-sdk/cc/utils"
)

var lockPadding = map[string]int{
	"aggregatable": 0x20CE,
	"hybrid":       0x2102,
	"dispersed":    0x20e4,
}
var unlockPadding = map[string]int{
	"aggregatable": 0x2e37,
	"hybrid":       0x2e6b,
	"dispersed":    0x2e4e,
}

type opcodeData struct {
	start byte
	end   byte
}

var opCodes = map[string]opcodeData{
	"issuance": {start: 0x00, end: 0x0f},
	"Transfer": {start: 0x10, end: 0x1f},
}

var EncLookup = map[byte]string{}

func InitENCLookup() {
	for name, transactionType := range opCodes {
		for j := transactionType.start; j < transactionType.end; j++ {
			EncLookup[j] = name
		}
	}
}

type AssetSlice []*Asset

func (as AssetSlice) IndexExists(index int) bool {
	return len(as) > index
}

func (as AssetSlice) AssetByIndex(index int) *Asset {
	return as[index]
}

type AssetDoubleSlice []AssetSlice

func (ads AssetDoubleSlice) IndexExists(index int) bool {
	return len(ads) > index
}

func (ads AssetDoubleSlice) SliceByIndex(index int) []*Asset {
	return ads[index]
}

type Asset struct {
	AssetId           string `json:"assetId"`
	Amount            int    `json:"amount"`
	IssueTxid         string `json:"issueTxId"`
	Divisibility      int    `json:"divisibility"`
	LockStatus        bool   `json:"lockStatus"`
	AggregationPolicy string `json:"aggregationPolicy"`
}

func (a *Asset) String() string {
	bytes, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func (a *Asset) Id() string {
	if a == nil {
		return ""
	}
	return a.AssetId
}

func isTransfer(assets map[int]map[int]*Asset, payments []*utils.PaymentData, tx *CCTransaction) bool {
	var _payments = make([]*utils.PaymentData, len(payments))
	for i, p := range payments {
		_p := *p
		_payments[i] = &_p
	}
	var _inputs = make(map[int]*CCVin, len(tx.Input))
	for key, value := range tx.Input {
		if len(value.Assets) == 0 {
			value.Assets = make(map[int]*Asset)
		}
		cpBytes, _ := json.Marshal(value)
		var copiedCCVin = CCVin{}
		json.Unmarshal(cpBytes, &copiedCCVin)
		_inputs[key] = &copiedCCVin

	}
	var currentInputIndex = 0
	var currentAssetIndex = 0
	var payment *utils.PaymentData
	var currentAsset *Asset
	var currentAmount int
	var lastPaymentIndex = -1 // aggregate only if paying the same payment

	for i := 0; i < len(_payments); i++ {
		payment = _payments[i]

		if !isPaymentSimple(payment) {
			return false
		}

		if int(payment.Input) >= len(tx.Input) {
			return false
		}

		if int(payment.Output) >= len(tx.Output) {
			return false
		}
		if payment.Amount <= 0 {
			continue
		}
		if currentInputIndex < int(payment.Input) {
			currentInputIndex = int(payment.Input)
			currentAssetIndex = 0
		}

		_, inputOk := _inputs[currentInputIndex]
		var assetOk = false
		if inputOk {
			_, assetOk = _inputs[currentInputIndex].Assets[currentAssetIndex]
		}

		if !assetOk || !inputOk || currentInputIndex >= len(_inputs) {
			return false
		}

		currentAsset = _inputs[currentInputIndex].Assets[currentAssetIndex]
		currentAmount = int(math.Min(float64(payment.Amount), float64(currentAsset.Amount)))

		if !payment.Skip {
			if _, ok := assets[int(payment.Output)]; !ok {
				assets[int(payment.Output)] = map[int]*Asset{}
			}

			if lastPaymentIndex == i {
				curAssetSliceByPaymentOutput := assets[int(payment.Output)]
				if len(curAssetSliceByPaymentOutput) == 0 {
					return false
				}
				currentAssetByPaymentOutput, ok := curAssetSliceByPaymentOutput[len(curAssetSliceByPaymentOutput)-1]
				if !ok {
					return false
				}
				if currentAssetByPaymentOutput.AssetId != currentAsset.AssetId {

					return false
				}
				if currentAsset.AggregationPolicy != "aggregatable" {
					return false
				}
				assets[int(payment.Output)][len(assets[int(payment.Output)])-1].Amount += currentAmount
			} else {
				realIndex := len(assets[int(payment.Output)])
				assets[int(payment.Output)][realIndex] = &Asset{
					AssetId:           currentAsset.AssetId,
					Amount:            currentAmount,
					IssueTxid:         currentAsset.IssueTxid,
					Divisibility:      currentAsset.Divisibility,
					LockStatus:        currentAsset.LockStatus,
					AggregationPolicy: currentAsset.AggregationPolicy,
				}

			}
		}
		currentAsset.Amount -= currentAmount
		payment.Amount -= currentAmount

		if currentAsset.Amount == 0 {
			currentAssetIndex++
			checkIfAssetOk := func(currentInputIndex, currentAssetIndex int) bool {
				return currentAssetIndex > len(_inputs[currentInputIndex].Assets)-1
			}
			checkIfInputOk := func(currentInputIndex int) bool {
				_, ok := _inputs[currentInputIndex]
				return ok
			}
			for checkIfInputOk(currentInputIndex) && checkIfAssetOk(currentInputIndex, currentAssetIndex) {
				currentAssetIndex = 0
				currentInputIndex++
			}
		}
		lastPaymentIndex = i
		if payment.Amount > 0 {
			i--
		}
	}
	transferToLastOutput(assets, _inputs, len(tx.Output)-1)
	return true
}

func transferToLastOutput(assets map[int]map[int]*Asset, inputs map[int]*CCVin, index int) {
	var assetsToTransfer = map[int]*Asset{}
	for _, a := range inputs {
		for _, assetItem := range a.Assets {
			assetsToTransfer[len(assetsToTransfer)] = assetItem
		}
	}
	var assetsIndexes = map[string]int{}
	var lastOutputAssets = map[int]*Asset{}
	for _, curAsset := range assetsToTransfer {
		as, ok := assetsIndexes[curAsset.AssetId]
		if curAsset.AggregationPolicy == "aggregatable" && ok {
			lastOutputAssets[as].Amount += curAsset.Amount
		} else if curAsset.Amount > 0 {
			if _, ok := assetsIndexes[curAsset.AssetId]; !ok {
				assetsIndexes[curAsset.AssetId] = len(lastOutputAssets)
			}
			lastOutputAssets[len(lastOutputAssets)] = &Asset{
				AssetId:           curAsset.AssetId,
				Amount:            curAsset.Amount,
				IssueTxid:         curAsset.IssueTxid,
				Divisibility:      curAsset.Divisibility,
				LockStatus:        curAsset.LockStatus,
				AggregationPolicy: curAsset.AggregationPolicy,
			}
		}
	}
	if _, ok := assets[index]; !ok {
		assets[index] = map[int]*Asset{}
	}
	for _, v := range lastOutputAssets {
		assets[index][len(assets[index])] = v
	}
}

func isPaymentSimple(payment *utils.PaymentData) bool {
	return !payment.Range && !payment.Percent
}
