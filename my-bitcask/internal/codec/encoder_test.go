package codec

// import (
// 	"bytes"
// 	"testing"
//
// 	"github.com/santhozkumar/bitcask/internal"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )
//
// func TestEncode(t *testing.T) {
// 	var buf bytes.Buffer
// 	// msg := internal.NewEntry([]byte("mykey"), []byte("myvalue"))
//
// 	msg := internal.Entry{
// 		Key:      []byte("mykey"),
// 		Value:    []byte("myvalue"),
// 		Checksum: 414141,
// 	}
// 	enc := NewEncode(&buf)
// 	expected := []byte{0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7, 0x6d, 0x79, 0x6b, 0x65, 0x79, 0x6d, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0, 0x6, 0x51, 0xbd}
//
// 	_, err := enc.Encode(msg)
// 	require.NoError(t, err)
//
//     actual := buf.Bytes()
// 	assert.EqualValues(t, expected, actual)
// }
