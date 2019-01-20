package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"strings"

	"github.com/peterbourgon/diskv"
)

//MyKillfile is the storage where we'll put domains to block
var MyKillfile *diskv.Diskv

func init() {

	flatTransform := func(s string) []string {

		var d []string

		h := md5sum(s)

		d = append(d, h[0:4])

		return d

	}

	MyKillfile = diskv.New(diskv.Options{
		BasePath:     "killfile",
		Transform:    flatTransform,
		CacheSizeMax: 2048 * 2048,
	})

	if MyKillfile != nil {
		fmt.Println("Killfile folder created: ", MyKillfile.BasePath)

	} else {
		fmt.Println("FAILED to create queue!")
	}

}

//DomainKill stores a domain name inside the killfile
func DomainKill(s, durl string) {

	if len(s) > 2 {

		s = strings.ToLower(s)

		MyKillfile.WriteString(strings.Trim(s, " "), durl)

	}
}

func md5sum(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
