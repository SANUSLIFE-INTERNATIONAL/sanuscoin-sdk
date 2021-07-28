package cc

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"sanus/sanus-sdk/cc/burn"
	"sanus/sanus-sdk/cc/regular"
)

type MultiSigData struct {
	Index    int
	HashType string
}



func Decode(data []byte) (*ColoredData, error) {
	var (
		err error

		cData = &ColoredData{}

		consume = consumer(data[2:])
	)

	protocolStr := hex.EncodeToString(consume(2))
	cData.Protocol, err = strconv.ParseInt(protocolStr, 16, 64)
	if err != nil {
		return nil, err
	}
	versionStr := hex.EncodeToString(consume(1))
	cData.Version, err = strconv.ParseInt(versionStr, 16, 64)
	if err != nil {
		return nil, err
	}
	var paymentEncoder = ""
	opcode := consume(1)
	if (opcode[0] & TypeMask) == TransferMask {
		paymentEncoder = "transfer"
	} else if (opcode[0] & TypeMask) == BurnMask {
		paymentEncoder = "burn"
	} else {
		return nil, fmt.Errorf("unrecognized code")
	}
	if opcode[0] == TransferOPCodes[0][0] || opcode[0] == BurnOPCodes[0][0] {
		cData.TorrentHash = consume(20)
		cData.Sha2 = consume(32)
	} else if opcode[0] == TransferOPCodes[1][0] || opcode[0] == BurnOPCodes[1][0] {
		cData.TorrentHash = consume(20)
		cData.MultiSig = append(cData.MultiSig, MultiSigData{
			Index:    1,
			HashType: "sha2",
		})
	} else if opcode[0] == TransferOPCodes[2][0] || opcode[0] == BurnOPCodes[2][0] {
		cData.MultiSig = append(cData.MultiSig, MultiSigData{
			Index:    1,
			HashType: "sha2",
		})
		cData.MultiSig = append(cData.MultiSig, MultiSigData{
			Index:    2,
			HashType: "torrentHash",
		})
	} else if opcode[0] == TransferOPCodes[3][0] || opcode[0] == BurnOPCodes[3][0] {
		cData.TorrentHash = consume(20)
		cData.NoRules = false
	} else if opcode[0] == TransferOPCodes[4][0] || opcode[0] == BurnOPCodes[4][0] {
		cData.TorrentHash = consume(20)
		cData.NoRules = true
	} else if opcode[0] == TransferOPCodes[5][0] || opcode[0] == BurnOPCodes[5][0] {
	} else {
		return nil, fmt.Errorf("unrecognized code")
	}

	switch paymentEncoder {
	case "transfer":
		cData.Payments = regular.DecodeBulk(consume, nil)
	case "burn":
		burn.DecodeBulk(consume)
	}
	return cData, err
}

func consumer(buff []byte) func(int) []byte {
	var curr = 0
	return func(lgt int) []byte {
		if len(buff) < lgt+curr {
			return []byte{}
		}
		bytes := buff[curr : curr+lgt]
		curr = curr + lgt
		return bytes
	}
}
