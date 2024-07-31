package metadata

import (
	"os"

	"github.com/santhozkumar/bitcask/internal"
)

type MetaData struct {
	IndexUpToDate    bool  `json:"index_up_to_date"`
	ReclaimableSpace int64 `json:"reclaimable_space"`
}

func (m *MetaData) Save(path string, mode os.FileMode) error {
	return internal.SaveToJSONFile(m, path, mode)

}

func Load(path string) (*MetaData, error) {
	var m MetaData
	err := internal.LoadFromJSONFile(m, path)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
