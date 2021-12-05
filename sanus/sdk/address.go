package sdk

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"sanus/sanus-sdk/cc/transfer"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/waddrmgr"
	btcWallet "github.com/btcsuite/btcwallet/wallet"
	"github.com/btcsuite/btcwallet/walletdb"
)

var (
	scopeBucketName       = []byte("scope")
	addrBucketName        = []byte("addr")
	addrAcctIdxBucketName = []byte("addracctidx")
	nullVal               = []byte{0}
)

const scopeKeySize = 8

// dbAddressRow houses common information stored about an address in the
// database.
type dbAddressRow struct {
	addrType   uint8
	account    uint32
	addTime    uint64
	syncStatus uint8
	rawData    []byte // Varies based on address type field.
}

// managedAddress represents a public key address.  It also may or may not have
// the private key associated with the public key.
type managedAddress struct {
	manager          *waddrmgr.ScopedKeyManager
	derivationPath   waddrmgr.DerivationPath
	address          btcutil.Address
	imported         bool
	internal         bool
	compressed       bool
	used             bool
	addrType         waddrmgr.AddressType
	pubKey           *btcec.PublicKey
	privKeyEncrypted []byte
	privKeyCT        []byte // non-nil if unlocked
	privKeyMutex     sync.Mutex
}

func (w *BTCWallet) List() ([]string, error) {
	var addrList []string
	if accounts, err := w.wlt.Accounts(waddrmgr.KeyScopeBIP0044); err == nil {
		for _, a := range accounts.Accounts {
			if addresses, err := w.wlt.AccountAddresses(a.AccountNumber); err == nil {
				for _, addr := range addresses {
					addrList = append(addrList, addr.String())
				}
			}
			return addrList, err
		}
	} else {
		return addrList, err
	}
	return addrList, nil
}

// NewAddress method generates a new BIP44 address
func (w *BTCWallet) NewAddress() (btcutil.Address, error) {
	addr, err := w.wlt.NewAddress(0, waddrmgr.KeyScopeBIP0044)
	if err != nil {
		return nil, err
	}
	w.rescan(addr)
	return addr, nil
}

// ImportAddress method imports a new address based on privat key
func (w *BTCWallet) ImportAddress(publicKey string) (btcutil.Address, error) {
	publickKeyHash, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	manager, err := w.wlt.Manager.FetchScopedKeyManager(waddrmgr.KeyScopeBIP0044)
	if err != nil {
		return nil, err
	}

	pk, err := btcec.ParsePubKey(publickKeyHash, btcec.S256())
	if err != nil {
		return nil, err
	}

	var addr btcutil.Address
	var props *waddrmgr.AccountProperties
	err = walletdb.Update(w.wlt.Database(), func(tx walletdb.ReadWriteTx) error {
		addrmgrNs := tx.ReadWriteBucket([]byte("waddrmgr"))

		// Prevent duplicates.
		serializedPubKey := pk.SerializeCompressed()
		pubKeyHash := btcutil.Hash160(serializedPubKey)
		alreadyExists := w.existsAddress(addrmgrNs, pubKeyHash)
		if alreadyExists {
			return fmt.Errorf("already exists")
		}

		encryptedPubKey, err := w.wlt.Manager.Encrypt(waddrmgr.CKTPublic, serializedPubKey)
		if err != nil {
			return err
		}
		encryptedPrivKey := []byte{}
		ks := manager.Scope()
		w.putImportedAddress(
			addrmgrNs, &ks, pubKeyHash, 0, 0,
			encryptedPubKey, encryptedPrivKey,
		)

		// The full derivation path for an imported key is incomplete as we
		// don't know exactly how it was derived.
		importedDerivationPath := waddrmgr.DerivationPath{
			Account: 0,
		}

		managedAddr, err := newManagedAddressWithoutPrivKey(
			manager, importedDerivationPath, pk, true,
			manager.AddrSchema().ExternalAddrType,
		)
		addr = managedAddr.address
		//////////////////
		props, err = manager.AccountProperties(
			addrmgrNs, waddrmgr.ImportedAddrAccount,
		)
		return err
	})

	w.rescan(addr)

	if err := w.wlt.ChainClient().NotifyReceived([]btcutil.Address{addr}); err != nil {
		return nil, fmt.Errorf("failed to subscribe for address ntfns for "+
			"address %s: %s", addr.EncodeAddress(), err)
	}
	go func() {
		//an := w.wlt.NtfnServer.AccountNotifications()
		//an.C <- &btcWallet.AccountNotification{
		//	AccountNumber:    props.AccountNumber,
		//	AccountName:      props.AccountName,
		//	ExternalKeyCount: props.ExternalKeyCount,
		//	InternalKeyCount: props.InternalKeyCount,
		//	ImportedKeyCount: props.ImportedKeyCount,
		//}
	}()
	return addr, nil
}

