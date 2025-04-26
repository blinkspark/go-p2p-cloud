package client

import (
	"github.com/dgraph-io/badger/v4"
)

type DataStore interface {
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Delete(key []byte) error
	// ListPrefix(prefix []byte) error
	Close() error
}

type BadgerDataStore struct {
	db *badger.DB
}

var _ DataStore = (*BadgerDataStore)(nil)

func NewBadgerDataStore(dsPath string, encKey []byte) (*BadgerDataStore, error) {
	badgerOpts := badger.DefaultOptions(dsPath)
	if encKey != nil {
		badgerOpts.IndexCacheSize = 64 << 20
		badgerOpts.EncryptionKey = encKey
	}
	db, err := badger.Open(badgerOpts)
	if err != nil {
		return nil, err
	}
	return &BadgerDataStore{db: db}, nil
}

func (ds *BadgerDataStore) Put(key []byte, value []byte) error {
	return ds.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (ds *BadgerDataStore) Get(key []byte) ([]byte, error) {
	var value []byte
	err := ds.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(nil)
		return err
	})
	return value, err
}

func (ds *BadgerDataStore) Delete(key []byte) error {
	return ds.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (ds *BadgerDataStore) ListPrefix(prefix []byte) ([][]byte, error) {
	keys := make([][]byte, 0)
	err := ds.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.KeyCopy(nil)
			keys = append(keys, k)
		}
		return nil
	})
	return keys, err
}

func (ds *BadgerDataStore) Close() error {
	return ds.db.Close()
}
