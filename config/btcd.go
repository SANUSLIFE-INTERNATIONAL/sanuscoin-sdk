package config

type BTCD struct {
	Testnet    bool
	ConfigFile string `short:"C" long:"configfile" description:"Path to configuration file"`
	DataDir    string `short:"b" long:"datadir" description:"Directory to store data"`
}
