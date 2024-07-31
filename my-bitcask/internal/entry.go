package internal

import "hash/crc32"

type Entry struct {
    Checksum uint32
    Key   []byte
    Value []byte
}

func NewEntry(key, value []byte) Entry {
	checksum := crc32.ChecksumIEEE(value)

    return Entry{
        Key: key,
        Value: value,
        Checksum: checksum,
    }
}

