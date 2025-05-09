package main

import (
	"log"

	"github.com/blinkspark/go-p2p-cloud/client"
)

func playDS() {
	ds, err := client.NewBadgerDataStore("tmp", nil)
	if err != nil {
		log.Panic(err)
	}
	defer ds.Close()
	// err = ds.Put([]byte("hello.1"), []byte("value1"))
	// if err != nil {
	// 	log.Panic(err)
	// }
	// err = ds.Put([]byte("hello.2"), []byte("value2"))
	// if err != nil {
	// 	log.Panic(err)
	// }
	// err = ds.Put([]byte("hello"), []byte("value3"))
	// if err!= nil {
	// 	log.Panic(err)
	// }
	keys, err := ds.ListPrefix([]byte("hello"))
	if err != nil {
		log.Panic(err)
	}
	for _, key := range keys {
		log.Println(string(key))
	}
}
