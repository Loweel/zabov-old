package main

import (
	"fmt"

	config "github.com/micro/go-config"
	configfile "github.com/micro/go-config/source/file"
	"github.com/miekg/dns"
)

func init() {

	MyConf = config.NewConfig()

	MyConf.Load(configfile.NewSource(
		configfile.WithPath("./config.json"),
	))

	// now we read configuration file
	fmt.Println("Reading configuration file...")
	ZabovPort := MyConf.Get("zabov", "port").String("53")
	ZabovType := MyConf.Get("zabov", "proto").String("udp")
	ZabovAddr := MyConf.Get("zabov", "ipaddr").String("")
	ZabovUpDNS = MyConf.Get("zabov", "upstream").String("127.0.0.1")
	ZabovSingleBL = MyConf.Get("zabov", "singlefilters").String("")
	ZabovDoubleBL = MyConf.Get("zabov", "doublefilters").String("")
	ZabovAddBL = MyConf.Get("zabov", "blackholeip").String("127.0.0.1")
	ZabovCacheTTL = MyConf.Get("zabov", "cachettl").Int(12)
	ZabovHostsFile =  MyConf.Get("zabov","hostsfile").String("./hosts.txt")

	zabovString := ZabovAddr + ":" + ZabovPort

	MyDNS = new(dns.Server)
	MyDNS.Addr = zabovString
	MyDNS.Net = ZabovType

	MyConf.Close()

}
