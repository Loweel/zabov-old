package main

import (
	"fmt"
	"net"

	"strings"

	"github.com/miekg/dns"
)

//ForwardQuery forwards the query to the upstream server
//first server to answer wins
func ForwardQuery(query *dns.Msg) *dns.Msg {

	r := new(dns.Msg)

	// initializing r, just in case no dns works.
	r.Answer = append(r.Answer, &dns.A{
		Hdr: dns.RR_Header{Name: query.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
		A:   net.ParseIP("127.0.0.1"),
	})

	upl := strings.Split(ZabovUpDNS, ",")
	fmt.Println("Servers: ", upl)
	c := new(dns.Client)

	for _, d := range upl {
		fmt.Println("Query DNS: ", d)
		in, _, err := c.Exchange(query, d)
		if err != nil {
			fmt.Printf("Problem with DNS %s : %s\n", d, err.Error())
			continue
		} else {
			r = in
			return r
		}
	}
	return r

}

func init() {

	fmt.Println("DNS client engine starting")
	m := new(dns.Msg)
	// RFC2606 test domain, should always work.
	m.SetQuestion(dns.Fqdn("example.com"), dns.TypeA)
	ForwardQuery(m)
	fmt.Println("DNS client tested")
}
