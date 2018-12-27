package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/miekg/dns"

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
		fmt.Println("Cache folder created: ", MyCachefile.BasePath)
		MyCachefile.EraseAll()

	} else {
		fmt.Println("FAILED to create cache!")
	}

	

	go cacheCleanThread()

}

//DomainCache stores a domain name inside the cache
func DomainCache(s string, resp *dns.Msg) {

	a, err := resp.Pack()
	if err != nil {
		fmt.Println("Problems packing the response: ", err.Error())
	}

	if len(s) > 2 {

		MyCachefile.Write(strings.Trim(s, " "), a)

	}
}

//GetDomainFromCache stores a domain name inside the cache
func GetDomainFromCache(s string) *dns.Msg {

	ret := new(dns.Msg)

	record, _ := MyCachefile.Read(s)
	err := ret.Unpack(record)
	if err != nil {
		fmt.Println("Problem unpacking response: ", err.Error())
	}

	return ret

}

func cacheCleanThread() {
	fmt.Println("Starting updater of Cache, each ", ZabovCacheTTL)
	for {

		time.Sleep(time.Duration(ZabovCacheTTL) * time.Hour)
		MyCachefile.EraseAll()
	}

}
