package codec

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/santhozkumar/bitcask/internal"
)

var (
	errCantDecodeOnNilEntry = errors.New("Entry is nil")
	errTruncated            = errors.New("Data Truncated")
	errInvalidKeyValueSize  = errors.New("Invalid key/value size")
)

type Decoder struct {
	r            io.Reader
	maxKeySize   uint32
	maxValueSize uint64
}

func NewDecoder(r io.Reader, maxKeySize uint32, maxValueSize uint64) *Decoder {
	return &Decoder{r: r, maxKeySize: maxKeySize, maxValueSize: maxValueSize}
}

func (d *Decoder) Decode(v *internal.Entry) (int64, error) {
	if v == nil {
		return -1, errCantDecodeOnNilEntry
	}
	// find  the keySize and valueSize
	prefixBuf := make([]byte, keySize+valueSize)
	if _, err := io.ReadFull(d.r, prefixBuf); err != nil {
		return -1, errors.New("can't read prefix")
	}
	fmt.Printf("%#v", prefixBuf)
	fmt.Printf("\n")

	actualKeySize, actualValueSize, _ := getKeyandValueSize(prefixBuf, d.maxKeySize, d.maxValueSize)
	fmt.Printf("%T, %T, %T", actualKeySize, actualValueSize, ChecksumSize)
	fmt.Println(keySize, valueSize)

	buf := make([]byte, uint64(actualKeySize)+actualValueSize+ChecksumSize)

	if _, err := io.ReadFull(d.r, buf); err != nil {
		return -1, errors.New("can't read prefix")
	}

	DecodeEntryWithoutPrefix(buf, actualKeySize, v)
	fmt.Printf("%#v", buf)
	fmt.Printf("\n")

	return int64(10), nil
}

func DecodeEntryWithoutPrefix(buf []byte, valueOffset uint32, v *internal.Entry) {
	v.Key = buf[:valueOffset]
	v.Value = buf[valueOffset : len(buf)-ChecksumSize]
	v.Checksum = binary.BigEndian.Uint32(buf[len(buf)-ChecksumSize:])
}

func getKeyandValueSize(prefix []byte, maxKeySize uint32, maxValueSize uint64) (uint32, uint64, error) {

	actualkeysize := binary.BigEndian.Uint32(prefix[:keySize])
	actualvaluesize := binary.BigEndian.Uint64(prefix[keySize:])

	if (actualkeysize > 0 && actualkeysize > maxKeySize) || (actualvaluesize > 0 && actualvaluesize > maxValueSize) || actualkeysize == 0 {
		return 0, 0, errInvalidKeyValueSize
	}

	return actualkeysize, actualvaluesize, nil
}

func DecodeEntry(buf []byte, v *internal.Entry, maxKeySize uint32, maxValueSize uint64) error {
	valueOffset, _, err := getKeyandValueSize(buf, maxKeySize, maxValueSize)
	if err != nil {
		return errors.New("key value error")
	}
	DecodeEntryWithoutPrefix(buf[keySize+valueSize:], valueOffset, v)
	return nil
}
