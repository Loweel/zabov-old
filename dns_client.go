package main

import (
	"fmt"
	"math/rand"
	"time"

	"strings"

	"github.com/miekg/dns"
)

//ForwardQuery forwards the query to the upstream server
//first server to answer wins
func ForwardQuery(query *dns.Msg) *dns.Msg {

	ZabovStats["ForwardQueries"]++

	r := new(dns.Msg)
	r.SetReply(query)
	r.Authoritative = true

	fqdn := strings.TrimRight(query.Question[0].Name, ".")

	lfqdn := fmt.Sprintf("%d", query.Question[0].Qtype) + "." + fqdn
	if MyCachefile.Has(lfqdn) {	
		incrementStats("CacheHit",1 )
		
		cached := GetDomainFromCache(lfqdn)
		cached.SetReply(query)
		cached.Authoritative = true
		return cached

	}

	upl := strings.Split(ZabovUpDNS, ",")
	
	c := new(dns.Client)
	rand.Seed(time.Now().Unix())
	for range upl {
		n := rand.Int() % len(upl)
		d := upl[n]
		
		
		in, _, err := c.Exchange(query, d)
		if err != nil {
			fmt.Printf("Problem with DNS %s : %s\n", d, err.Error())
			incrementStats("Problems " + d,1)
			
			continue
		} else {
			incrementStats(d,1)
			
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
