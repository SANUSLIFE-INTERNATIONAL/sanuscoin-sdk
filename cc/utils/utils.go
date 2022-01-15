package utils

import (
	"encoding/json"
	"math"
)

type PaymentData struct {
	Skip    bool  `json:"skip"`
	Range   bool  `json:"range"`
	Percent bool  `json:"percent"`
	Output  int64 `json:"output"`
	Input   int64 `json:"input"`
	Amount  int   `json:"amount"`
}

func (p *PaymentData) String() string {
	bytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func PadLeadingZeros(hex string, byteSize int) string {
	if byteSize == -1 {
		byteSize = int(math.Ceil(float64(len(hex) / 2)))
	}

	if len(hex) == byteSize*2 || byteSize == 0 {
		return hex
	}
	return PadLeadingZeros("0"+hex, byteSize)
}

func BytesConcat(slices ...[]byte) []byte {
	var hash = make([]byte, 0)
	for x := 0; x < len(slices); x++ {
		hash = append(hash, slices[x]...)
	}
	return hash
}
