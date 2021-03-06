package sdk

import (
	"fmt"
	"net"
	"time"

	"sanus/sanus-sdk/config"
	"sanus/sanus-sdk/kvdb/storage"
	"sanus/sanus-sdk/misc/log"
	"sanus/sanus-sdk/sanus/daemon"

	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/chain"
	wllt "github.com/btcsuite/btcwallet/wallet"
	_ "github.com/btcsuite/btcwallet/walletdb/bdb"

	"github.com/btcsuite/btcd/chaincfg"
)

const defaultLogFName = "wallet.log"

type BTCWallet struct {
	loader    *wllt.Loader
	rpcClient *chain.RPCClient
	wlt       *wllt.Wallet

	netParams *chaincfg.Params

	cfg *config.Config

	lock chan time.Time

	minAmount btcutil.Amount

	*log.Logger

	db *storage.DB

	rpcLaunched chan struct{}
}

// NewWallet creates a new BTCWallet instance
func NewWallet(cfg *config.Config, db *storage.DB) *BTCWallet {
	var param = &chaincfg.TestNet3Params
	if !cfg.Net.Testnet {
		param = &chaincfg.MainNetParams
	}
	loader := wllt.NewLoader(param, config.AppWalletPath(), false, 250)
	logger := log.NewLogger(cfg)
	logger.SetPrefix("WLT")
	return &BTCWallet{
		loader:      loader,
		rpcLaunched: make(chan struct{}),
		Logger:      logger,
		cfg:         cfg,
		lock:        make(chan time.Time, 1),
		netParams:   param,
		db:          db,
	}
}

func (w *BTCWallet) GetNetParams() *chaincfg.Params {
	return w.netParams
}

// Create method creates a new wallet
func (w *BTCWallet) Create(pubPassphrase, privPassphrase, seed []byte) (err error) {
	w.initLogger()
	w.wlt, err = w.loader.CreateNewWallet(pubPassphrase, privPassphrase, seed, time.Now())
	if err != nil {
		return err
	}
	if err = w.Unlock(privPassphrase); err != nil {
		return err
	}
	if err = w.lunchRPC(); err != nil {
		return err
	}
	w.sync()
	return nil
}

// Open method opens already existing wallet
func (w *BTCWallet) Open(pubPassphrase []byte) (err error) {
	if w.wlt != nil {
		return fmt.Errorf("wallet already opened")
	}
	w.wlt, err = w.loader.OpenExistingWallet(pubPassphrase, false)
	if err != nil {
		return err
	}
	if err = w.lunchRPC(); err != nil {
		return err
	}
	w.sync()
	return nil
}

// Unlock method unlocks already initialized wallet
func (w *BTCWallet) Unlock(privatePassphrase []byte) (err error) {
	if w.wlt == nil {
		return fmt.Errorf("wallet hasn't beel opened")
	}
	if w.wlt.Locked() {
		return w.wlt.Unlock(privatePassphrase, w.lock)
	}
	return nil
}

// Lock method locks already initialized wallet
func (w *BTCWallet) Lock() {
	w.lock <- time.Time{}
}

// Synced method returns true if wallet already synced with blockchain
func (w *BTCWallet) Synced() bool {
	return w.wlt.ChainSynced()
}

// Stop method stops wallet and rpc connection
func (w *BTCWallet) Stop() {
	w.rpcClient.Stop()
	w.wlt.Stop()
}

// initLogger method initialized logger for current service
func (w *BTCWallet) initLogger() {
	w.Logger.SetOutput(defaultLogFName, "WALLET")
}

// sync method sync current already initialized wallet with blockchain
func (w *BTCWallet) sync() {
	w.wlt.SynchronizeRPC(w.rpcClient)
	go func() {
		<-w.rpcLaunched
		w.Scan()
	}()
}

// lunchRPC method launches rpc client to connect to blockchain
func (w *BTCWallet) lunchRPC() (err error) {
	activeNet := daemon.ActiveNetParams
	if w.rpcClient, err = chain.NewRPCClient(
		activeNet.Params, net.JoinHostPort("", activeNet.RPCPort),
		daemon.DefaultRPCUser, daemon.DefaultRPCPassword, nil, true, 0); err != nil {
		return err
	}
	go func() {
		if err := w.rpcClient.Start(); err != nil {
			w.Errorf("error caused when trying to start RPC client | %v", err)
			return
		}
		w.rpcLaunched <- struct{}{}
		w.Infof("rpc client has been launched")
	}()
	return err
}
