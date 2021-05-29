package wallet_btc

type IAddress interface {
	PrivateKey()
	PublicKey()
}

type BtcAddress struct {
	xpub, xpriv string
}

func (a *BtcAddress) PrivateKey() {

}

func (a *BtcAddress) PublicKey() {

}
