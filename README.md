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

### Usage of HTTP API

1.Generate a new seed hash
```CURL
Url: /<APP_NAME>/v1/wallet/seed
Method: POST
Form:  
  mnemonic <mnemonic phrase>
  public <public password>
```

2.Create a new wallet
```CURL
Url: /<APP_NAME>/v1/wallet/create
Method: POST
Form:  
  public <public password>
  private <private password>
  seed <seed hash>
```

3.Open already existing wallet
```CURL
Url: /<APP_NAME>/v1/wallet/open
Method: POST
Form:  
  public <public password>
```

4.Unlock already existing wallet
```CURL
Url: /<APP_NAME>/v1/wallet/unlock
Method: POST
Form:  
  private <private password>
```

5.Lock already existing wallet
```CURL
Url: /<APP_NAME>/v1/wallet/lock
Method: POST
```


6.Check if network synced
```CURL
Url: /<APP_NAME>/v1/wallet/synced
Method: POST
```


7.Create a new address
```CURL
Url: /<APP_NAME>/v1/address/create
Method: POST
Form:  
  account <account name> (not required)
```


8.Get unspent transactions
```CURL
Url: /<APP_NAME>/v1/tx/unspent
Method: POST
Form:  
  address <account address>
```


9.Get unspent transactions
```CURL
Url: /<APP_NAME>/v1/network/status
Method: POST
```

10.Get test SNC Balance
```CURL
Url: /<APP_NAME>/v1/test/do
Method: POST
Form: 
    address <address account>
```