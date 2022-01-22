package rpc

import (
	"net"
	"net/http"
	"net/rpc"

	"sanus/sanus-sdk/config"
	"sanus/sanus-sdk/misc/log"
	"sanus/sanus-sdk/sanus/sdk"
)

const (
	defaultLogFName = "rpc.log"
)

type RPCServer struct {
	listener net.Listener
	*log.Logger
	wallet *sdk.BTCWallet
}

func New(cfg *config.Config, wallet *sdk.BTCWallet) *RPCServer {
	srv := &RPCServer{
		Logger: log.NewLogger(cfg),
		wallet: wallet,
	}
	srv.initLogger()
	rpc.HandleHTTP()
	// Create a TCP listener that will listen on `Port`
	listener, err := net.Listen("tcp", cfg.Net.RPC)
	if err != nil {
		srv.Errorf("error caused when trying to listen rpc server %v", err)
		return nil
	}
	srv.listener = listener
	srv.register()
	return srv
}

func (server *RPCServer) register() {
	rpc.Register(NewNetworkHandler(server.wallet))
	rpc.Register(NewTxHandler(server.wallet))
	rpc.Register(NewWalletHandler(server.wallet))
	rpc.Register(NewScriptHandler(server.wallet))
}

func (server *RPCServer) Serve(stopSig chan struct{}) {
	//server.initLogger()
	go func() {
		<-stopSig
		if err := server.listener.Close(); err != nil {
			server.Errorf("error when trying to stop rpc %v", err)
		}
		server.Info("rpc server has been stopped")

	}()
	server.Info("Starting RPC server")
	// Wait for incoming connections
	http.Serve(server.listener, nil)
}

func (server *RPCServer) Close() error {
	return server.listener.Close()
}

func (server *RPCServer) initLogger() {
	server.SetOutput(defaultLogFName, "RPC")
}
