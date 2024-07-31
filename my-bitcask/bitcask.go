package bitcask

import (
	"hash/crc32"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"

	"github.com/gofrs/flock"
	iradix "github.com/hashicorp/go-immutable-radix/v2"
	"github.com/santhozkumar/bitcask/internal"
	"github.com/santhozkumar/bitcask/internal/config"
	"github.com/santhozkumar/bitcask/internal/data"
	"github.com/santhozkumar/bitcask/internal/index"
	"github.com/santhozkumar/bitcask/internal/metadata"
	"golang.org/x/exp/slices"
)

const lockfile = "lock"

type bitcask struct {
	mu        sync.RWMutex
	current   data.DataFile
	flock     *flock.Flock
	path      string
	options   []Option
	config    *config.Config
	datafiles map[int]data.DataFile
	trie      *iradix.Tree[internal.Item]
	metadata  *metadata.MetaData
	indexer   index.Indexer[internal.Item]
	isMerging bool
}

func (b *bitcask) Get(key Key) (Value, error) {
	return b.Transaction().Get(key)
}

func (b *bitcask) Has(key Key) bool {
	return b.Transaction().Has(key)
}

func (b *bitcask) Merge() error {
    b.mu.Lock()

    if b.current.ReadOnly(){
        b.mu.Unlock()
        return ErrDatabaseReadOnly
    }
    if b.isMerging {
        b.mu.Unlock()
        return ErrMerginInProgress
    }
    b.isMerging = true
    b.mu.Unlock()

    defer func(){
        b.isMerging = false
    }()

    b.mu.Lock()
    err := b.closeCurrentDatafile()
    if err != nil {
        b.mu.Unlock()
        return err
    }

    filesToMerge := make([]int, 0, len(b.datafiles))
    for k := range b.datafiles {
        filesToMerge = append(filesToMerge, k)
    }
    err = b.openNewWritableFile()
    if err != nil {
        b.mu.Unlock()
        return err
    }
    b.mu.Unlock()
    slices.Sort(filesToMerge)


    temp, err := os.MkdirTemp(b.path, "merge")
    if err != nil {
        return err
    }

    defer os.RemoveAll(temp)
    mdb, err := Open(temp, WithAutoReadOnly(true))


    b.trie.Root().Walk( func(key []byte, item internal.Item) bool {

        if item.FileID > filesToMerge[len(filesToMerge) -1]{
            return false
        }

        e, err := b.read(key)
        if err != nil {
            return true
        }

        if err := mdb.Put(key, e.Value); err != nil {
            return true
        }
        return false
    })

    if err := mdb.Close(); err != nil {
        return err
    }

    b.mu.Lock()
    defer b.mu.Unlock()

	if err = b.close(); err != nil {
		return err
	}
    // delete the datafiles

    files, err := os.ReadDir(b.path)
    if err != nil {
        return err
    }

    for _, file := range files {

        if file.IsDir() || file.Name() == lockfile {
            continue
        }

        ids, err := internal.ParseIds([]string{file.Name()})
        if err != nil {
            return err
        }

        if len(ids) > 0 && ids[0] > filesToMerge[len(filesToMerge)-1] {
            continue
        }

        err = os.RemoveAll(path.Join(b.path, file.Name()))
        if err != nil {
            return err
        }

    }
    // Rename the merged datafiles to actual db path
    files, err = os.ReadDir(mdb.Path())
    if err != nil {
        return err
    }

    for _, file := range files {

        if file.Name() == lockfile {
            continue
        }
        err = os.Rename(path.Join([]string{b.path, file.Name()}...),
                path.Join([]string{mdb.Path(), file.Name()}...))

        if err != nil {
            return err
        }
    }

    b.metadata.ReclaimableSpace = 0
    return b.reopen(false)
}

func (b *bitcask) Path() string {
    return b.path
}

func (b *bitcask) Delete(key Key) error {
    txn := b.Transaction()
    defer txn.Discard()

    err := txn.Delete(key)
    if err != nil {
        return err
    }
    return txn.Commit()
}

func (b *bitcask) Put(key Key, value Value) error {

    b.mu.RLock()
    if b.current.ReadOnly() {
        b.mu.RUnlock()
        return ErrDatabaseReadOnly
    }

    b.mu.RUnlock()

    txn := b.Transaction()
    defer txn.Discard()

    if err:= txn.Put(key, value); err != nil {
        return err
    }

    return txn.Commit()
}


