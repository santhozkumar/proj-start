package data

import (
	"os"
	"sync"

	"golang.org/x/exp/mmap"
	"github.com/pkg/errors"
	"github.com/santhozkumar/bitcask/internal/codec"
    "github.com/santhozkumar/bitcask/internal"
)


type onDiskDataFile struct {
	sync.RWMutex

	id     int
	r      *os.File
	ra     *mmap.ReaderAt
	w      *os.File
	enc    *codec.Encoder
	dec    *codec.Decoder
	offset int64
    maxKeySize uint32
    maxValueSize uint64
}

func (df *onDiskDataFile) FileID() int {
	return df.id
}

func (df *onDiskDataFile) Name() string {
	return df.r.Name()
}

func (df *onDiskDataFile) Size() int64 {
    df.Lock()
    defer df.Unlock()
	return df.offset
}

func (df *onDiskDataFile) Sync() error {
    if df.w == nil{
        return nil
    }
    return df.w.Sync()
}

func (df *onDiskDataFile) Close() error {
    defer func(){
        if df.ra != nil {
            df.ra.Close()
        }
        df.r.Close()
    }()

    if df.w == nil {
        return nil
    }

    err := df.w.Sync()
    if err != nil {
        return errors.Wrap(err, "while syncing in close")
    }

    return df.w.Close()
}

func (df *onDiskDataFile) Read() (e internal.Entry, n int64, err error) {
    df.Lock()
    defer df.Unlock()

    n, err = df.dec.Decode(&e)
    if err != nil {
        return 
    }
    return
}


func (df *onDiskDataFile) ReadAt(index, size int64) (e internal.Entry, err error) {
    df.RLock()
    defer df.RUnlock()

    buf := make([]byte, size)
    if df.ra != nil {
        df.ra.ReadAt(buf, index)
    } else {
        df.r.ReadAt(buf, index)
    }

    err = codec.DecodeEntry(buf, &e, df.maxKeySize, df.maxValueSize)
    if err != nil {
        return 
    }
    return
}


func (df *onDiskDataFile) Write(e internal.Entry)(int64, int64, error){
    if df.w == nil {
        return -1, 0, errReadOnly
    }

    df.Lock()
    defer df.Unlock()
    offset := df.offset

    n, err := df.enc.Encode(e)
    if err != nil {
        return -1, 0, err
    }

    df.offset += n
    return offset, n, nil
}

func (df *onDiskDataFile) ReadOnly() bool {
	return df.w == nil
}

func (df *onDiskDataFile) ReopenReadonly() DataFile {
    df.RLock()
    defer df.RUnlock()

	return &onDiskDataFile{
		id:           df.id,
		r:            df.r,
		ra:           df.ra,
		w:            nil,
		offset:       df.offset,
		enc:          df.enc,
		dec:          df.dec,
		maxKeySize:   df.maxKeySize,
		maxValueSize: df.maxValueSize,
}
}
