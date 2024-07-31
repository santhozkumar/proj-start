package bitcask

import "errors"

var (
	ErrValueTooLarge = errors.New("Error: Value too large")

	ErrKeyTooLarge = errors.New("Error: Key too large")

	ErrEmptyKey = errors.New("Error: Empty key")

	ErrDatabaseLocked = errors.New("error: database locked")

	ErrMerginInProgress = errors.New("error: Already mergin in progress")

	ErrKeyNotFound = errors.New("error: Key not found in DB")

	ErrCheckSumFailed = errors.New("error: Checksum value failed")

	ErrDatabaseReadOnly = errors.New("error: Read only database")
)
