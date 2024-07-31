package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

type Config struct {
	MaxDataFileSize int `json:"max_datafile_size"`
	MaxKeySize      uint32 `json:"max_key_size"`
	MaxValueSize    uint64 `json:"max_value_size"`
	SyncWrites      bool `json:"sync"`
	AutoRecovery    bool `json:"autorecovery"`
	AutoReadOnly    bool `json:"auto_read_only"`
	DirMode         os.FileMode `json:"dir_mode"`
	FileMode        os.FileMode `json:"file_mode"`
}


func Load (path string) (*Config, error) {
    var cfg Config
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, errors.Wrap(err, "got error while loading config")
    }
    
    if err = json.Unmarshal(data, &cfg); err != nil {
        return nil, errors.Wrap(err, "got error while loading Unmarshalling")
    }

    return &cfg, nil
}

func (c *Config) Save(path string) error {
    data, err := json.Marshal(c)
    if err != nil {
        return errors.Wrap(err, "error while saving the file")
    }
    if err := os.WriteFile(path, data, c.FileMode); err != nil {
        return errors.Wrap(err, "cannot write to file")
    }
    return nil
}
