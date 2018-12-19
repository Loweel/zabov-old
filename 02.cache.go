package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/peterbourgon/diskv"
)

//MyCachefile is the storage where we'll put domains to block
var MyCachefile *diskv.Diskv

func init() {

	flatTransform := func(s string) []string {

		var d []string

		usen := strings.Split(s, ".")
		inverse := reverse(usen)
		d = append(d, inverse...)
		return d

	}

	MyCachefile = diskv.New(diskv.Options{
		BasePath:     "cache",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})

	if MyCachefile != nil {
		fmt.Println("Cache folder created", MyCachefile.BasePath)

	} else {
		fmt.Println("FAILED to create cache!")
	}

	go cacheCleanThread()

}

//DomainCache stores a domain name inside the cache
func DomainCache(s, addr string) {

	if len(s) > 2 {

		MyCachefile.WriteString(strings.Trim(s, " "), addr)

	}
}

//GetDomainFromCache stores a domain name inside the cache
func GetDomainFromCache(s string) string {

	record := MyCachefile.ReadString(s)

	return record

}

func cacheCleanThread() {
	fmt.Println("Starting updater of Cache")
	for {

		time.Sleep(12 * time.Hour)
		MyCachefile.EraseAll()
	}

}