func (b *bitcask) read (key []byte) (internal.Entry, error) {
    var df data.DataFile

    b.mu.RLock()
    item, found := b.trie.Root().Get(key)
    b.mu.RUnlock()
    if !found {
        return internal.Entry{}, ErrKeyNotFound
    }
    if b.current.FileID() == item.FileID{
        df = b.current
    }else {
        df = b.datafiles[item.FileID]
    }
    e, err := df.ReadAt(item.Offset, item.Size)
    if err != nil {
        return internal.Entry{}, err
    }

    checksum := crc32.ChecksumIEEE(e.Value)
    if checksum != e.Checksum {
        return internal.Entry{}, ErrKeyNotFound
    }

    return e, nil
}



func (b *bitcask) openNewWritableFile() error {
    id := b.current.FileID() + 1
    current, err := data.NewOnDiskDataFile(b.path,
        id, false, b.config.MaxKeySize,
        b.config.MaxValueSize,
        b.config.FileMode)
    if err != nil {
        return err
    }

    b.current = current
    return nil

}


func (b *bitcask) closeCurrentDatafile() error {
    if err := b.current.Close(); err != nil {
        return err
    }

    id := b.current.FileID()
    df, err := data.NewOnDiskDataFile(b.path, id, true, 
    b.config.MaxKeySize, b.config.MaxValueSize, b.config.FileMode)

    if err != nil {
        return err
    }

    b.datafiles[id] = df
    return nil
}

func Open(path string, options ...Option) (DB, error) {
	var (
		cfg  *config.Config
		err  error
		meta *metadata.MetaData
	)

	configPath := filepath.Join(path, "config.json")
	if internal.Exists(configPath) {
		cfg, err = config.Load(configPath)
		if err != nil {
			return nil, err
		}
	} else {
		cfg = newDefaultConfig()
	}

	for _, opt := range options {
		if err := opt(cfg); err != nil {
			return nil, err
		}
	}

	if err = os.MkdirAll(path, cfg.DirMode); err != nil {
		return nil, err
	}

	meta, err = loadMetaData(path)
	if err != nil {
		return nil, err
	}

	db := &bitcask{
		flock:    flock.New(filepath.Join(path, lockfile)),
		metadata: meta,
		config:   cfg,
		path:     path,
		options:  options,
		trie:     iradix.New[internal.Item](),
		indexer:  index.NewIndexer(),
	}

	ok, err := db.flock.TryLock()
	if err != nil {
		return nil, err
	}

	if !ok {
		if !db.config.AutoReadOnly {
			return nil, ErrDatabaseLocked
		}

		if err := db.reopen(true); err != nil {
			return nil, err
		}
		return db, nil
	}

	if err := cfg.Save(configPath); err != nil {
		return nil, err
	}

	if err := db.reopen(false); err != nil {
		return nil, err
	}

	return db, nil
}

func (b *bitcask) reopen(readonly bool) error {
	datafiles, lastId, err := loadDataFiles(b.path,
		b.config.MaxKeySize,
		b.config.MaxValueSize,
		b.config.DirMode)

	t, err := loadIndexes(b, datafiles)

	if err != nil {
		return err
	}
	current, err := data.NewOnDiskDataFile(b.path, lastId, readonly,
		b.config.MaxKeySize, b.config.MaxValueSize, b.config.FileMode)

	if err != nil {
		return err
	}

	b.current = current
	b.trie = t
	b.datafiles = datafiles

	return nil
}

func loadIndexes(b *bitcask, datafiles map[int]data.DataFile) (*iradix.Tree[internal.Item], error) {

    t, err := b.indexer.Load(filepath.Join(b.path, "index"), b.config.MaxKeySize)
    if err != nil {
        return loadIndexFromDatafiles(datafiles)
    }
    if !b.metadata.IndexUpToDate {
        return loadIndexFromDatafiles(datafiles)
    }

    return t, nil

}

func loadIndexFromDatafile(df data.DataFile, t *iradix.Tree[internal.Item]) (*iradix.Tree[internal.Item], error){
    var offset int64
    for {
        e, n, err := df.Read()
        if err != nil {
            if err == io.EOF {
                break 
            }
            return t, err
        }
        if len(e.Value) == 0 {
            t, _, _ = t.Delete(e.Key)
            offset += n
            continue
        }
        item := internal.Item{FileID: df.FileID(), Offset: offset, Size: n}
        t, _, _ = t.Insert(e.Key, item)
        offset += n
    }
    return t, nil
}


