package sdk

type IWallet interface {
	Create() IWallet
	Load() IWallet
	NewAddress() IAddress
	AddressList() IAddress
}

type BTCWallet struct {
	xpub string
}

func NewBtcWallet() IWallet {
	return &BTCWallet{}
}

func (w *BTCWallet) Create() IWallet {
	return w
}

func (w *BTCWallet) Load() IWallet {
	return w
}

func (w *BTCWallet) NewAddress() IAddress {
	return &BtcAddress{}
}

func (w *BTCWallet) AddressList() IAddress {
	return &BtcAddress{}
}
