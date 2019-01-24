package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
)

//MyKillfile is the storage where we'll put domains to block
var MyKillfile *bolt.DB

//MyLock to avoid having too many writes with no parallelism.
var MyLock = &sync.Mutex{}

func init() {

	var err error

	os.MkdirAll("./killfile", 0755)

	MyKillfile, err = bolt.Open("./killfile/killfile.db", 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("Problem Creating Killfile DB: ", err.Error())
	} else {
		fmt.Println("Killfile DB Created")

	}

}

//DomainKill stores a domain name inside the killfile
func DomainKill(s, durl string) {

	if len(s) > 2 {

		s = strings.ToLower(s)

		writeInBolt(s, durl)

	}

}

func md5sum(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func writeInBolt(key, value string) {

	MyLock.Lock()

	// store some data
	err := MyKillfile.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("killfile"))
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

	MyLock.Unlock()
}

func domainInKillfile(domain string) bool {

	var val []byte

	err := MyKillfile.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("killfile"))

		if bucket == nil {
			fmt.Println("Bucket killfile not found!")
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
