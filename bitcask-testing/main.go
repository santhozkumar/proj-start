package main

import (
	"fmt"
	"log"

	"go.mills.io/bitcask/v2"
)

func readLoop(m map[int]int) {
	for {
		for k, v := range m {
			fmt.Println(k, "-", v)
		}
	}
}

func old_main() {
	db, err := bitcask.Open("test.db", bitcask.WithAutoRecovery(true), bitcask.WithAutoReadonly(true))
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()

	if err := db.Put(bitcask.Key("hello"), bitcask.Value("world")); err != nil {
		log.Fatalf("got error while inserting: %v", err)
	}
	if err := db.Put(bitcask.Key("bar"), bitcask.Value("baz")); err != nil {
		log.Fatal(err)
	}
	if err := db.Put(bitcask.Key("hello"), bitcask.Value("world")); err != nil {
		log.Fatal(err)
	}

	l := db.List(bitcask.Key("fruits"))
	if err := l.Append(bitcask.Value("Apples")); err != nil {
		log.Fatal(err)
	}
	if err := l.Append(bitcask.Value("Bananas")); err != nil {
		log.Fatal(err)
	}
	if err := l.Append(bitcask.Value("Oranges")); err != nil {
		log.Fatal(err)
	}

	h := db.Hash(bitcask.Key("acronyms"))
	if err := h.Set(bitcask.Key("CPU"), bitcask.Value("Central Processing Unit")); err != nil {
		log.Fatal(err)
	}
	if err := h.Set(bitcask.Key("RAM"), bitcask.Value("Random Access Memory")); err != nil {
		log.Fatal(err)
	}
	if err := h.Set(bitcask.Key("HDD"), bitcask.Value("Hard Disk Drive")); err != nil {
		log.Fatal(err)
	}

	s := db.SortedSet(bitcask.Key("scores"))
	if _, err := s.Add(
		bitcask.Int64ToScore(100), bitcask.Key("Bob"),
		bitcask.Int64ToScore(200), bitcask.Key("Dan"),
		bitcask.Int64ToScore(300), bitcask.Key("Joe"),
	); err != nil {
		log.Fatal(err)
	}
	m := map[int]int{}

	go readLoop(m)
	// messages := make(chan string)
	block := make(chan struct{})
	<-block
	// <-messages

}
