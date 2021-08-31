package utils

type PaymentData struct {
	Skip    bool  `json:"skip"`
	Range   bool  `json:"range"`
	Percent bool  `json:"percent"`
	Output  int64 `json:"output"`
	Amount  int   `json:"amount"`
}

func PadLeadingZeros(hex string, byteSize int) string {
	if len(hex) == byteSize*2 {
		return hex
	}
	return PadLeadingZeros("0"+hex, byteSize)
}

func ByteToBool(data []byte) bool {
	res := make([]bool, len(data)*8)
	for i := range res {
		res[i] = data[i/8]&(0x80>>byte(i&0x7)) != 0
	}
	return res[0]
}

func BytesConcat(slices ...[]byte) []byte {
	var hash = make([]byte, 0)
	for x := 0; x < len(slices); x++ {
		hash = append(hash, slices[x]...)
	}
	return hash
}
