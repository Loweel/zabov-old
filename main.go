package main

import (
	"log"

	config "github.com/micro/go-config"
	"github.com/miekg/dns"
)

//MyDNS is my dns server
var MyDNS *dns.Server

//MyConf is the config file
var MyConf config.Config

//ZabovUpDNS keeps the name of upstream DNSs
var ZabovUpDNS string

//ZabovSingleBL list of urls returning a file with just names of domains
var ZabovSingleBL string

//ZabovDoubleBL list of urls returning a file with  IP<space>domain
var ZabovDoubleBL string

//ZabovAddBL is the IP we want to send all the clients to. Usually is 127.0.0.1
var ZabovAddBL string

type handler struct{}

func main() {

	MyDNS.Handler = &handler{}
	if err := MyDNS.ListenAndServe(); err != nil {
		log.Printf("Failed to set udp listener %s\n", err.Error())
	} else {
		log.Printf("Listener running \n")
	}
}