// newManagedAddressWithoutPrivKey returns a new managed address based on the
// passed account, public key, and whether or not the public key should be
// compressed.
func newManagedAddressWithoutPrivKey(m *waddrmgr.ScopedKeyManager,
	derivationPath waddrmgr.DerivationPath, pubKey *btcec.PublicKey, compressed bool,
	addrType waddrmgr.AddressType) (*managedAddress, error) {

	// Create a pay-to-pubkey-hash address from the public key.
	var pubKeyHash []byte
	if compressed {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeCompressed())
	} else {
		pubKeyHash = btcutil.Hash160(pubKey.SerializeUncompressed())
	}

	var address btcutil.Address
	var err error

	switch addrType {

	case waddrmgr.NestedWitnessPubKey:
		// For this address type we'l generate an address which is
		// backwards compatible to Bitcoin nodes running 0.6.0 onwards, but
		// allows us to take advantage of segwit's scripting improvments,
		// and malleability fixes.

		// First, we'll generate a normal p2wkh address from the pubkey hash.
		witAddr, err := btcutil.NewAddressWitnessPubKeyHash(
			pubKeyHash, m.ChainParams(),
		)
		if err != nil {
			return nil, err
		}

		// Next we'll generate the witness program which can be used as a
		// pkScript to pay to this generated address.
		witnessProgram, err := txscript.PayToAddrScript(witAddr)
		if err != nil {
			return nil, err
		}

		// Finally, we'll use the witness program itself as the pre-image
		// to a p2sh address. In order to spend, we first use the
		// witnessProgram as the sigScript, then present the proper
		// <sig, pubkey> pair as the witness.
		address, err = btcutil.NewAddressScriptHash(
			witnessProgram, m.ChainParams(),
		)
		if err != nil {
			return nil, err
		}

	case waddrmgr.PubKeyHash:
		address, err = btcutil.NewAddressPubKeyHash(
			pubKeyHash, m.ChainParams(),
		)
		if err != nil {
			return nil, err
		}

	case waddrmgr.WitnessPubKey:
		address, err = btcutil.NewAddressWitnessPubKeyHash(
			pubKeyHash, m.ChainParams(),
		)
		if err != nil {
			return nil, err
		}
	}

	return &managedAddress{
		manager:          m,
		address:          address,
		derivationPath:   derivationPath,
		imported:         false,
		internal:         false,
		addrType:         addrType,
		compressed:       compressed,
		pubKey:           pubKey,
		privKeyEncrypted: nil,
		privKeyCT:        nil,
	}, nil
}

func (w *BTCWallet) existsAddress(bucket walletdb.ReadWriteBucket, hash []byte) bool {
	return false
}

func (w *BTCWallet) putImportedAddress(ns walletdb.ReadWriteBucket, scope *waddrmgr.KeyScope,
	addressID []byte, account uint32, status uint8,
	encryptedPubKey, encryptedPrivKey []byte) error {

	rawData := serializeImportedAddress(encryptedPubKey, encryptedPrivKey)
	addrRow := dbAddressRow{
		addrType:   1,
		account:    account,
		addTime:    uint64(time.Now().Unix()),
		syncStatus: status,
		rawData:    rawData,
	}

	return putAddress(ns, scope, addressID, &addrRow)
}

