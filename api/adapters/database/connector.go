package database

import (
	"os"

	"github.com/dgraph-io/badger/v4"
)

func GetDatabase(opts *badger.Options) (*badger.DB, error) {
	var db *badger.DB
	var err error
	if opts == nil {
		wd, _ := os.Getwd()
		db, err = badger.Open(badger.DefaultOptions(wd))
	} else {
		db, err = badger.Open(*opts)
	}
	if err != nil {
		return nil, err
	}
	return db, err
}
