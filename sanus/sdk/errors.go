package sdk

import (
	"fmt"
)

var (
	ErrDbDoesNotExist       = fmt.Errorf("database does not exist")
	ErrorDustAmount         = fmt.Errorf("amount is below network dust threshold")
	ErrorDustAmountChange   = fmt.Errorf("amount change is below network dust threshold")
	ErrorInsufficientFunds  = fmt.Errorf("insufficient funds in wallet")
	ErrorWalletDoesNotExist = fmt.Errorf("wallet does not exist")
	ErrorWalletLocked       = fmt.Errorf("wallet is locked")
	ErrKeyNotFound          = fmt.Errorf("key not found")
	ErrPrevOutScript        = fmt.Errorf("prev-out script not found")
	ErrFailedToSignTx       = fmt.Errorf("failed to sign the transaction")
)
