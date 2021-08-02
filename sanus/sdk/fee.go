package sdk

import (
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/wallet/txrules"
)

type FeeLevel int

const (
	PRIORITY_FEE FeeLevel = 0
	NORMAL_FEE            = 1
	ECONOMIC_FEE          = 2
	FEE_BUMP_FEE          = 3
)

func (w *BTCWallet) estimateFee(feeLevel FeeLevel) btcutil.Amount {
	nBlocks := 6 // Default fee level - ECONOMIC
	switch feeLevel {
	case NORMAL_FEE:
		nBlocks = 3
	case PRIORITY_FEE:
		nBlocks = 1
	}

	estimatedFee, err := w.rpcClient.EstimateFee(int64(nBlocks))
	if err != nil || estimatedFee <= txrules.DefaultRelayFeePerKb.ToBTC() {
		return txrules.DefaultRelayFeePerKb
	}

	feeAmount, err := btcutil.NewAmount(estimatedFee)
	if err != nil {
		return txrules.DefaultRelayFeePerKb
	}

	return feeAmount
}