// putAddress stores the provided address information to the database.  This is
// used a common base for storing the various address types.
func putAddress(ns walletdb.ReadWriteBucket, scope *waddrmgr.KeyScope,
	addressID []byte, row *dbAddressRow) error {

	scopedBucket, err := fetchWriteScopeBucket(ns, scope)
	if err != nil {
		return err
	}

	bucket := scopedBucket.NestedReadWriteBucket(addrBucketName)

	// Write the serialized value keyed by the hash of the address.  The
	// additional hash is used to conceal the actual address while still
	// allowed keyed lookups.
	addrHash := sha256.Sum256(addressID)
	err = bucket.Put(addrHash[:], serializeAddressRow(row))
	if err != nil {

		return fmt.Errorf("failed to store address %x", addressID)
	}
	// Update address account index
	return putAddrAccountIndex(ns, scope, row.account, addrHash[:])
}

// putAddrAccountIndex stores the given key to the address account index of the
// database.
func putAddrAccountIndex(ns walletdb.ReadWriteBucket, scope *waddrmgr.KeyScope,
	account uint32, addrHash []byte) error {

	scopedBucket, err := fetchWriteScopeBucket(ns, scope)
	if err != nil {
		return err
	}

	bucket := scopedBucket.NestedReadWriteBucket(addrAcctIdxBucketName)

	// Write account keyed by address hash
	err = bucket.Put(addrHash, uint32ToBytes(account))
	if err != nil {
		return nil
	}

	bucket, err = bucket.CreateBucketIfNotExists(uint32ToBytes(account))
	if err != nil {
		return err
	}

	// In account bucket, write a null value keyed by the address hash
	err = bucket.Put(addrHash, nullVal)
	if err != nil {
		return fmt.Errorf("failed to store address account index key %s", addrHash)
	}
	return nil
}

// uint32ToBytes converts a 32 bit unsigned integer into a 4-byte slice in
// little-endian order: 1 -> [1 0 0 0].
func uint32ToBytes(number uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, number)
	return buf
}

// deserializeAddressRow deserializes the passed serialized address
// information.  This is used as a common base for the various address types to
// deserialize the common parts.
func deserializeAddressRow(serializedAddress []byte) (*dbAddressRow, error) {
	// The serialized address format is:
	//   <addrType><account><addedTime><syncStatus><rawdata>
	//
	// 1 byte addrType + 4 bytes account + 8 bytes addTime + 1 byte
	// syncStatus + 4 bytes raw data length + raw data

	// Given the above, the length of the entry must be at a minimum
	// the constant value sizes.
	if len(serializedAddress) < 18 {
		return nil, fmt.Errorf("malformed serialized address")
	}

	row := dbAddressRow{}
	row.addrType = serializedAddress[0]
	row.account = binary.LittleEndian.Uint32(serializedAddress[1:5])
	row.addTime = binary.LittleEndian.Uint64(serializedAddress[5:13])
	row.syncStatus = serializedAddress[13]
	rdlen := binary.LittleEndian.Uint32(serializedAddress[14:18])
	row.rawData = make([]byte, rdlen)
	copy(row.rawData, serializedAddress[18:18+rdlen])

	return &row, nil
}

// serializeAddressRow returns the serialization of the passed address row.
func serializeAddressRow(row *dbAddressRow) []byte {
	// The serialized address format is:
	//   <addrType><account><addedTime><syncStatus><commentlen><comment>
	//   <rawdata>
	//
	// 1 byte addrType + 4 bytes account + 8 bytes addTime + 1 byte
	// syncStatus + 4 bytes raw data length + raw data
	rdlen := len(row.rawData)
	buf := make([]byte, 18+rdlen)
	buf[0] = byte(row.addrType)
	binary.LittleEndian.PutUint32(buf[1:5], row.account)
	binary.LittleEndian.PutUint64(buf[5:13], row.addTime)
	buf[13] = byte(row.syncStatus)
	binary.LittleEndian.PutUint32(buf[14:18], uint32(rdlen))
	copy(buf[18:18+rdlen], row.rawData)
	return buf
}

