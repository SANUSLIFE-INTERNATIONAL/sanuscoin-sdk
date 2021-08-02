package sdk

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
)

type coinBase struct {
	TxHash       *chainhash.Hash
	TxIndex      uint32
	TxValue      btcutil.Amount
	TxNumConfs   int64
	ScriptPubKey []byte
}

func (c *coinBase) Hash() *chainhash.Hash {
	return c.TxHash
}

func (c *coinBase) Index() uint32 {
	return c.TxIndex
}

func (c *coinBase) Value() btcutil.Amount {
	return c.TxValue
}

func (c *coinBase) PkScript() []byte {
	return c.ScriptPubKey
}

func (c *coinBase) NumConfs() int64 {
	return c.TxNumConfs
}

func (c *coinBase) ValueAge() int64 {
	return int64(c.TxValue) * c.TxNumConfs
}
