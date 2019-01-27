package main

import (
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
)

var zabovKbucket = []byte("killfile")

type killfileItem struct {
	Kdomain string
	Ksource string
}

var bChannel chan killfileItem

func init() {

	bChannel = make(chan killfileItem, 1024)
	fmt.Println("Initializing kill channel engine.")

	go bWriteThread()

}

func bWriteThread() {

	for item := range bChannel {

		writeInBolt(item.Kdomain, item.Ksource)
		go incrementStats("BL domains from "+item.Ksource+": ", 1)

	}

}

//DomainKill stores a domain name inside the killfile
func DomainKill(s, durl string) {

	if len(s) > 2 {

		s = strings.ToLower(s)

		var k killfileItem

		k.Kdomain = s
		k.Ksource = durl

		bChannel <- k

	}

}

func writeInBolt(key, value string) {

	// store some data
	err := MyZabovDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(zabovKbucket)

		if bucket == nil {
			fmt.Printf("Bucket %s not found!\n", zabovKbucket)
			return nil
		}
		berr := bucket.Put([]byte(key), []byte(value))
		if berr != nil {
			return berr
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
			fmt.Printf("Bucket %s not found!\n", zabovKbucket)
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
