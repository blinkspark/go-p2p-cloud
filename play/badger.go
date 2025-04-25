package main

import (
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v4"
)

func playBadger() {
	db, err := badger.Open(badger.DefaultOptions("tmp").WithIndexCacheSize(64 << 20))
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	err = db.Update(func(txn *badger.Txn) error {
		err = txn.Set([]byte("hello"), []byte("world"))
		if err != nil {
			err = fmt.Errorf("txn.Set: %w", err)
			return err
		}

		return err
	})
	if err != nil {
		log.Panic(err)
	}
	err = db.View(func(txn *badger.Txn) error {
		var itm *badger.Item
		itm, err = txn.Get([]byte("hello"))
		if err != nil {
			err = fmt.Errorf("txn.Get: %w", err)
			return err
		}
		valSize := itm.ValueSize()
		log.Printf("valSize: %d", valSize)
		var valCopy []byte = make([]byte, valSize)
		var val []byte
		val, err = itm.ValueCopy(valCopy)
		if err != nil {
			err = fmt.Errorf("itm.ValueCopy: %w", err)
			return err
		}
		log.Println("var:" + string(val))
		log.Println("var:" + string(valCopy))
		log.Printf("valcap: %d", cap(val))
		log.Printf("valcopycap: %d", cap(valCopy))
		return err
	})
}
