package sdk

import (
	"github.com/tyler-smith/go-bip39"
)

// NewSeed method generates seed phrase using mnemonic and public password
func NewSeed(mnemonic, pass string) []byte {
	return bip39.NewSeed(mnemonic, pass)
}
