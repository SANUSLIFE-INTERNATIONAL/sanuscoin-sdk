package badger

import (
	"github.com/dgraph-io/badger/v3"
)

// Open tries open data store.
func (b *driver) Open(path string) (err error) {
	// init database options
	opts := badger.DefaultOptions(path)
	opts.Logger = nil // mute unwanted messages

	// try to open the Badger database
	// it will be created if doesn't exist
	b.DB, err = badger.Open(opts)

	return err
}

// Del deletes entry by key.
func (b *driver) Del(key string) error {
	return b.del([]byte(key))
}

// Get retrieves value by key.
func (b *driver) Get(key string) ([]byte, error) {
	return b.get([]byte(key))
}

// Has retrieves key exists in key store.
func (b *driver) Has(key string) bool {
	return b.has([]byte(key))
}

// Set sets specified key-value pair to data store.
func (b *driver) Set(key string, val []byte) error {
	return b.set([]byte(key), val)
}

// Put puts specified key to data store with empty value.
func (b *driver) Put(key string) error {
	return b.set([]byte(key), nil)
}
