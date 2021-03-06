package issue_flags

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"sanus/sanus-sdk/cc/utils"
)

var aggregationPolicies = [3]string{
	"aggregatable",
	"hybrid",
	"dispersed",
}

type Flags struct {
	Divisibility      int
	LockStatus        bool
	AggregationPolicy string
}

func (f *Flags) Encode() ([]byte, error) {
	var divisibility = 0
	if f.Divisibility != 0 {
		divisibility = f.Divisibility
	}
	var lockStatus = f.LockStatus
	var aggregationPolicy = aggregationPolicies[0]
	if f.AggregationPolicy != "" {
		aggregationPolicy = f.AggregationPolicy
	}
	if divisibility < 0 || divisibility > 7 {
		return nil, fmt.Errorf("divisibility not in range")
	}

	var aggregationPolicyIndex = -1

	for x := 0; x < len(aggregationPolicies); x++ {
		if aggregationPolicy == aggregationPolicies[x] {
			aggregationPolicyIndex = x
		}
	}

	if aggregationPolicyIndex == -1 {
		return nil, fmt.Errorf("invalid aggregation policy")
	}

	var result = divisibility << 1
	var lockStatusFlag = 0
	if lockStatus {
		lockStatusFlag = 1
	}
	result = result | lockStatusFlag
	result = result << 2
	result = result | aggregationPolicyIndex
	result = result << 2

	var resultString = strconv.FormatInt(int64(result), 16)
	resultString = utils.PadLeadingZeros(resultString, 1)
	return hex.DecodeString(resultString)
}

func Decode(consume func(int) []byte) *Flags {
	consumeFirst := consume(1)
	if len(consumeFirst) == 0 {
		return nil
	}
	var number = consumeFirst[0]
	number = number >> 2 // least significant 2 bits unused
	var aggregationPolicy string
	if int(number&0x3) > len(aggregationPolicies)-1 {
		aggregationPolicy = "undefined"
	} else {
		aggregationPolicy = aggregationPolicies[number&0x3]
	}

	number = number >> 2

	var lockStatus = !(number&1 == 0)
	number = number >> 1
	var divisibility = number & 0x7
	return &Flags{
		Divisibility:      int(divisibility),
		LockStatus:        lockStatus,
		AggregationPolicy: aggregationPolicy,
	}
}
