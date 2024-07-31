package internal

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

// Exists returns `true` if the given `path` on the current file system exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirSize returns the space occupied by the given `path` on disk on the current
// file system.
func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// GetDatafiles returns a list of all data files stored in the database path
// given by `path`. All datafiles are identified by the the glob `*.data` and
// the basename is represented by a monotonic increasing integer.
// The returned files are *sorted* in increasing order.
func GetDatafiles(path string) ([]string, error) {
	fns, err := filepath.Glob(fmt.Sprintf("%s/*.data", path))
	if err != nil {
		return nil, err
	}
	slices.Sort(fns)
	return fns, nil
}

// ParseIds will parse a list of datafiles as returned by `GetDatafiles` and
// extract the id part and return a slice of ints.
func ParseIds(fns []string) ([]int, error) {
	var ids []int
	for _, fn := range fns {
		fn = filepath.Base(fn)
		ext := filepath.Ext(fn)
		if ext != ".data" {
			continue
		}
		id, err := strconv.ParseInt(strings.TrimSuffix(fn, ext), 10, 32)
		if err != nil {
			return nil, err
		}
		ids = append(ids, int(id))
	}
	slices.Sort(ids)
	return ids, nil
}

func SaveToJSONFile(v any, path string, mode os.FileMode) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, mode)
}

func LoadFromJSONFile(v any, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func CalculateMD5(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	hash := md5.New()

	_, err = io.Copy(hash, f)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

func Copy(src, dst string, exclude []string) error {

// type WalkFunc func(path string, info fs.FileInfo, err error) error
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
        relPath := strings.Replace(path, src, "", 1)

        if relPath == ""{
            return nil
        }
        // the files in exlcude array won't be copied
        for _, e := range exclude {
            match, err := filepath.Match(e, info.Name())
            if err != nil {
                return err 
            }

            if match {
                return nil 
            }
        }

        if info.IsDir(){
            os.Mkdir(filepath.Join(dst, relPath), info.Mode())
        }

        var data, err1 = os.ReadFile(filepath.Join(src, relPath))

        if err1 != nil {
            return err 
        }

        return os.WriteFile(filepath.Join(dst, relPath), data, info.Mode())
	})

}





