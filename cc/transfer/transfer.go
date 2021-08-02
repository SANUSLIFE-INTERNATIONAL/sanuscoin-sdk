package transfer

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"sanus/sanus-sdk/cc/encdec"
	"sanus/sanus-sdk/cc/utils"
)

type MultiSigData struct {
	Index    int
	HashType string
}

type ColoredData struct {
	Type        string
	Protocol    int64
	Version     int64
	TorrentHash []byte
	Sha2        []byte
	NoRules     bool
	MultiSig    []MultiSigData
	Payments    []*utils.PaymentData
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
		cData.Payments = encdec.TransferDecodeBulk(consume, nil)
	case "burn":
		encdec.BurnDecodeBulk(consume)
	}
	return cData, err
}

func (cd *ColoredData) Encode(byteSize int) (hash []byte, err error) {
	hash = make([]byte, 0)
	if cd.Payments == nil || len(cd.Payments) == 0 {
		return
	}
	var opCode []byte
	var opCodes [][]byte
	if cd.Type == "burn" {
		opCodes = BurnOPCodes
	} else {
		opCodes = TransferOPCodes
	}

	protocolByteStr := strconv.FormatInt(cd.Protocol, 16)
	protocolByte, err := hex.DecodeString(protocolByteStr)
	if err != nil {
		return
	}
	versionByteStr := strconv.FormatInt(cd.Version, 16)
	versionByte, err := hex.DecodeString(versionByteStr)
	if err != nil {
		return
	}
	transferHeader := append(protocolByte, versionByte...)
	var paymentByte []byte
	if cd.Type == "burn" {
		paymentByte = encdec.EncodeBulk(cd.Payments)
	} else {
		paymentByte = encdec.TransferEncodeBulk(cd.Payments)
	}

	var issueByteSize = len(transferHeader) + len(paymentByte) + 1

	if issueByteSize > byteSize {
		return nil, fmt.Errorf("data code is bigger then the allowed byte size")
	}

	if len(cd.Sha2) == 0 {
		if len(cd.TorrentHash) > 0 {
			if cd.NoRules {
				opCode = opCodes[4]
			} else {
				opCode = opCodes[3]
			}
			if issueByteSize+len(cd.TorrentHash) > byteSize {
				return nil, fmt.Errorf("can't fit Torrent Hash in byte size")
			}
			return utils.BytesConcat(transferHeader, opCode, cd.TorrentHash, paymentByte), nil
		}

		return utils.BytesConcat(transferHeader, opCodes[5], hash, paymentByte), nil
	}

	if len(cd.TorrentHash) == 0 {
		return nil, fmt.Errorf("torrent Hash is missing")
	}

	var leftover = [][]byte{
		cd.TorrentHash,
		cd.Sha2,
	}

	opCode = opCodes[2]
	issueByteSize = issueByteSize + len(cd.TorrentHash)
	if issueByteSize <= byteSize {
		hash = utils.BytesConcat(hash, leftover[0])
		opCode = opCodes[1]
		issueByteSize = issueByteSize + len(cd.Sha2)
	}
	if issueByteSize <= byteSize {
		hash = utils.BytesConcat(hash, leftover[1])
		opCode = opCodes[0]
	}
	return utils.BytesConcat(transferHeader, opCode, hash, paymentByte), nil
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
