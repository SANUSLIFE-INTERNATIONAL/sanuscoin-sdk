package sdk

import (
	"fmt"
	"net"
	"path/filepath"
	"time"

	"sanus/sanus-sdk/config"
	"sanus/sanus-sdk/misc/log"
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

const defaultLogFName = "wallet.log"

type BTCWallet struct {
	loader    *wllt.Loader
	rpcClient *chain.RPCClient
	wlt       *wllt.Wallet

	cfg *config.Config

	lock chan time.Time

	*log.Logger
}

// NewWallet creates a new BTCWallet instance
func NewWallet(cfg *config.Config) *BTCWallet {
	var param = &chaincfg.TestNet3Params
	if !cfg.Net.Testnet {
		param = &chaincfg.MainNetParams
	}
	loader := wllt.NewLoader(param, defaultWalletDir, false, 250)
	return &BTCWallet{
		loader: loader,
		Logger: log.NewLogger(cfg),
		cfg:    cfg,
	}
}

func (w *BTCWallet) initLogger() {
	w.Logger.SetOutput(defaultLogFName, "WALLET")
}

func (w *BTCWallet) unlock(privatePass []byte, lock chan time.Time) error {
	return w.wlt.Unlock(privatePass, lock)
}

func (w *BTCWallet) sync() {
	w.wlt.SynchronizeRPC(w.rpcClient)
}

func (w *BTCWallet) Create(pubPassphrase, privPassphrase, seed []byte) (err error) {
	w.initLogger()
	w.wlt, err = w.loader.CreateNewWallet(pubPassphrase, privPassphrase, seed, time.Now())
	if err != nil {
		return err
	}
	w.lock = make(chan time.Time, 1)
	if err = w.unlock(privPassphrase, w.lock); err != nil {
		return err
	}
	if err = w.lunchRPC(); err != nil {
		return err
	}
	w.sync()
	return nil
}

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
		w.Infof("rpc client has been launched")
	}()
	return err
}

func (w *BTCWallet) Open(pubPassphrase []byte) (err error) {
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

func (w *BTCWallet) Unlock(privatePassphrase []byte) (err error) {
	return w.wlt.Unlock(privatePassphrase, w.lock)
}

func (w *BTCWallet) Lock() {
	w.lock <- time.Time{}
}

func (w *BTCWallet) NewAddress(account string) (btcutil.Address, error) {
	idx, err := w.wlt.NextAccount(waddrmgr.KeyScopeBIP0044, account)
	if err != nil {
		return nil, err
	}
	addr, err := w.wlt.NewAddress(idx, waddrmgr.KeyScopeBIP0044)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	w.Infof("Generated new address %v", addr.EncodeAddress())
	return addr, nil
}

func (w *BTCWallet) UnspentTx(addr btcutil.Address) ([]string, error) {
	list, err := w.wlt.ListUnspent(3, 99999, map[string]struct{}{
		addr.EncodeAddress(): {},
	})
	if err != nil {
		return nil, err
	}
	var txs = make([]string, len(list), len(list))
	for k, tx := range list {
		txs[k] = tx.TxID
	}
	return txs, err
}

func (w *BTCWallet) Stop() {
	w.rpcClient.Stop()
	w.wlt.Stop()
}

func (w *BTCWallet) start() {
	w.Info("Starting service")

	if w.wlt != nil {
		go w.wlt.Start()
		w.Info("Wallet has been started")
	}
	//go w.startRPCClient()
	w.Info("RPCClient has been started")

	time.Sleep(3 * time.Second)
	w.Infof("wallet locked %v", w.wlt.Locked())
	return
}

func (w *BTCWallet) AddressList() IAddress {
	return &BtcAddress{}
}
