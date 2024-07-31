package codec

// import (
// 	"bytes"
// 	"fmt"
// 	"testing"
//
// 	"github.com/santhozkumar/bitcask/internal"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )
//
// func TestDecode(t *testing.T) {
// 	buf := bytes.NewBuffer([]byte{0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7, 0x6d, 0x79, 0x6b, 0x65, 0x79, 0x6d, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0, 0x6, 0x51, 0xbd})
// 	dec := NewDecoder(buf, 16, 32)
//
// 	expected := internal.Entry{
// 		Key:      []byte("mykey"),
// 		Value:    []byte("myvalue"),
// 		Checksum: 414141,
// 	}
// 	actual := internal.Entry{}
//
// 	_, err := dec.Decode(&actual)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	require.NoError(t, err)
// 	assert.EqualValues(t, expected, actual)
// }
//
// func TestDecodeEntry(t *testing.T){
// 	buf := bytes.NewBuffer([]byte{0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7, 0x6d, 0x79, 0x6b, 0x65, 0x79, 0x6d, 0x79, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x0, 0x6, 0x51, 0xbd})
// 	expected := internal.Entry{
// 		Key:      []byte("mykey"),
// 		Value:    []byte("myvalue"),
// 		Checksum: 414141,
// 	}
// 	actual := internal.Entry{}
//     err := DecodeEntry(buf.Bytes(), &actual)
// 	require.NoError(t, err)
// 	assert.EqualValues(t, expected, actual)
// }
