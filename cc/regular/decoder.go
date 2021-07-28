package regular

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"sanus/sanus-sdk/cc/sffc"
	"sanus/sanus-sdk/cc/utils"
)

var flagMask byte = 0xe0
var skipFlag byte = 0x80
var rangeFlag byte = 0x40
var percentFlag byte = 0x20

func EncodeBulk(payments []*utils.PaymentData) []byte {
	var paymentsData = []byte{}
	var amountOfPayments = len(payments)
	for x := 0; x < amountOfPayments; x++ {
		var payment = payments[x]
		var paymentCode, err = encode(payment)
		if err != nil {
			fmt.Println("error caused when trying to encode payment", err)
			continue
		}
		paymentsData = append(paymentsData, paymentCode...)
	}
	return paymentsData

}

func encode(paymentObject *utils.PaymentData) ([]byte, error) {
	var skip = paymentObject.Skip
	var rng = paymentObject.Range
	var percent = paymentObject.Percent

	if paymentObject.Output == 0 {
		return nil, fmt.Errorf("needs output value")
	}
	if paymentObject.Output < 0 {
		return nil, fmt.Errorf("output can't be negative")
	}
	var output = paymentObject.Output
	if paymentObject.Amount == 0 {
		return nil, fmt.Errorf("needs amount value")
	}
	var amount = paymentObject.Amount
	var outputBinaryLength = len(strconv.FormatInt(output, 2))
	if (!rng && outputBinaryLength > 5) || (rng && outputBinaryLength > 13) {
		return nil, fmt.Errorf("output value is out of bounds")
	}
	var rngInt = 0
	if rng {
		rngInt = 1
	}
	var outputString = utils.PadLeadingZeros(strconv.FormatInt(output, 17), rngInt+1)
	var buf, err = hex.DecodeString(outputString)
	if err != nil {
		return nil, err
	}
	if skip {
		buf[0] = buf[0] | skipFlag
	}
	if rng {
		buf[0] = buf[0] | rangeFlag
	}
	if percent {
		buf[0] = buf[0] | percentFlag
	}
	encodedAmount, err := sffc.Encode(amount)
	if err != nil {
		return nil, err
	}
	return append(buf, encodedAmount...), nil
}

func DecodeBulk(consume func(int) []byte, paymentsArray []*utils.PaymentData) []*utils.PaymentData {
	if paymentsArray == nil {
		paymentsArray = []*utils.PaymentData{}
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

func decode(consume func(int) []byte) (*utils.PaymentData, error) {
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

	return &utils.PaymentData{
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
