package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"

	"github.com/boltdb/bolt"
)

var zabovKbucket = []byte("killfile")

//DomainKill stores a domain name inside the killfile
func DomainKill(s, durl string) {

	if len(s) > 2 {

		s = strings.ToLower(s)

		go writeInBolt(s, durl)

	}

}

func md5sum(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func writeInBolt(key, value string) {

	// store some data
	err := MyZabovDB.Batch(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(zabovKbucket)
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(key), []byte(value))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Println("Failed to write inside db: ", err.Error())
	}

}

func domainInKillfile(domain string) bool {

	var val []byte

	err := MyZabovDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(zabovKbucket)

		if bucket == nil {
			fmt.Printf("Bucket %s not found!", zabovKbucket)
			return nil
		}

		val = bucket.Get([]byte(domain))

		return nil
	})

	if err != nil {
		return false
	}

	return val != nil
}
