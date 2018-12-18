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
		le := fmt.Sprintf("%d", len(s))
		fi := fmt.Sprintf("%d", strings.Count(s, "."))
		en := string(s[len(s)-1])
		d = append(d, s[0:1], le, en, fi, s)
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
func DomainKill(s string) {

	if len(s) > 2 {

		MyKillfile.WriteString(s, "127.0.0.1")

	}
}
