package index

import (
	iradix "github.com/hashicorp/go-immutable-radix/v2"
	"github.com/santhozkumar/bitcask/internal"
)

func getSamples() (*iradix.Tree[internal.Item], int) {
	at := iradix.New[internal.Item]()
	keys := [][]byte{[]byte("abcd"), []byte("abce"), []byte("abcf"), []byte("abgd")}
	expectedSerializedSize := 0
	for i := range keys {
		at.Insert(keys[i], internal.Item{FileID: i, Offset: int64(i), Size: int64(i)})
		expectedSerializedSize += int32Size + len(keys[i]) + fileIDSize + offsetSize + sizeSize
	}

	return at, expectedSerializedSize
}
