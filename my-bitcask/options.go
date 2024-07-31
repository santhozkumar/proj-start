package bitcask

import (
	"os"

	"github.com/santhozkumar/bitcask/internal/config"
)

const (
	DefaultMaxDataFileSize = 1 << 20
	DefaultKeySize         = uint32(64)
	DefaultValueSize       = uint64(1 << 16)
	DefaultSyc             = false
	DefaultAutoRecovery    = false
	DefaultAutoReadOnly    = false
	DefaultDirMode         = os.FileMode(0700)
	DefaultFileMode        = os.FileMode(0600)
)

type Option func(*config.Config) error

func WithMaxDataFileSize(size int) Option {
	return func(c *config.Config) error {
		c.MaxDataFileSize = size
		return nil
	}
}

func WithMaxKeySize(size uint32) Option {
	return func(c *config.Config) error {
		c.MaxKeySize = size
		return nil
	}
}

func WithValueSize(size uint64) Option {
	return func(c *config.Config) error {
		c.MaxValueSize = size
		return nil
	}
}

func WithSyncWrites(enabled bool) Option {
	return func(c *config.Config) error {
		c.SyncWrites = enabled
		return nil
	}
}

func WithAutoRecovery(enabled bool) Option {
	return func(c *config.Config) error {
		c.AutoRecovery = enabled
		return nil
	}
}

func WithAutoReadOnly(enabled bool) Option {
	return func(c *config.Config) error {
		c.AutoReadOnly = enabled
		return nil
	}
}

func WithDirMode(mode os.FileMode) Option {
	return func(c *config.Config) error {
		c.DirMode = mode
		return nil
	}
}

func WithFileMode(mode os.FileMode) Option {
	return func(c *config.Config) error {
		c.FileMode = mode
		return nil
	}
}

func newDefaultConfig() *config.Config {
	return &config.Config{
		MaxDataFileSize: DefaultMaxDataFileSize,
		MaxKeySize:      DefaultKeySize,
		MaxValueSize:    DefaultValueSize,
		// Sync:            DefaultSyc,
		AutoRecovery: DefaultAutoRecovery,
		AutoReadOnly: DefaultAutoReadOnly,
		DirMode:      DefaultDirMode,
		FileMode:     DefaultFileMode,
	}
}
