package data

import (
	"sync"
    "fmt"

	"github.com/mattetti/filebuffer"
	"github.com/pkg/errors"
	"github.com/santhozkumar/bitcask/internal"
	"github.com/santhozkumar/bitcask/internal/codec"
)

type InMemoryDataFile struct {
    sync.RWMutex
    id int
    buf *filebuffer.Buffer
    offset int64
    enc *codec.Encoder
    dec *codec.Decoder
    maxKeySize uint32
    maxValueSize uint64
}

func (df *InMemoryDataFile) FileID() int {
    return df.id
}

func (df *InMemoryDataFile) Name() string {
	return fmt.Sprintf("in-memory-%d", df.id)
}


func (df *InMemoryDataFile) Size() int64 {
    df.Lock()
    df.Unlock()
    return df.offset
}

func (df *InMemoryDataFile) Close() error {
    return nil
}

func (df *InMemoryDataFile) Sync() error {
    return nil
}


func (df *InMemoryDataFile) Read() (e internal.Entry, n int64, err error) {

    df.RLock()
    defer df.RUnlock()

    n, err = df.dec.Decode(&e)
    if err != nil {
        return 
    }
    return 
}


func (df *InMemoryDataFile) ReadAt(index, size int64) (e internal.Entry, err error) {
    b := make([]byte, size)
    df.RLock()
    defer df.RUnlock()

    n, err := df.buf.ReadAt(b, index)
    if err != nil {
        return 
    }

    if int64(n) != size {
        err = errReadError
        return
    }

    err = codec.DecodeEntry(b, &e, df.maxKeySize, df.maxValueSize)
    if err != nil {
        return 
    }
    return
}

func (df *InMemoryDataFile) Write(e internal.Entry) (int64, int64, error){
    df.Lock()
    defer df.Unlock()

    offset := df.offset

    n, err := df.enc.Encode(e)
    if err != nil {
        return -1, 0, errors.Wrap(err, "unable to write for in memory file")
    }
    df.offset += n

    return offset, df.offset, nil
}


func (df *InMemoryDataFile) ReadOnly() bool {
    return true
}

func (df *InMemoryDataFile) ReopenReadonly() DataFile {
    return df
}
