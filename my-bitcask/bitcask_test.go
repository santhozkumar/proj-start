package bitcask

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestAll(t *testing.T) {
    var (
        db DB 
        testDir string 
        err error 
    )

    testDir, err = os.MkdirTemp("", "bitcask")
    defer os.RemoveAll(testDir)
    assert.NoError(t, err)

    t.Run("Open", func(t *testing.T){
        db, err = Open(testDir)
        assert.NoError(t, err)
    })

    t.Run("Put", func(t *testing.T){
        err = db.Put(Key("foo"), Value("bar"))
        assert.NoError(t, err)
    })

    t.Run("Get", func(t *testing.T){
        actual, err := db.Get(Key("foo"))
        assert.NoError(t, err)
        assert.Equal(t, Value("bar"), actual)
    })

    t.Run("Len", func(t *testing.T){
        l := db.Len()
        assert.Equal(t, 1, l)
    })

    t.Run("Has", func(t *testing.T){
        assert.True(t, db.Has(Key("foo")))
    })

    t.Run("Delete", func(t *testing.T){
        err = db.Delete(Key("foo"))
        assert.NoError(t, err)
        _, err := db.Get(Key("foo"))
        assert.Error(t, err)
        assert.EqualError(t, err, ErrKeyNotFound.Error())
    })

    t.Run("Sync", func(t *testing.T){
        assert.NoError(t, db.Sync())
    })

    t.Run("ForEach", func(t *testing.T){
        var (
            keys []Key
            values []Value
        )

        err := db.ForEach(func(key Key) error {
            value, err := db.Get(key)
            if err != nil {
                return err 
            }

            keys = append(keys,key)
            values = append(values,value)
            return nil
        })

        assert.NoError(t, err)
        assert.Equal(t, []Key{[]byte("foo")}, keys)
        assert.Equal(t, []Value{[]byte("bar")}, keys)

    })

    t.Run("Backup", func(t *testing.T){
        backupDir, err := os.MkdirTemp("", "backup")
        fmt.Println(backupDir)
        defer os.RemoveAll(backupDir)
        assert.NoError(t, err)
        err = db.Backup(filepath.Join(backupDir, "db-backup"))
        assert.NoError(t, err)

    })

    t.Run("Close", func(t *testing.T){
        assert.NoError(t, db.Close())
    })
}


func TestLoadIndexes(t *testing.T){
    testDir, err := os.MkdirTemp("", "bitcask")
	defer os.RemoveAll(testDir)
    assert.NoError(t, err)

    var db DB

    t.Run("Setup", func(t *testing.T){ 
        db, err = Open(testDir)
        assert.NoError(t, err)


        for i :=0; i < 5; i ++ {
            key := fmt.Sprintf("key %d", i)
            value := fmt.Sprintf("value %d", i)
            err := db.Put([]byte(key), []byte(value))
            assert.NoError(t, err)
        }

        assert.NoError(t, db.Close())

    })

    t.Run("OpenAgain", func(t *testing.T) {
        db, err = Open(testDir)
        assert.NoError(t, err)
        assert.Equal(t, t, db.Len())
    })
}


func TestLocking(t *testing.T) {
    testDir, err := os.MkdirTemp("", "bitcask")
    assert.NoError(t, err)

    db, err := Open(testDir)
    assert.NoError(t, err)
    defer db.Close()

    require.NoError(t, db.Put(Key("foo"), Value("bar")))

    _, err = Open(testDir)
    require.Error(t, err)
    assert.EqualError(t, err, ErrDatabaseLocked.Error())

    rdb, err := Open(testDir, WithAutoReadOnly(true))
    require.NoError(t, err)
    assert.True(t, rdb.Readonly())
    defer rdb.Close()


    err = rdb.Put(Key("foo"), Value("bar"))
    require.Error(t, err)
    assert.EqualError(t, err, ErrDatabaseReadOnly.Error())


    err = rdb.Delete(Key("foo"))
    require.Error(t, err)
    assert.EqualError(t, err, ErrDatabaseReadOnly.Error())

    actual, err := db.Get(Key("foo"))
    expected := Value("bar")
    require.NoError(t, err)
    assert.Equal(t, expected, actual)

}
