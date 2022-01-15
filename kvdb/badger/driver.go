package badger

import (
	"crypto/subtle"

	sncDriver "sanus/sanus-sdk/kvdb/driver"

	bdg "github.com/dgraph-io/badger/v3"
)

type (
	// driver implements Driver for data store.
	driver struct {
		*bdg.DB
	}
)

var (
	// Make sure driver implements Driver interface.
	_ sncDriver.Driver = (*driver)(nil)
)

// DataStore returns driver driver interface
// with Badger DB instance implementation.
func DataStore() sncDriver.Driver {
	return &driver{DB: new(bdg.DB)}
}

// del deletes entry by key.
func (b *driver) del(key []byte) error {
	return b.Update(
		func(txn *bdg.Txn) error {
			return txn.Delete(key)
		})
}

// get retrieves value by key.
func (b *driver) get(key []byte) (val []byte, err error) {
	err = b.View(
		func(txn *bdg.Txn) error {
			item, err := txn.Get(key)
			if err == bdg.ErrKeyNotFound {
				return nil
			}
			val, err = item.ValueCopy(val)

			return err
		})

	return val, err
}

// has retrieves key exists in key store.
func (b *driver) has(key []byte) bool {
	has := false
	err := b.View(
		func(txn *bdg.Txn) error {
			opts := bdg.DefaultIteratorOptions
			opts.PrefetchValues = false // key-only iteration
			iter := txn.NewIterator(opts)
			defer iter.Close()
			for iter.Seek(key); iter.ValidForPrefix(key); iter.Next() {
				if has = subtle.ConstantTimeCompare(key, iter.Item().Key()) == 1; has {
					return nil
				}
			}
			return nil
		})

	return err == nil && has
}

//set method sets specified key-value pair to data store.
func (b *driver) set(key, val []byte) error {
	return b.Update(
		func(txn *bdg.Txn) error {
			return txn.Set(key, val)
		})
}
