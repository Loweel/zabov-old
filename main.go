package main

import (
	"fmt"
	"log"

	config "github.com/micro/go-config"
	configfile "github.com/micro/go-config/source/file"
	"github.com/miekg/dns"
)

//MyDNS is my dns server
var MyDNS *dns.Server

//MyConf is the config file
var MyConf config.Config

type handler struct{}

func init() {

	MyConf = config.NewConfig()

	MyConf.Load(configfile.NewSource(
		configfile.WithPath("./config.json"),
	))

	// now we read configuration file
	fmt.Println("Reading configuration file...")
	ZabovPort := MyConf.Get("zabov", "port").String("53")
	ZabovType := MyConf.Get("zabov", "proto").String("udp")
	ZabovAddr := MyConf.Get("zabov", "ipaddr").String("127.0.0.1")

	zabovString := ZabovAddr + ":" + ZabovPort

	MyDNS = new(dns.Server)
	MyDNS.Addr = zabovString
	MyDNS.Net = ZabovType

	MyConf.Close()

}

func main() {

	MyDNS.Handler = &handler{}
	if err := MyDNS.ListenAndServe(); err != nil {
		log.Printf("Failed to set udp listener %s\n", err.Error())
	} else {
		log.Printf("Listener running \n")
	}
}
