package sffc

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"

	"sanus/sanus-sdk/cc/utils"
)

func intToFloatArray(number int, n interface{}) []int {
	nInt, ok := n.(int)
	if !ok {
		nInt = 0
	}

	if number%10 > 0 {
		return []int{number, nInt}
	}
	return intToFloatArray(number/10, nInt+1)
}

func Encode(amount int) ([]byte, error) {
	if amount < 0 {
		return []byte{}, fmt.Errorf("number is out of bounds")
	}
	if amount > maxSafeInt {
		return []byte{}, fmt.Errorf("number is out of bounds")
	}
	if amount < 32 {
		return []byte{byte(amount)}, nil
	}
	var floatingNumberArray = intToFloatArray(amount, nil)
	var encodingObjFormatInt = strconv.FormatInt(int64(floatingNumberArray[0]), 2)
	var encodingObject, ok = mantisLookup[len(encodingObjFormatInt)]
	if !ok {
		return []byte{}, fmt.Errorf("number is out of bounds")
	}
	for true {
		if math.Pow(2, float64(encodingObject.Exponent))-1 > float64(floatingNumberArray[1]) {
			break
		}
		floatingNumberArray[0] = floatingNumberArray[0] * 10
		floatingNumberArray[1] = floatingNumberArray[1] - 1
	}
	var shiftedNumber = float64(floatingNumberArray[0]) *
		math.Pow(2, float64(encodingObject.Exponent))

	var shiftedNumberStr = strconv.FormatInt(int64(shiftedNumber), 16)

	var numberString = utils.PadLeadingZeros(shiftedNumberStr, encodingObject.ByteSize)
	buf, err := hex.DecodeString(numberString)
	if err != nil {
		return nil, err
	}

	buf[0] = buf[0] | encodingObject.Flag
	buf[len(buf)-1] = buf[len(buf)-1] | byte(floatingNumberArray[1])

	return buf, nil
}
