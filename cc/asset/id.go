package asset

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"

	"sanus/sanus-sdk/cc/utils"

	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

func createIdFromPreviousOutputScriptPubKey(payload string, padding, divisibility int) (string, error) {
	hexHash, err := hex.DecodeString(payload)
	if err != nil {
		return "", err
	}
	hexStr := hex.EncodeToString(hexHash)
	return hashAndBase58CheckEncode(hexStr, padding, divisibility)

}

func createIdFromTxidIndex(vout wire.OutPoint, padding int, divisibility int) (string, error) {
	indexToString := strconv.Itoa(int(vout.Index))
	var str = vout.Hash.String() + ":" + indexToString
	return hashAndBase58CheckEncode(str, padding, divisibility)
}

func hashAndBase58CheckEncode(payloadToHash string, padding int, divisibility int) (string, error) {
	var hash256 = sha256.New()
	hash256.Write([]byte(payloadToHash))
	hash160 := ripemd160.New()
	hash160.Write(hash256.Sum(nil))
	ripemdHash := hash160.Sum(nil)
	paddingStrFormat := strconv.FormatInt(int64(padding), 16)
	paddingPadLeading := utils.PadLeadingZeros(paddingStrFormat, -1)
	pdd, err := hex.DecodeString(paddingPadLeading)
	if err != nil {
		return "", nil
	}

	divisibilityStrFormat := strconv.FormatInt(int64(divisibility), 16)
	divisibilityPadLeading := utils.PadLeadingZeros(divisibilityStrFormat, 2)
	dvb, err := hex.DecodeString(divisibilityPadLeading)
	if err != nil {
		return "", nil
	}
	hash := append([]byte{}, pdd...)
	hash = append(hash, ripemdHash...)
	hash = append(hash, dvb...)

	return bs58check(hash), nil
}

func bs58check(payload []byte) string {
	hash256F := sha256.New()
	hash256F.Write(payload)
	hash256S := sha256.New()
	hash256S.Write(hash256F.Sum(nil))
	checksum := hash256S.Sum(nil)
	concat := append(payload, checksum...)
	return base58.Encode(concat[:len(payload)+4])
}

func assetId(tx *CCTransaction) (string, error) {
	if tx.Issuance == nil {
		return "", fmt.Errorf("missing Colored Coin metadata")
	}
	if tx.Type != "issuance" {
		return "", fmt.Errorf("not an issuance transaction")
	}
	var lockStatus = tx.Issuance.LockStatus
	var aggregationPolicy = "aggregatable"
	if tx.Issuance.AggregationPolicy != "" {
		aggregationPolicy = tx.Issuance.AggregationPolicy
	}
	var divisibility = tx.Issuance.Divisibility
	var firstInput = tx.Input[0]
	var padding int

	if lockStatus {
		padding = lockPadding[aggregationPolicy]
		return createIdFromTxidIndex(
			firstInput.GetInput().PreviousOutPoint, padding, divisibility)
	}

	padding = unlockPadding[aggregationPolicy]

	return createIdFromPreviousOutputScriptPubKey(firstInput.GetInput().PreviousOutPoint.Hash.String(), padding, divisibility)

}
