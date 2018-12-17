package main

import (
	"fmt"

	"strings"

	"github.com/miekg/dns"
)

//ForwardQuery forwards the query to the upstream server
func ForwardQuery(query *dns.Msg) *dns.Msg {

	r := new(dns.Msg)

	upl := strings.Split(ZabovUpDNS, ",")
	fmt.Println("Servers: ", upl)
	c := new(dns.Client)

	for _, d := range upl {
		fmt.Println("DNS: ", d)
		in, _, err := c.Exchange(query, d)
		if err != nil {
			fmt.Printf("Problem with DNS %s : %s\n", d, err.Error())
			continue
		} else {
			r = in
		}
	}
	return r

}

func init() {

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)

	ForwardQuery(m)
}
