package main

import (
	"fmt"

	"strings"

	"github.com/miekg/dns"
)

//ForwardQuery forwards the query to the upstream server
//first server to answer wins
func ForwardQuery(query *dns.Msg) *dns.Msg {

	r := new(dns.Msg)
	r.SetReply(query)
	r.Authoritative = true

	fqdn := strings.TrimRight(query.Question[0].Name, ".")

	lfqdn := fmt.Sprintf("%d", query.Question[0].Qtype) + "." + fqdn
	if MyCachefile.Has(lfqdn) {

		fmt.Println("Cache hit: ", fqdn)
		cached := GetDomainFromCache(lfqdn)
		cached.SetReply(query)
		cached.Authoritative = true
		return cached

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
			in.SetReply(query)
			in.Authoritative = true
			DomainCache(lfqdn, in)
			return in

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
