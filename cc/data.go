package cc

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"sanus/sanus-sdk/cc/burn"
	issue_flags "sanus/sanus-sdk/cc/issue-flags"
	"sanus/sanus-sdk/cc/regular"
	"sanus/sanus-sdk/cc/sffc"
	"sanus/sanus-sdk/cc/utils"
)

type ColoredData struct {
	Amount int
	*issue_flags.Flags
	Type        string
	Protocol    int64
	Version     int64
	TorrentHash []byte
	Sha2        []byte
	NoRules     bool
	MultiSig    []MultiSigData
	Payments    []*utils.PaymentData
}

type EncodedColoredData struct {
	Data     []byte
	LeftOver [][]byte
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
		paymentByte = burn.EncodeBulk(cd.Payments)
	} else {
		paymentByte = regular.EncodeBulk(cd.Payments)
	}
	issueFlagsByte, err := cd.Flags.Encode()
	if err != nil {
		return nil, err
	}
	amountByte, err := sffc.Encode(cd.Amount)
	if err != nil {
		return nil, err
	}
	var issueTail = utils.BytesConcat(amountByte, paymentByte, issueFlagsByte)

	var issueHeader = utils.BytesConcat(protocolByte, versionByte)

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
		if cd.NoRules {
			opCode = opCodes[5]
		} else {
			opCode = opCodes[6]
		}
		return utils.BytesConcat(issueHeader, opCode, issueTail), nil
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
		hash = leftover[0]
		issueByteSize = issueByteSize + len(cd.Sha2)
	}
	if issueByteSize <= byteSize {
		hash = leftover[1]
		opCode = opCodes[1]
	}
	hash = utils.BytesConcat(issueHeader, opCode, hash, issueTail)
	return
}
