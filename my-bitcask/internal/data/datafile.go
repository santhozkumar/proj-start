package data

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattetti/filebuffer"
	"github.com/pkg/errors"
	"github.com/santhozkumar/bitcask/internal"
	"github.com/santhozkumar/bitcask/internal/codec"
	"golang.org/x/exp/mmap"
)

var (
	errReadOnly  = errors.New("Cannot write to readonly file")
	errReadError = errors.New("Read Error")
)

type DataFile interface {
	FileID() int
	Name() string
	Size() int64
	Close() error
	Sync() error
	Read() (internal.Entry, int64, error)
	ReadAt(index, size int64) (internal.Entry, error)
	Write(internal.Entry) (int64, int64, error)

	ReadOnly() bool
	ReopenReadonly() DataFile
}

func NewOnDiskDataFile(path string, id int, readOnly bool, maxKeySize uint32, maxValueSize uint64, fileMode os.FileMode) (DataFile, error) {
	var (
		w   *os.File
		r   *os.File
		ra  *mmap.ReaderAt
		err error
	)

	fn := filepath.Join(path, fmt.Sprintf("%09d.data", id))

	if !readOnly {
		w, _ = os.OpenFile(fn, os.O_WRONLY|os.O_APPEND|os.O_CREATE, fileMode)
		fmt.Println("write file descriptor", w.Fd())
	}

	r, err = os.Open(fn)
	if err != nil {
		return nil, err
	}

	fmt.Println("read file descriptor", r.Fd())
	stat, err := os.Stat(fn)
	if err != nil {
		return nil, errors.Wrap(err, "stat failed")
	}

	if readOnly {
		w = nil
		ra, err = mmap.Open(fn)
		if err != nil {
			return nil, errors.Wrap(err, "mmap not loading")
		}
	}

	offset := stat.Size()

	dec := codec.NewDecoder(r, maxKeySize, maxValueSize)
	enc := codec.NewEncoder(w)

	return &onDiskDataFile{
		id:           id,
		r:            r,
		ra:           ra,
		w:            w,
		offset:       offset,
		enc:          enc,
		dec:          dec,
		maxKeySize:   maxKeySize,
		maxValueSize: maxValueSize,
	}, nil
}

func NewInMemoryDataFile(id int, maxKeySize uint32, maxValueSize uint64) DataFile {

	buf := filebuffer.New(nil)
	dec := codec.NewDecoder(buf, maxKeySize, maxValueSize)
	enc := codec.NewEncoder(buf)
	return &InMemoryDataFile{
		id:           id,
		buf:          buf,
		enc:          enc,
		dec:          dec,
		maxKeySize:   maxKeySize,
		maxValueSize: maxValueSize,
	}
}
