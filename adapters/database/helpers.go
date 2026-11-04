package database

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/dgraph-io/badger/v4"
	"github.com/pedrodbflores/personal-finance-tracker/app_errors"
)

func IsKeyPresent(db *badger.DB, key []byte) (bool, error) {
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		return err
	})
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetByKey[T any](db *badger.DB, key []byte) (T, error) {
	var output T
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				typeName := reflect.TypeOf((*T)(nil)).Elem().Name()
				return app_errors.NewResourceNotFoundError(typeName, string(key))
			}
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &output)
		})
	})
	if err != nil {
		return output, err
	}
	return output, nil
}

func ListAllByKey[T any](db *badger.DB, prefix []byte) ([]T, error) {
	var output []T
	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			Prefix: []byte(prefix),
		})
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				var out T
				json.Unmarshal(val, &out)
				output = append(output, out)
				return nil
			})
		}
		return nil
	})
	if err != nil {
		return []T{}, err
	}
	return output, nil
}
