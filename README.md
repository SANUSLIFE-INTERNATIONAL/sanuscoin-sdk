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
