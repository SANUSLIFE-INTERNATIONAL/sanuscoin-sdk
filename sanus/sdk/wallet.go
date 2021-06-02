package sdk

import (
	"fmt"
	"net"
	"path/filepath"
	"time"

	"sanus/sanus-sdk/sanus/daemon"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/chain"
	"github.com/btcsuite/btcwallet/waddrmgr"
	wllt "github.com/btcsuite/btcwallet/wallet"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb"

	"github.com/btcsuite/btcd/chaincfg"
)

var (
	defaultWalletDir = filepath.Join(daemon.DefaultHomeDir, "wallet")
)

type IWallet interface {
	Create(pubPassphrase, privPassphrase, seed []byte) error
	Open(pubPassphrase []byte) error
	NewAddress() (btcutil.Address, error)
	AddressList() IAddress
}

type BTCWallet struct {
	loader    *wllt.Loader
	rpcClient *chain.RPCClient
	wlt       *wllt.Wallet
}

// NewWallet creates a new BTCWallet instance
func NewWallet(params *chaincfg.Params) IWallet {
	loader := wllt.NewLoader(params, defaultWalletDir, false, 250)
	return &BTCWallet{
		loader: loader,
	}
}

func (w *BTCWallet) unlock(privatePass []byte) error {
	var lock = make(chan time.Time, 1)
	defer func() {
		lock <- time.Time{}
	}()
	return w.wlt.Unlock(privatePass, lock)
}

func (w *BTCWallet) sync() error {
	rpcClient, err := chain.NewRPCClient(
		daemon.ActiveNetParams.Params, net.JoinHostPort("", daemon.ActiveNetParams.DefaultPort),
		daemon.DefaultRPCUser, daemon.DefaultRPCPassword, nil, true, 0)
	if err != nil {
		return err
	}
	if err = rpcClient.Start(); err != nil {
		return err
	}
	w.wlt.SynchronizeRPC(rpcClient)
	return nil
}

func (w *BTCWallet) Create(pubPassphrase, privPassphrase, seed []byte) (err error) {
	w.wlt, err = w.loader.CreateNewWallet(pubPassphrase, privPassphrase, seed, time.Now())
	if err != nil {
		return err
	}
	if err = w.unlock(privPassphrase); err != nil {
		return err
	}
	if err = w.sync(); err != nil {
		return err
	}
	fmt.Println("Wallet syncing with daemon", w.wlt.SynchronizingToNetwork())
	return nil
}

func (w *BTCWallet) Open(pubPassphrase []byte) (err error) {
	var lock = make(chan time.Time, 1)
	defer func() {
		lock <- time.Time{}
	}()
	w.wlt, err = w.loader.OpenExistingWallet(pubPassphrase, false)
	if err != nil {
		return err
	}
	return w.unlock(pubPassphrase)
}

func (w *BTCWallet) NewAddress() (btcutil.Address, error) {
	idx, err := w.wlt.NextAccount(waddrmgr.KeyScopeBIP0044, "hello")
	if err != nil {
		return nil, err
	}
	return w.wlt.NewAddress(idx, waddrmgr.KeyScopeBIP0044)
}

func (w *BTCWallet) AddressList() IAddress {
	return &BtcAddress{}
}
