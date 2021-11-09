### Install
```Bash
$ go instal
```

### Sanuscoin (sanus-sdk) usage
```Bash
$ sanus start (starts btcd daemon)
$ sanus init (initialize configuration)
$ sanus version (returns current sanus-sdk version)
```


### .env Usage

##### Set some .env variables:

```Bash
APP_NAME="Sanuscoin regular node"
APP_NICK=Sanuscoin
APP_DEBUG=false
APP_VERBOSE=true
```

Results:

```Bash
&config.config{
  App: &struct{
    Name    string
    Nick    string
    Degug   bool
    Verbose bool
  }{
    Name:    "Sanuscoin regular node",
    Nick:    "Sanuscoin",
    Degug:   false,
    Verbose: true,
  },
}
```

### JSON-RPC Usage
````
 Wallet.Create | {"seed":"<generated-seed-hash>","password":"<wallet-password>"}
 Wallet.Open   | {"password":"<wallet-password>"}
 Wallet.Unlock | {"password":"<wallet-password>"}
 Wallet.Lock   | {"success":<bool>}
 Wallet.Synced | {"success":<bool>}
 
 Wallet.NewAddress | {"account":"<account-name>"}
 Wallet.Balance    | {"address":"<wallet-address>","coin":"<btc/snc>"}
 
 
 Tx.Unspent | {"address":"<wallet-address>"}
 Tx.Send    | {"address":"<wallet-address>","amount":<amount>,"pkScript":"<pkscript-hash>"}

 Script.Issuance | {See cc/issuance.ColoredData}
 Script.Transfer | {See cc/transfer.ColoredData}
 
 Network.Status | {}
 
````
