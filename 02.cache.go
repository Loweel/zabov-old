package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/miekg/dns"
)

var zabovCbucket = []byte("cache")

type cacheItem struct {
	query []byte
	date  time.Time
}

//DomainCache stores a domain name inside the cache
func DomainCache(s string, resp *dns.Msg) {

	var domain2cache cacheItem
	var err error
	var dom2 bytes.Buffer
	enc := gob.NewEncoder(&dom2)

	domain2cache.query, err = resp.Pack()
	if err != nil {
		fmt.Println("Problems packing the response: ", err.Error())
	}
	domain2cache.date = time.Now()

	err = enc.Encode(domain2cache)

	if err != nil {
		fmt.Println("Cannot GOB the domain to cache", err.Error())
	}

	cacheInBolt(s, dom2.Bytes())

}

func cacheInBolt(key string, domain []byte) {

	MyZabovLock.Lock()

	// store some data
	err := MyZabovDB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(zabovCbucket)
		if err != nil {

			return err
		}

		err = bucket.Put([]byte(key), domain)
		if err != nil {

			return err
		}
		return nil
	})

	if err != nil {
		fmt.Println("Failed to write inside db: ", err.Error())
	}

	MyZabovLock.Unlock()
}

//GetDomainFromCache stores a domain name inside the cache
func GetDomainFromCache(s string) *dns.Msg {

	ret := new(dns.Msg)
	var cache bytes.Buffer
	dec := gob.NewDecoder(&cache)
	var record cacheItem
	var conf []byte

	MyZabovLock.Lock()

	if domainInCache(s) != true {
		MyZabovLock.Unlock()
		return nil
	}

	if err := MyZabovDB.View(func(tx *bolt.Tx) error {
		conf = tx.Bucket(zabovCbucket).Get([]byte(s))

		return nil
	}); err != nil {
		fmt.Println("Error getting data from cache: ", err.Error())
		return nil
	}

	MyZabovLock.Unlock()

	if conf == nil {
		return nil
	}

	fmt.Println("Conf: ", conf)

	cache.Write(conf)

	err := dec.Decode(&record)
	if err != nil {
		fmt.Println("Decode error :", err.Error())
		return nil
	}

	if time.Since(record.date) > (time.Duration(ZabovCacheTTL) * time.Hour) {
		return nil
	}

	err = ret.Unpack(record.query)
	if err != nil {
		fmt.Println("Problem unpacking response: ", err.Error())
		return nil
	}

	return ret

}



func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func domainInCache(domain string) bool {

	var val []byte

	err := MyZabovDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(zabovCbucket)

		if bucket == nil {
			fmt.Printf("Bucket %s not found!", zabovCbucket)
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

