package codec

import (
	"bufio"
	"encoding/binary"
	"io"
	"sync"

	"github.com/pkg/errors"
	"github.com/santhozkumar/bitcask/internal"
)

var bufPool = sync.Pool{
	New: func() any {
		return make([]byte, keySize+valueSize)
	},
}

const (
	keySize      = 4
	valueSize    = 8
	ChecksumSize = 4
)

const MetaInfoSize = keySize + valueSize + ChecksumSize

type Encoder struct {
	w *bufio.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{bufio.NewWriter(w)}
}

func (e *Encoder) Encode(msg internal.Entry) (int64, error) {

	var bufKeyValue = bufPool.Get().([]byte)
	binary.BigEndian.PutUint32(bufKeyValue[:keySize], uint32(len(msg.Key)))
	binary.BigEndian.PutUint64(bufKeyValue[keySize:], uint64(len(msg.Value)))

	// fmt.Printf("Before: %#v", bufKeyValue)
	// fmt.Println("")
	if _, err := e.w.Write(bufKeyValue); err != nil {
		return 0, errors.New("Key did not write")
	}
	// fmt.Printf("Aftet: %#v", bufKeyValue)
	// fmt.Println("")
	if _, err := e.w.Write(msg.Key); err != nil {
		return 0, errors.New("Key did not write")
	}
	if _, err := e.w.Write(msg.Value); err != nil {
		return 0, errors.New("Key did not write")
	}
	var bufChecksum = bufKeyValue[:ChecksumSize]
	// fmt.Printf("Before: %#v", bufKeyValue)
	// fmt.Println("")
	binary.BigEndian.PutUint32(bufChecksum, msg.Checksum)
	if _, err := e.w.Write(bufChecksum); err != nil {
		return 0, errors.New("Key did not write")
	}
	// fmt.Printf("After: %#v", bufKeyValue)
	// fmt.Println("")
	// fmt.Printf("After: %#v", bufChecksum)
	// fmt.Println("")

	if err := e.w.Flush(); err != nil {
		return 0, errors.New("failed flushing data")
	}
	return int64(keySize + valueSize + len(msg.Key) + len(msg.Value) + ChecksumSize), nil
}

// func main() {
//     // var s string
//     // for i := range(b){
//     //     s += fmt.Sprintf("%#v ", b[i])
//     // }
//
//
//     // fmt.Printf("%#x", uint(5))
//     // fmt.Printf("%#x", b)
//     // fmt.Println()
//     // fmt.Println(s)
//
//     // s = ""
//     // for i := range(bufKeyValue){
//     //     s += fmt.Sprintf("%#v ",bufKeyValue[i])
//     // }
//     // fmt.Printf("%x\n", bufKeyValue)
//     // fmt.Println()
//     // fmt.Println(s)
//
//     msg := internal.NewEntry([]byte("hello"), []byte("world"))
//     enc := NewEncode(os.Stdout)
//     _, err := enc.Encode(msg)
//     if err != nil {
//         fmt.Println("got error")
//     }
// }
