package bitcask

import (
	"hash/crc32"

	iradix "github.com/hashicorp/go-immutable-radix/v2"
	"github.com/santhozkumar/bitcask/internal"
	"github.com/santhozkumar/bitcask/internal/config"
	"github.com/santhozkumar/bitcask/internal/data"
)

type transaction struct {
	db        DB
	trie      *iradix.Txn[internal.Item]
	current   data.DataFile
	previous  data.DataFile
	datafiles map[int]data.DataFile
	batch     Batch
	opts      *transactionOptions
}

type transactionOptions struct {
}

func defaultTransactionOptions(*config.Config) *transactionOptions {
	return &transactionOptions{}
}

type TransactionOption func(*transaction) error

func (txn *transaction) Discard() {
}

func (txn *transaction) Commit() error {
	return txn.db.WriteBatch(txn.batch)
}

func (txn *transaction) Has(key Key) bool {
	return true
}

func (txn *transaction) Put(key Key, value Value) error {
	e, err := txn.batch.Put(key, value)
	if err != nil {
		return err
	}

	offset, n, err := txn.current.Write(e)

	if err != nil {
		return err
	}

	item := internal.Item{FileID: txn.current.FileID(), Offset: offset, Size: n}
	_, _ = txn.trie.Insert(key, item)

	return nil
}

func (txn *transaction) Get(key Key) (Value, error) {
	entry, err := txn.get(key)
	if err != nil {
		return nil, err
	}
	return entry.Value, nil

}

func (txn *transaction) get(key Key) (internal.Entry, error) {
    var df data.DataFile

    item, found := txn.trie.Get(key)
    if !found{
        return internal.Entry{}, ErrKeyNotFound 
    }

    switch item.FileID {
    case txn.current.FileID():
        df = txn.current
    case txn.previous.FileID():
        df = txn.previous
    default:
        df = txn.datafiles[item.FileID]
    }

    e, err := df.ReadAt(item.Offset, item.Size)
    if err != nil {
        return internal.Entry{}, err
    }

    checksum := crc32.ChecksumIEEE(e.Value)

    if checksum != e.Checksum {
        return internal.Entry{}, ErrCheckSumFailed
    }

	return e, nil
}

func (txn *transaction) Delete(key Key) error {

    entry, err := txn.batch.Delete(key)
    if err != nil {
        return err
    }
    _, _, err = txn.current.Write(entry)
    if err != nil {
        return err
    }

    _, _ = txn.trie.Delete(key)

	return nil
}

func (b *bitcask) Transaction(opts ...TransactionOption) Transaction {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var current = data.NewInMemoryDataFile(-1, b.config.MaxKeySize, b.config.MaxValueSize)
	var previous = b.current.ReopenReadonly()

	txn := &transaction{
		db:        b,
		batch:     b.Batch(),
		current:   current,
		previous:  previous,
		trie:      b.trie.Txn(),
		datafiles: b.datafiles,
		opts:      defaultTransactionOptions(b.config),
	}

	for _, opt := range opts {
		opt(txn)
	}

	return txn
}

func (txn *transaction) ForEach(f KeyFunc) (err error) {
// type WalkFn[T any] func(k []byte, v T) bool
    txn.trie.Root().Walk(func(key []byte, v internal.Item) bool {
        if err = f(key); err != nil {
            return true
        }
        return false
    })
    return 
}