func fetchWriteScopeBucket(ns walletdb.ReadWriteBucket,
	scope *waddrmgr.KeyScope) (walletdb.ReadWriteBucket, error) {

	rootScopeBucket := ns.NestedReadWriteBucket(scopeBucketName)

	scopeKey := scopeToBytes(scope)
	scopedBucket := rootScopeBucket.NestedReadWriteBucket(scopeKey[:])
	if scopedBucket == nil {

		return nil, fmt.Errorf("unable to find scope %v", scope)
	}

	return scopedBucket, nil
}

// scopeToBytes transforms a manager's scope into the form that will be used to
// retrieve the bucket that all information for a particular scope is stored
// under
func scopeToBytes(scope *waddrmgr.KeyScope) [scopeKeySize]byte {
	var scopeBytes [scopeKeySize]byte
	binary.LittleEndian.PutUint32(scopeBytes[:], scope.Purpose)
	binary.LittleEndian.PutUint32(scopeBytes[4:], scope.Coin)

	return scopeBytes
}

// serializeImportedAddress returns the serialization of the raw data field for
// an imported address.
func serializeImportedAddress(encryptedPubKey, encryptedPrivKey []byte) []byte {
	// The serialized imported address raw data format is:
	//   <encpubkeylen><encpubkey><encprivkeylen><encprivkey>
	//
	// 4 bytes encrypted pubkey len + encrypted pubkey + 4 bytes encrypted
	// privkey len + encrypted privkey
	pubLen := uint32(len(encryptedPubKey))
	privLen := uint32(len(encryptedPrivKey))
	rawData := make([]byte, 8+pubLen+privLen)
	binary.LittleEndian.PutUint32(rawData[0:4], pubLen)
	copy(rawData[4:4+pubLen], encryptedPubKey)
	offset := 4 + pubLen
	binary.LittleEndian.PutUint32(rawData[offset:offset+4], privLen)
	offset += 4
	copy(rawData[offset:offset+privLen], encryptedPrivKey)
	return rawData
}

// get current address for default (number 0) account,
// a new address will be generated if not exists yet
func (w *BTCWallet) defaultAddress() (address btcutil.Address, err error) {
	return w.wlt.CurrentAddress(0, waddrmgr.KeyScopeBIP0044)
}

func (w *BTCWallet) SNCBalance(address btcutil.Address) (int64, error) {
	txs, err := w.wlt.ListUnspent(3, 9999999, map[string]struct{}{
		address.EncodeAddress(): {},
	})

	if err != nil {
		return 0, err
	}

	var balance int64 = 0
	for _, tx := range txs {
		h, err := chainhash.NewHashFromStr(tx.TxID)
		if err != nil {
			return 0, err
		}
		txDetail, err := btcWallet.UnstableAPI(w.wlt).TxDetails(h)
		if err != nil {
			return 0, err
		}

		for _, out := range txDetail.MsgTx.TxOut {
			pkScript := out.PkScript
			if pkScript[0] == txscript.OP_RETURN {
				pkScriptData, err := transfer.Decode(pkScript)
				fmt.Printf("%#v", pkScriptData)
				fmt.Printf("%#v", pkScriptData.Payments[0])
				if err != nil {
					w.Errorf("Error caused when trying to fetch data from PkScript | %v", err)
				}

				for _, p := range pkScriptData.Payments {
					balance += int64(p.Amount)
				}
			}
		}
	}
	return balance, err
}

func (w *BTCWallet) BTCBalance(address btcutil.Address) (float64, error) {
	txs, err := w.wlt.ListUnspent(3, 9999999, map[string]struct{}{
		address.EncodeAddress(): {},
	})
	if err != nil {
		return 0, nil
	}
	var amount float64
	for _, tx := range txs {
		amount += tx.Amount
	}
	return amount, nil
}
