package db

import (
	"github.com/boltdb/bolt"
)

var g_db *bolt.DB

func init() {
	g_db, _ = bolt.Open("dns.db", 0600, nil)
}

func Save(table string, k string, v string) {
	g_db.Update(func(tx *bolt.Tx) error {
		b, _ := tx.CreateBucketIfNotExists([]byte(table))
		b.Put([]byte(k), []byte(v))
		return nil
	})
}
func Get(table string, k string) string {
	var data string
	g_db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(table))
		if b != nil {
			data = string(b.Get([]byte(k)))
		}
		return nil
	})
	return data
}
