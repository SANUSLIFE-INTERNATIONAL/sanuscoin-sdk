package driver

type Driver interface {
	// Close tries close data store.
	Close() error

	// Open tries open data store.
	Open(path string) error

	// Del deletes entry by key.
	Del(key string) error

	// Get retrieves value by key.
	Get(key string) ([]byte, error)

	// Has retrieves key exists in key store.
	Has(key string) bool

	// Set sets specified key-value pair to data store.
	Set(key string, val []byte) error

	// Put puts specified key to data store with empty value.
	Put(key string) error
}
