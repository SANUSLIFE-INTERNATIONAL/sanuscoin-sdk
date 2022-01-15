package entity

type Entity interface {
	Key() []byte
	Value() []byte
	From([]byte, []byte) error
}
