package store

import (
	"log"
	"time"

	"github.com/tidwall/buntdb"
)

type StoreDB struct {
	cfg struct {
		key          string
		mode         string
		dbpath       string
		isPresistant bool
	}
	db *buntdb.DB
}

type StoreKV struct {
	key   string
	value string
	old   string
	err   error
}

// Begin opens a new transaction.
// Multiple read-only transactions can be opened at the same time but there can
// only be one read/write transaction at a time. Attempting to open a read/write
// transactions while another one is in progress will result in blocking until
// the current read/write transaction is completed.
//
// All transactions must be closed by calling Commit() or Rollback() when done.
//func (db *DB) Begin(writable bool) (*Tx, error) {
// Get returns a value for a key. If the item does not exist or if the item
// has expired then ErrNotFound is returned. If ignoreExpired is true, then
// the found value will be returned even if it is expired.
//func (tx *Tx) Get(key string, ignoreExpired ...bool) (val string, err error) {

func (st *StoreDB) GetKV(key string) (kv StoreKV) {
	tx, err := st.db.Begin(false)
	if err != nil {
		log.Println(err)
	}
	kv.value, err = tx.Get(key)
	if err != nil {
		log.Println(err)
	}
	if err := tx.Rollback(); err != nil {
		log.Println(err)
	}

	kv.key = key
	return kv
}

func (st *StoreDB) SetKV(key, value string, ttl time.Duration) (kv StoreKV) {
	tx, err := st.db.Begin(true)
	//defer tx.Commit()
	if err != nil {
		log.Println(err)
	}
	var opts *buntdb.SetOptions
	if ttl != 0 {
		opts.Expires = true
		opts.TTL = ttl
		// time.Duration
		//tx.Set("mykey", "myval", &buntdb.SetOptions{Expires:true, TTL:time.Second})
	}

	kv.old, _, kv.err = tx.Set(key, value, opts)
	if err := tx.Commit(); err != nil {
		log.Println(err)
	}

	kv.key = key
	kv.value = value
	return kv
}

func (st *StoreDB) DelKV(key string) (kv StoreKV) {
	tx, err := st.db.Begin(true)
	//defer tx.Commit()
	if err != nil {
		log.Println(err)
	}
	kv.old, kv.err = tx.Delete(key)
	if err := tx.Commit(); err != nil {
		log.Println(err)
	}
	kv.key = key
	return kv
}

func (st *StoreDB) IndexJSON(name, pattern, jpath string) error {
	// return s.DB.CreateIndex(name, pattern, buntdb.IndexJSON(jpath))
	return st.db.ReplaceIndex(name, pattern, buntdb.IndexJSON(jpath))
}

// func (s *StoreDB) ensureCollateJSON(name, pattern, langmode, jpath string) {
// 	err := s.DB.CreateIndex(name, pattern, collate.IndexJSON(langmode, jpath))
// 	if err != nil {
// 		log.Print("ERR :: Store :: ", s.cfg.dbpath, " :: ensureCollateJSON :: ", name, pattern, langmode, jpath, err)
// 	}
// 	//   Case insensitive: add the CI tag to the name
// 	//   Case sensitive: add the CS tag to the name
// 	//   For loosness: add the LOOSE tag to the name
// 	//   Ignores diacritics, case and weight
// 	//db.CreateIndex("last_name", "*", collate.IndexJSON("CHINESE_CI", "name.last"))
// }

// TruncateKVStore truncates entire store database
func (st *StoreDB) TruncateDB() error {
	tx, err := st.db.Begin(true)
	//defer tx.Commit()
	if err != nil {
		return err
	}
	if err := tx.DeleteAll(); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	if err := st.db.Shrink(); err != nil {
		return err
	}
	return nil
}

// Shrink will make the database file smaller by removing redundant
// log entries. This operation does not block the database.
func (st *StoreDB) ShrinkDB() error {
	return st.db.Shrink()
}

// Close releases all database resources.
// All transactions must be closed before closing the database.
func (st *StoreDB) CloseDB() error {
	return st.db.Close()
}
