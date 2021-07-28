package utils

type PaymentData struct {
	Skip    bool
	Range   bool
	Percent bool
	Output  int64
	Amount  int
}

func PadLeadingZeros(hex string, byteSize int) string {
	if len(hex) == byteSize*2 {
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
