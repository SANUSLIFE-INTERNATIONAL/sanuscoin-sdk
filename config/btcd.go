package config

type BTCD struct {
	Testnet    bool
	ConfigFile string `short:"C" long:"configfile" description:"Path to configuration file"`
	DataDir    string `short:"b" long:"datadir" description:"Directory to store data"`
	DataExtDir string `long:"dataextdir" description:"Directory to store data"`
	LogDir     string `long:"logdir" description:"Directory to logger output."`
	RpcUser    string `short:"u" long:"rpcuser" description:"Username for RPC connections"`
	RpcPass    string `short:"P" long:"rpcpass" default-mask:"-" description:"Password for RPC connections"`
	RpcCert    string `long:"rpccert" description:"File containing the certificate file"`
	RpcKey     string `long:"rpckey" description:"File containing the certificate key"`
	Proxy      string `long:"proxy" description:"Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	ProxyUser  string `long:"proxyuser" description:"Username for proxy server"`
	ProxyPass  string `long:"proxypass" default-mask:"-" description:"Password for proxy server"`
	CpuProfile string `long:"cpuprofile" description:"Write CPU profile to the specified file"`
	DebugLevel string `short:"d" long:"debuglevel" description:"Logging level for all subsystems {trace, debug, info, warn, error, critical} -- You may also specify <subsystem>=<level>,<subsystem2>=<level>,... to set the logger level for individual subsystems -- Use show to list available subsystems"`
}
