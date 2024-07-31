package index

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"

	iradix "github.com/hashicorp/go-immutable-radix/v2"
	"github.com/santhozkumar/bitcask/internal"
)

var (
	errTruncatedKeySize = errors.New("key size data may be truncated")
	errKeySizeTooLarge =  errors.New("key size large")
	errTruncatedKeyData = errors.New("key data may be truncated")
	errTruncatedData = errors.New("data may be truncated")
)

const (
	int32Size  = 4
	int64Size  = 4
	fileIDSize = int32Size
	offsetSize = int64Size
	sizeSize   = int64Size
)

func readKeyBytes(r io.Reader, maxKeySize uint32) ([]byte, error) {

	b := make([]byte, int32Size)
	_, err := io.ReadFull(r, b)
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, errors.Wrap(errTruncatedKeySize, err.Error())
	}

    size := binary.BigEndian.Uint32(b)

    if maxKeySize >0 && size > maxKeySize {
        return nil, errKeySizeTooLarge
    }

	k := make([]byte, size)
	_, err = io.ReadFull(r, k)
	if err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, errors.Wrap(errTruncatedKeyData, err.Error())
	}

    return k, nil
}

func writeItem(w io.Writer, item internal.Item) error {
    if err := binary.Write(w, binary.BigEndian, item.FileID); err != nil {
        return err
    }
    if err := binary.Write(w, binary.BigEndian, item.Offset); err != nil {
        return err
    }
    if err := binary.Write(w, binary.BigEndian, item.Size); err != nil {
        return err
    }
   return nil
}


func readItem(r io.Reader) (internal.Item, error) {
    buf := make([]byte, fileIDSize+offsetSize+sizeSize)
    _, err := io.ReadFull(r, buf)
	if err != nil {
        return internal.Item{}, errors.Wrap(errTruncatedData, err.Error()) 
    }

    return internal.Item{
        FileID: int(binary.BigEndian.Uint32(buf[:fileIDSize])),
        Offset: int64(binary.BigEndian.Uint64(buf[fileIDSize:(fileIDSize+offsetSize)])),
        Size: int64(binary.BigEndian.Uint64(buf[(fileIDSize+offsetSize):])),
    }, nil
}


func readIndex(r io.Reader, t *iradix.Tree[internal.Item], maxKeySize uint32) (*iradix.Tree[internal.Item], error) {
	for {
		key, err := readKeyBytes(r, maxKeySize)
        if err != nil {
            if err == io.EOF {
                break
            }
            return t, err
        }
        item, err := readItem(r)
		if err != nil {
			return t, err
		}
        t, _,_ = t.Insert(key, item)
	}

	return t, nil
}

func writeBytes(b []byte, w io.Writer) error {
    err := binary.Write(w, binary.BigEndian, uint32(len(b)))

    if err != nil {
        return err
    }

    _, err = w.Write(b)
    if err != nil {
        return err
    }
    return nil
}

func WriteIndex(w io.Writer, t *iradix.Tree[internal.Item]) (err error) {
    t.Root().Walk(func (k []byte, v internal.Item) bool{
        err = writeBytes(k, w)
        if err != nil {
            return true
        }
        err = writeItem(w, v)
        return err != nil
    })
    return 
}