func loadIndexFromDatafiles(datafiles map[int]data.DataFile)(*iradix.Tree[internal.Item], error) {
    t := iradix.New[internal.Item]()

    sortedDatafiles := getSortedDataFiles(datafiles)

    for _, df:= range sortedDatafiles {
        t, err := loadIndexFromDatafile(df,t)
        if err != nil {
            return t, nil
        }
    }
    return t, nil
}

func getSortedDataFiles(datafiles map[int]data.DataFile) []data.DataFile {
    out := make([]data.DataFile, len(datafiles))
    idx := 0
    for _, df  := range datafiles {
        out[idx] = df
        idx++
    }

    sort.Slice(out, func(i, j int) bool {
        return out[i].FileID() < out[j].FileID()
    })

    return out
}

func loadDataFiles(path string, maxKeySize uint32, maxValueSize uint64, fileModeBeforeUmask os.FileMode) (datafiles map[int]data.DataFile, lastId int, err error) {
	fns, err := internal.GetDatafiles(path)
	if err != nil {
		return nil, 0, err
	}

	ids, err := internal.ParseIds(fns)
	if err != nil {
		return nil, 0, err
	}

	datafiles = make(map[int]data.DataFile, len(ids))
	for _, id := range ids {
		datafiles[id], err = data.NewOnDiskDataFile(path,
			id, true, maxKeySize, maxValueSize, fileModeBeforeUmask)
		if err != nil {
			return
		}
	}
	if len(ids) > 0 {
		lastId = ids[len(ids)-1]
	}
	return
}

func (b *bitcask) maybeRotate() error {
	size := b.current.Size()
	if size < int64(b.config.MaxDataFileSize) {
		return nil
	}
	err := b.current.Close()
	if err != nil {
		return err
	}

	id := b.current.FileID()

	df, err := data.NewOnDiskDataFile(b.path, id, true,
		b.config.MaxKeySize,
		b.config.MaxValueSize,
		b.config.FileMode)

	b.datafiles[id] = df

	id = b.current.FileID() + 1
	current, err := data.NewOnDiskDataFile(b.path, id, false,
		b.config.MaxKeySize,
		b.config.MaxValueSize,
		b.config.FileMode)

	if err != nil {
		return err
	}
	b.current = current

	return nil
}

func loadMetaData(path string) (*metadata.MetaData, error) {
	if !internal.Exists(filepath.Join(path, "meta.json")) {
		meta := new(metadata.MetaData)
		return meta, nil
	}

	return metadata.Load(filepath.Join(path, "meta.json"))
}

func (b *bitcask) Close() error {
    b.mu.Lock()

    defer func() {
        b.flock.Unlock()
        b.mu.Unlock()
    }()

    return b.close()
}

func (b *bitcask) saveIndexes() error {
    return b.indexer.Save(b.trie, path.Join(b.path, "index"))
}

func (b *bitcask) saveMetadata() error {
    return b.metadata.Save(path.Join(b.path, "meta.json"), b.config.FileMode)
}
func (b *bitcask) Readonly() bool {
    b.mu.RLock()
    defer b.mu.RUnlock()

    return b.current.ReadOnly()
}

func (b *bitcask) close() error {
    if !b.current.ReadOnly() {
        if err := b.saveIndexes(); err != nil {
            return err
        }

        b.metadata.IndexUpToDate = true
        if err := b.saveMetadata(); err != nil {
            return err
        }
    }

    for _, datafile := range b.datafiles {
        if err := datafile.Close(); err != nil {
            return err
        }
    }

    return b.current.Close()
}

func (b *bitcask) Len() int {
    b.mu.RLock()
    defer b.mu.RUnlock()

    return b.trie.Len()
}


func (b *bitcask) Sync() error {
    if b.Readonly() {
        return nil
    }
    if err := b.saveMetadata(); err != nil {
        return err
    }
    return b.current.Sync()
}


func (b *bitcask) Backup(path string) error {
    if !internal.Exists(path){
        if err := os.Mkdir(path, b.config.DirMode); err != nil {
            return err
        }
    }

    return internal.Copy(b.path, path, []string{lockfile})
}

func (b *bitcask) ForEach(f KeyFunc) error {
    return b.Transaction().ForEach(f)
}





