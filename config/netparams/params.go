// Copyright (C) 2017-2019 The EVEN Network Developers

package netparams

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/go-errors/errors"
)

type Params chaincfg.Params

var (
	// MainnetParams defines the network parameters for the daemon network.
	MainnetParams = Params{
		Name: "mainnet",
		Net:  wire.MainNet,

		// To encoding magics
		PubKeyHashAddrID: 0x00, // starts with 1
		ScriptHashAddrID: 0x05, // starts with 3
		PrivateKeyID:     0x80, // starts with 5 (uncompressed) or K (compressed)

		// BIP32 hierarchical deterministic extended key magics
		HDPrivateKeyID: [4]byte{0x04, 0x88, 0xad, 0xe4}, // starts with xprv
		HDPublicKeyID:  [4]byte{0x04, 0x88, 0xb2, 0x1e}, // starts with xpub

		// BIP44 coin type used in the hierarchical deterministic path for address generation.
		HDCoinType: 0,
	}

	// Testnet defines the network parameters for the test network is sometimes simply called "testnet".
	TestnetParams = Params{
		Name: "testnet",
		Net:  wire.TestNet3,

		// To encoding magics
		PubKeyHashAddrID: 0x6f, // starts with m or n
		ScriptHashAddrID: 0xc4, // starts with 2
		PrivateKeyID:     0xef, // starts with 9 (uncompressed) or c (compressed)

		// BIP32 hierarchical deterministic extended key magics
		HDPrivateKeyID: [4]byte{0x04, 0x35, 0x83, 0x94}, // starts with tprv
		HDPublicKeyID:  [4]byte{0x04, 0x35, 0x87, 0xcf}, // starts with tpub

		// BIP44 coin type used in the hierarchical deterministic path for address generation.
		HDCoinType: 1,
	}
)

var (
	// ErrDuplicateNet describes an error where the parameters for a Bitcoin
	// network could not be set due to the network already being a standard
	// network or previously-registered into this package.
	ErrDuplicateNet = errors.New("duplicate chain network")

	registeredNets          = make(map[wire.BitcoinNet]struct{})
	publicKeyHashAddrIds    = make(map[byte]struct{})
	scriptHashAddrIds       = make(map[byte]struct{})
	hdPrivateToPublicKeyIds = make(map[[4]byte][]byte)
	bech32SegwitPrefixes    = make(map[string]struct{})
)

// Register registers the network parameters for a Bitcoin network.
// This may error with ErrDuplicateNet if the network is already registered
// (either due to a previous Register call, or the network being one of the default networks).
//
// Network parameters should be registered into this package by a daemon package as early as possible.
// Then, library packages may lookup networks or network parameters based on inputs and work regardless
// of the network being standard or not.
func Register(params *Params) error {
	if _, ok := registeredNets[params.Net]; ok {
		return ErrDuplicateNet
	}
	registeredNets[params.Net] = struct{}{}
	scriptHashAddrIds[params.ScriptHashAddrID] = struct{}{}
	publicKeyHashAddrIds[params.PubKeyHashAddrID] = struct{}{}
	hdPrivateToPublicKeyIds[params.HDPrivateKeyID] = params.HDPublicKeyID[:]

	// A valid Bech32 encoded segwit address always has as prefix
	// the human-readable part for the given net followed by '1'.
	bech32SegwitPrefixes[params.Bech32HRPSegwit+"1"] = struct{}{}

	return nil
}

// mustRegister performs the same function as Register except it panics if there is an error.
// This should only be called from package init functions.
func mustRegister(params *Params) {
	if err := Register(params); err != nil {
		panic("failed to register network: " + err.Error())
	}
}

// IsPubKeyHashAddrID returns whether the id is an identifier known to prefix a
// pay-to-publicKey-hash address on any default or registered network.
// This is used when decoding an address string into a specific address type.
// It is up to the caller to check both this and IsScriptHashAddrID and decide whether an address
// is a publicKey hash address, script hash address, neither, or undeterminable (if both return true).
func IsPubKeyHashAddrID(id byte) bool {
	_, exists := publicKeyHashAddrIds[id]
	return exists
}

// IsScriptHashAddrID returns whether the id is an identifier known to prefix a
// pay-to-script-hash address on any default or registered network.
// This is used when decoding an address string into a specific address type.
// It is up to the caller to check both this and IsPubKeyHashAddrID and decide whether an address
// is a publicKey hash address, script hash address, neither, or undeterminable (if both return true).
func IsScriptHashAddrID(id byte) bool {
	_, exists := scriptHashAddrIds[id]
	return exists
}

func init() {
	// Register all default networks when the package is initialized.
	mustRegister(&MainnetParams)
	mustRegister(&TestnetParams)
}
