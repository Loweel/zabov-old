package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/miekg/dns"

	"github.com/peterbourgon/diskv"
)

//MyCachefile is the storage where we'll put domains to block
var MyCachefile *diskv.Diskv

const myCachePath = "cache"

func init() {

	flatTransform := func(s string) []string {

		var d []string

		d = append(d, hourPrint())

		usen := strings.Split(s, ".")
		inverse := reverse(usen)
		d = append(d, inverse...)
		return d

	}

	MyCachefile = diskv.New(diskv.Options{
		BasePath:     myCachePath,
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})

	if MyCachefile != nil {
		fmt.Println("Cache folder created: ", MyCachefile.BasePath)

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
		cleanCache()
	}

}

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func hourPrint() (htamp string) {

	t := time.Now()

	htamp = fmt.Sprintf("%s", t.Format("2006010215"))

	return

}

func printCache() {

	fmt.Println("Printing Cache Keys")
	for k := range MyCachefile.Keys(nil) {
		fmt.Println(k)
	}

}

func cleanCache() {

	var cutoff = 61 * time.Minute

	fileInfo, err := ioutil.ReadDir(myCachePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	now := time.Now()
	for _, info := range fileInfo {
		if diff := now.Sub(info.ModTime()); diff > cutoff {
			if removeErr := os.RemoveAll(info.Name()); removeErr == nil {
				fmt.Printf("Deleted %s which is %s old\n", info.Name(), diff)
			} else {
				fmt.Printf("Error removing %s: %s", info.Name(), removeErr.Error())
			}
		}
	}

}
