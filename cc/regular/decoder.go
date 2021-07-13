package regular

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"sanus/sanus-sdk/cc/sffc"
)

type PaymentData struct {
	Skip    bool
	Range   bool
	Percent bool
	Output  int64
	Amount  int
}

var flagMask byte = 0xe0
var skipFlag byte = 0x80
var rangeFlag byte = 0x40
var percentFlag byte = 0x20

func EncodeBulk(data []*PaymentData) []byte {
	return nil
}

func DecodeBulk(consume func(int) []byte, paymentsArray []*PaymentData) []*PaymentData {
	if paymentsArray == nil {
		paymentsArray = []*PaymentData{}
	}
	for true {
		paymentData, err := decode(consume)
		if err != nil {
			return paymentsArray
		}
		paymentsArray = append(paymentsArray, paymentData)
		DecodeBulk(consume, paymentsArray)
	}
	return paymentsArray
}

func decode(consume func(int) []byte) (*PaymentData, error) {
	flagData := consume(1)
	if len(flagData) == 0 {
		return nil, fmt.Errorf("no flags are found")
	}

	flagsBuffer := flagData[0]
	output := []byte{flagsBuffer & ^flagsBuffer}
	flags := flagsBuffer & flagMask

	skipB := flags & skipFlag
	rangeB := flags & rangeFlag
	percentB := flags & percentFlag

	skip := byteToBool([]byte{skipB})
	rangeF := byteToBool([]byte{rangeB})
	percent := byteToBool([]byte{percentB})

	if rangeF {
		output = append(output, consume(1)...)
	}

	outputInt, err := strconv.ParseInt(hex.EncodeToString(output), 16, 64)
	if err != nil {
		return nil, err
	}

	amount, err := sffc.Decode(consume)
	if err != nil {
		return nil, err
	}

	return &PaymentData{
		Skip:    skip,
		Range:   rangeF,
		Percent: percent,
		Output:  outputInt,
		Amount:  amount,
	}, nil
}

func byteToBool(data []byte) bool {
	res := make([]bool, len(data)*8)
	for i := range res {
		res[i] = data[i/8]&(0x80>>byte(i&0x7)) != 0
	}
	return res[0]
}
