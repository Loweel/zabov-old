package main

import (
	"fmt"
	"strings"

	"github.com/peterbourgon/diskv"
)

//MyKillfile is the storage where we'll put domains to block
var MyKillfile *diskv.Diskv

func init() {

	flatTransform := func(s string) []string {

		var d []string

		usen := strings.Split(s, ".")
		inverse := reverse(usen)
		d = append(d, inverse...)
		return d

	}

	MyKillfile = diskv.New(diskv.Options{
		BasePath:     "killfile",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})

	if MyKillfile != nil {
		fmt.Println("Killfile folder created", MyKillfile.BasePath)

	} else {
		fmt.Println("FAILED to create queue!")
	}

}

//DomainKill stores a domain name inside the killfile
func DomainKill(s, durl string) {

	if len(s) > 2 {

		MyKillfile.WriteString(strings.Trim(s, " "), durl)

	}
}

func reverse(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
