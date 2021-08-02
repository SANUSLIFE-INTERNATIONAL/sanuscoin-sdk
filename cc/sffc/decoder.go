package sffc

import (
	"encoding/hex"
	"math"
	"strconv"
)

type SchemeTable struct {
	Flag     byte
	Exponent int
	ByteSize int
	Mantis   int
}

var (
	flagMask byte = 0xe0

	encodingSchemeTable = []SchemeTable{
		{
			Flag:     0x20,
			Exponent: 4,
			ByteSize: 2,
			Mantis:   9,
		},
		{
			Flag:     0x40,
			Exponent: 4,
			ByteSize: 3,
			Mantis:   17,
		},
		{
			Flag:     0x60,
			Exponent: 4,
			ByteSize: 4,
			Mantis:   25,
		},
		{
			Flag:     0x80,
			Exponent: 3,
			ByteSize: 5,
			Mantis:   34,
		},
		{
			Flag:     0xa0,
			Exponent: 3,
			ByteSize: 6,
			Mantis:   42,
		},
		{
			Flag:     0xc0,
			Exponent: 0,
			ByteSize: 7,
			Mantis:   54,
		},
	}

	flagLookup   = map[byte]SchemeTable{}
	mantisLookup = map[int]SchemeTable{}
)

const (
	maxSafeInt = 4294967295
)

func init() {
	for i, _ := range encodingSchemeTable {
		var flagObject = encodingSchemeTable[i]
		flagLookup[flagObject.Flag] = flagObject
	}
	var currentIndex = 0
	var currentMantis = encodingSchemeTable[currentIndex].Mantis
	var endMantis = encodingSchemeTable[len(encodingSchemeTable)-1].Mantis

	for i := 1; i <= endMantis; i++ {
		if i > currentMantis {
			currentIndex++
			currentMantis = encodingSchemeTable[currentIndex].Mantis
		}
		mantisLookup[i] = encodingSchemeTable[currentIndex]
	}
}

func Decode(consume func(int) []byte) (int, error) {
	var flagByte = consume(1)[0]
	var flag = flagByte & flagMask
	if flag == 0 {
		return int(flagByte), nil
	}
	if flag == 0xe0 {
		flag = 0xc0
	}
	var encodingObject = flagLookup[flag]
	var headOfNumber = []byte{flagByte & (^flag)}
	var tailOfNumber = consume(encodingObject.ByteSize - 1)
	var fullNumber = append(headOfNumber, tailOfNumber...)
	var number, err = strconv.ParseInt(hex.EncodeToString(fullNumber), 16, 64)
	if err != nil {
		return 0, err
	}
	var exponentShift = math.Pow(2, float64(encodingObject.Exponent))
	var exponent = int(number) % int(exponentShift)
	var mantis = math.Floor(float64(number) / exponentShift)
	return int(mantis * math.Pow(10, float64(exponent))), nil
}
