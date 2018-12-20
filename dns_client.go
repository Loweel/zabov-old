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
	r.SetReply(query)

	fqdn := strings.TrimRight(query.Question[0].Name, ".")

	if MyCachefile.Has(fqdn) {
		r.Authoritative = true
		fmt.Println("Cache hit: ", fqdn)
		cached := GetDomainFromCache(fqdn)
		fmt.Println("Cached IP: ", cached)
		if net.ParseIP(cached) != nil {
			r.Answer = append(r.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: query.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(cached),
			})
			return r
		}
	}

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
			if query.Question[0].Qtype == dns.TypeA {
				ip := r.Answer[0].(*dns.A)
				DomainCache(fqdn, ip.A.String())
			}
			return r
		}
	}
	return r

}

func init() {

	fmt.Println("DNS client engine starting")
	m := new(dns.Msg)
	// RFC2606 test domain, should always work, unless internet is down.
	m.SetQuestion(dns.Fqdn("example.com"), dns.TypeA)
	ForwardQuery(m)
	fmt.Println("DNS client tested")
}
