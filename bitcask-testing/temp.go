package main

import (
	"fmt"
	"sync"
)

type Temp struct {
    mu sync.RWMutex
}

func temp() {
	fmt.Println("in ")

	defer func() {
		fmt.Println("1")
	}()
	defer func() {
		fmt.Println("2")
	}()

	defer func() {
		fmt.Println("3")
	}()
}

func main() {
	// temp()
    // fns, _ := os.ReadDir("/home/santhosh/opensource/bitcask-testing/test.db/")
    // for _, fn := range fns {
    //     fmt.Println(fn.Name())
    // }

    t := Temp{}
    t.mu.RLock()
    t.mu.RLock()
    t.mu.RUnlock()

    t.mu.Lock()
}
