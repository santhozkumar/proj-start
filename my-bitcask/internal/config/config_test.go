package config

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)



func TestLoad(t *testing.T){
    got, _  := Load("./test_config.json")
    want := &Config{
        MaxDataFileSize: 1048576,
        MaxKeySize: 64,
        MaxValueSize: 65536,
        Sync: false,
        AutoRecovery: false,
        DirMode: 448,
        FileMode: 384,
    }
    if !reflect.DeepEqual(*got, *want){
    // if *got != *want{
        spew.Dump("got:", *got)
        spew.Dump("want:", *want)
        t.Errorf("got %v want %v", got, want)
    }
}
