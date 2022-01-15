package issuance

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"sanus/sanus-sdk/cc/encdec"
	issue_flags "sanus/sanus-sdk/cc/issue-flags"
	"sanus/sanus-sdk/cc/sffc"
	"sanus/sanus-sdk/cc/transfer"
	"sanus/sanus-sdk/cc/utils"
)

var opCodes = [][]byte{
	[]byte{0x00}, // wild-card to be defined
	[]byte{0x01}, // All Hashes in OP_RETURN - Pay-to-PubkeyHash
	[]byte{0x02}, // SHA2 in Pay-to-Script-Hash multi-sig output (1 out of 2)
	[]byte{0x03}, // All Hashes in Pay-to-Script-Hash multi-sig outputs (1 out of 3)
	[]byte{0x04}, // Low security issue no SHA2 for torrent data. SHA1 is always inside OP_RETURN in this case.
	[]byte{0x05}, // No rules, no torrent, no meta data ( no one may add rules in the future, anyone can add metadata )
	[]byte{0x06}, // No meta data (anyone can add rules and/or metadata  in the future)

}

type ColoredData struct {
	Amount            int                     `json:"amount"`
	Protocol          int64                   `json:"protocol"`
	Version           int64                   `json:"version"`
	Divisibility      int                     `json:"divisibility"`
	LockStatus        bool                    `json:"lock_status"`
	AggregationPolicy string                  `json:"aggregation_policy"`
	MultiSig          []transfer.MultiSigData `json:"multi_sig"`
	NoRules           bool                    `json:"no_rules"`
	Sha2              []byte                  `json:"sha_2"`
	TorrentHash       []byte                  `json:"torrent_hash"`
	Payments          []*utils.PaymentData    `json:"payments"`
}

func (cd *ColoredData) String() string {
	bytes, err := json.Marshal(cd)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

type Response struct {
	Code     []byte
	Leftover [][]byte
}

func (cd *ColoredData) Encode(byteSize int) ([]byte, error) {
	hash := make([]byte, 0)
	opCode := []byte{}
	var protocolString = utils.PadLeadingZeros(strconv.FormatInt(cd.Protocol, 16), 2)
	var protocol, err = hex.DecodeString(protocolString)
	if err != nil {
		return hash, err
	}
	version := []byte{byte(cd.Version)}
	issueHeader := utils.BytesConcat(protocol, version)
	amount, err := sffc.Encode(cd.Amount)
	if err != nil {
		return hash, err
	}
	payments := encdec.TransferEncodeBulk(cd.Payments)
	flg := issue_flags.Flags{
		Divisibility:      cd.Divisibility,
		LockStatus:        cd.LockStatus,
		AggregationPolicy: cd.AggregationPolicy,
	}
	issueFlagsByte, err := flg.Encode()
	if err != nil {
		return hash, err
	}
	issueTail := utils.BytesConcat(amount, payments, issueFlagsByte)
	issueByteSize := len(issueHeader) + len(issueTail) + 1
	if issueByteSize > byteSize {
		return hash, fmt.Errorf("data code is bigger then the allowed byte size")
	}
	if len(cd.Sha2) == 0 {
		if len(cd.TorrentHash) > 0 {
			if issueByteSize+len(cd.TorrentHash) > byteSize {
				return hash, fmt.Errorf("can't fit Torrent Hash in byte size")
			}
			return utils.BytesConcat(issueHeader, opCodes[4], cd.TorrentHash, issueTail), nil
		}
		opCode = opCodes[5]
		if !cd.NoRules {
			opCode = opCodes[6]
		}
		return utils.BytesConcat(issueHeader, opCode, hash, issueTail), nil
	}
	if len(cd.TorrentHash) == 0 {
		return hash, fmt.Errorf("torrent Hash is missing")
	}
	opCode = opCodes[3]
	issueByteSize = issueByteSize + len(cd.TorrentHash)
	if issueByteSize <= byteSize {
		hash = utils.BytesConcat(hash, cd.TorrentHash)
		opCode = opCodes[2]
		issueByteSize = issueByteSize + len(cd.Sha2)
	}
	if issueByteSize <= byteSize {
		hash = utils.BytesConcat(hash, cd.Sha2)
		opCode = opCodes[1]
	}

	return utils.BytesConcat(issueHeader, opCode, hash, issueTail), nil
}

func Decode(hexByte []byte) (data *ColoredData, err error) {
	var byteSize = len(hexByte)
	var lastByte = []byte{hexByte[len(hexByte)-1]}
	var issueTail = issue_flags.Decode(consumer(lastByte))
	data = &ColoredData{}
	data.Divisibility = issueTail.Divisibility
	data.LockStatus = issueTail.LockStatus
	data.AggregationPolicy = issueTail.AggregationPolicy
	var consume = consumer(hexByte[0 : byteSize-1])
	protocolStr := hex.EncodeToString(consume(2))
	data.Protocol, err = strconv.ParseInt(protocolStr, 16, 64)
	if err != nil {
		return nil, err
	}
	versionStr := hex.EncodeToString(consume(1))
	data.Version, err = strconv.ParseInt(versionStr, 16, 64)
	if err != nil {
		return nil, err
	}
	data.MultiSig = []transfer.MultiSigData{}

	var opcode = consume(1)
	if len(opcode) == 0 {
		return nil, fmt.Errorf("missing data to get opcode")
	}
	if opcode[0] == opCodes[1][0] {
		data.TorrentHash = consume(20)
		data.Sha2 = consume(32)
	} else if opcode[0] == opCodes[2][0] {
		data.TorrentHash = consume(20)
		data.MultiSig = append(data.MultiSig, transfer.MultiSigData{Index: 1, HashType: "sha2"})
	} else if opcode[0] == opCodes[3][0] {
		data.MultiSig = append(data.MultiSig, transfer.MultiSigData{Index: 1, HashType: "sha2"})
		data.MultiSig = append(data.MultiSig, transfer.MultiSigData{Index: 2, HashType: "torrentHash"})
	} else if opcode[0] == opCodes[4][0] {
		data.TorrentHash = consume(20)
	} else if opcode[0] == opCodes[5][0] {
		data.NoRules = true
	} else if opcode[0] == opCodes[6][0] {
		data.NoRules = false
	} else {
		return nil, fmt.Errorf("unrecognized code")
	}
	data.Amount, err = decodeAmountByVersion(data.Version, consume, data.Divisibility)
	if err != nil {
		return nil, err
	}
	data.Payments = encdec.TransferDecodeBulk(consume, nil)

	return data, nil
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

func decodeAmountByVersion(version int64, consume func(int) []byte, divisibility int) (int, error) {
	var decodedAmount, err = sffc.Decode(consume)
	if err != nil {
		return 0, err
	}
	if byte(version) == 0x01 {
		return int(float64(decodedAmount) / math.Pow(10, float64(divisibility))), nil
	}
	return decodedAmount, nil
}
