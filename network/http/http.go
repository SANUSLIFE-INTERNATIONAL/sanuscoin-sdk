package http

import (
	"fmt"
	coreHttp "net/http"

	"sanus/sanus-sdk/misc/log"

	"sanus/sanus-sdk/config"
	"sanus/sanus-sdk/sanus/sdk"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	defaultLogFName = "http.log"
)

type HTTPServer struct {
	cfg *config.Config

	*log.Logger

	wallet *sdk.BTCWallet
}

func NewHTTP(cfg *config.Config, wallet *sdk.BTCWallet) *HTTPServer {
	logger := log.NewLogger(cfg)
	return &HTTPServer{
		cfg: cfg,

		Logger: logger,

		wallet: wallet,
	}
}

func (server *HTTPServer) router() *mux.Router {
	routers := mux.NewRouter().StrictSlash(true)

	basePath := fmt.Sprintf("/%s/%s", server.cfg.App.Name, ProtocolVersion)

	r := routers.PathPrefix(basePath).Subrouter()

	// Routes for receiving messages from Camunda
	wallet := r.PathPrefix("/wallet").Subrouter()
	wallet.Path("/seed").Methods("POST").Handler(appHandler(server.Seed))
	wallet.Path("/create").Methods("POST").Handler(appHandler(server.CreateWallet))
	wallet.Path("/open").Methods("POST").Handler(appHandler(server.OpenWallet))
	wallet.Path("/unlock").Methods("POST").Handler(appHandler(server.Unlock))
	wallet.Path("/lock").Methods("POST").Handler(appHandler(server.Lock))
	wallet.Path("/synced").Methods("POST").Handler(appHandler(server.Synced))

	address := r.PathPrefix("/address").Subrouter()
	address.Path("/create").Methods("POST").Handler(appHandler(server.NewAddress))
	address.Path("/balance").Methods("POST").Handler(appHandler(server.Balance))
	address.Path("/import").Methods("POST").Handler(appHandler(server.ImportAddress))
	address.Path("/list").Methods("POST").Handler(appHandler(server.ListAddresses))

	tx := r.PathPrefix("/tx").Subrouter()
	tx.Path("/unspent").Methods("POST").Handler(appHandler(server.UnspentTX))
	tx.Path("/send").Methods("POST").Handler(appHandler(server.SendTx))
	tx.Path("/pk-script").Methods("POST").Handler(appHandler(server.Script))

	network := r.PathPrefix("/network").Subrouter()
	network.Path("/status").Methods("POST").Handler(appHandler(server.NetworkStatus))

	coreHttp.Handle("/", handlers.CombinedLoggingHandler(server.Out(), routers))
	return routers

}

func (server *HTTPServer) initLogger() {
	server.SetOutput(defaultLogFName, "HTTP")
}

func (server *HTTPServer) Serve(stopSignal chan struct{}) {
	server.initLogger()
	router := server.router()

	srv := coreHttp.Server{
		Addr:    server.cfg.Net.Http,
		Handler: router,
	}

	go func() {
		<-stopSignal
		if err := srv.Shutdown(nil); err != nil {
			server.Infof("Can't stop server | %v", err)
		}
		server.Info("Server has been stopped")
	}()

	server.Infof("Starting server %v", server.cfg.Net.Http)

	if err := srv.ListenAndServe(); err != nil && err != coreHttp.ErrServerClosed {
		server.Errorf("Can't start server | %v", err)
	}
}
