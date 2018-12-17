package main

import (
	"fmt"

	"github.com/peterbourgon/diskv"
)

//MyKillfile is the storage where we'll put domains to block
var MyKillfile *diskv.Diskv

func init() {

	flatTransform := func(s string) []string {

		var d []string
		d = append(d, s[0:1], s[1:2], s[2:3], s)
		fmt.Println("Returning: ", d)
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

	DomainKill("coglione.com")
	DomainKill("coglionissimo.com")

}

//DomainKill stores a domain name inside the killfile
func DomainKill(s string) {

	MyKillfile.WriteString(s, "127.0.0.1")

}
