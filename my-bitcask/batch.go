package bitcask

import (
	"sync"

	"github.com/santhozkumar/bitcask/internal"
	"github.com/santhozkumar/bitcask/internal/codec"
	"github.com/santhozkumar/bitcask/internal/config"
)

type batch struct {
	db      DB
	mu      sync.RWMutex
	entries []internal.Entry
	opts    *BatchOptions
}

func DefualtBatchOptions(config *config.Config) *BatchOptions {
	return &BatchOptions{maxKeySize: config.MaxKeySize, maxValueSize: config.MaxValueSize}
}

type BatchOption func(*batch)

type BatchOptions struct {
	maxKeySize   uint32
	maxValueSize uint64
}

func (b *batch) Entries() []internal.Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.entries
}

func (b *batch) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = nil
}

func (b *batch) Delete(key Key) (internal.Entry, error) {
	entry := internal.NewEntry(key, Value(nil))

	b.mu.Lock()
	b.entries = append(b.entries, entry)
	b.mu.Unlock()
	return entry, nil
}

func (b *batch) Put(key Key, value Value) (internal.Entry, error) {
	if len(key) == 0 {
		return internal.Entry{}, ErrEmptyKey
	}

	if b.opts.maxValueSize > 0 && uint64(len(value)) > b.opts.maxValueSize {
		return internal.Entry{}, ErrValueTooLarge
	}
	if b.opts.maxKeySize > 0 && uint32(len(key)) > b.opts.maxKeySize {
		return internal.Entry{}, ErrKeyTooLarge
	}
	entry := internal.NewEntry(key, value)
	b.mu.Lock()
	b.entries = append(b.entries, entry)
	b.mu.Unlock()

	return entry, nil
}

func (b *bitcask) Batch(opts ...BatchOption) Batch {
	bat := &batch{
		db:   b,
		opts: DefualtBatchOptions(b.config)}

	for _, opt := range opts {
		opt(bat)
	}
	return bat
}

func (b *bitcask) WriteBatch(batch Batch) error {

    b.mu.Lock()
    defer b.mu.Unlock()

    if b.current.ReadOnly() {
        return ErrDatabaseReadOnly
    }

    b.metadata.IndexUpToDate = false

	for _, entry := range batch.Entries() {
		if err := b.maybeRotate(); err != nil {
			return nil
		}
		offset, n, err := b.current.Write(entry)

		if err != nil {
			return err
		}

		if b.config.SyncWrites {
			if err := b.current.Sync(); err != nil {
				return err
			}
		}

        b.metadata.IndexUpToDate = false

		if entry.Value != nil {
			if oldItem, found := b.trie.Root().Get(entry.Key); found {
				b.metadata.ReclaimableSpace += oldItem.Size
			}
			item := internal.Item{FileID: b.current.FileID(), Offset: offset, Size: n}
			b.trie, _, _ = b.trie.Insert(entry.Key, item)
		} else {
			if oldItem, found := b.trie.Root().Get(entry.Key); found {
				b.metadata.ReclaimableSpace += oldItem.Size + codec.MetaInfoSize + int64(len(entry.Key))
			}
			b.trie, _, _ = b.trie.Delete(entry.Key)
		}
	}
	return nil
}
