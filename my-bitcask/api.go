package bitcask

import "github.com/santhozkumar/bitcask/internal"

type Key []byte

type Value []byte

type KeyFunc func(Key) error

type DB interface {
	Keys
	Batch(...BatchOption) Batch

	WriteBatch(Batch) error

    Merge() error
    Close() error
    Sync() error

    Path() string
    Backup(string) error

    Readonly() bool
    Len() int

    ForEach(KeyFunc) error
}

type Keys interface {
	Get(Key) (Value, error)
	Has(Key) bool
	Put(Key, Value) error
    Delete(Key) error
}

type Transaction interface {
	Keys

	Discard()
	Commit() error

    ForEach(KeyFunc) error
}

type Batch interface {
	Put(Key, Value) (internal.Entry, error)
	Delete(Key) (internal.Entry, error)
	Clear()
	Entries() []internal.Entry
}
