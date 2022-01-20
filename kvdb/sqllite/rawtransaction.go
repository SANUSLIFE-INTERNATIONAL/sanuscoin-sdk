package sqllite

import (
	_ "github.com/mattn/go-sqlite3"
)

type RawTransactionEntity struct {
	TxId   string
	Assets []string
}
