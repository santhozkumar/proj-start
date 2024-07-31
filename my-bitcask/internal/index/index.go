package index

import (
	"os"

	iradix "github.com/hashicorp/go-immutable-radix/v2"
	"github.com/santhozkumar/bitcask/internal"
)

type Indexer[T any] interface {
	Load(string, uint32) (*iradix.Tree[T], error)
	Save(*iradix.Tree[T], string) error
}

type indexer struct{}

func NewIndexer() Indexer[internal.Item] {
	return indexer{}
}

func (i indexer) Load(path string, maxKeySize uint32) (*iradix.Tree[internal.Item], error) {
	t := iradix.New[internal.Item]()

	f, err := os.Open(path)
	if err != nil {
		return t, err
	}
	defer f.Close()

	t, err = readIndex(f, t, maxKeySize)

	if err != nil {
		return t, err
	}
	return t, nil
}

func (i indexer) Save(t *iradix.Tree[internal.Item], path string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return WriteIndex(f, t)
}
