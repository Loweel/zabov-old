package main

import (
	"fmt"
	"time"

	"math/rand"
	"strings"

	"github.com/miekg/dns"
)

//ForwardQuery forwards the query to the upstream server
//first server to answer wins
func ForwardQuery(query *dns.Msg) *dns.Msg {

	go incrementStats("ForwardQueries", 1)

	r := new(dns.Msg)
	r.SetReply(query)
	r.Authoritative = true

	fqdn := strings.TrimRight(query.Question[0].Name, ".")

	lfqdn := fmt.Sprintf("%d", query.Question[0].Qtype) + "." + fqdn
	if MyCachefile.Has(lfqdn) {
		go incrementStats("CacheHit", 1)

		cached := GetDomainFromCache(lfqdn)
		cached.SetReply(query)
		cached.Authoritative = true
		return cached

	}

	c := new(dns.Client)

	for {
		// round robin with retry

		d := oneTimeDNS()

		in, _, err := c.Exchange(query, d)
		if err != nil {
			fmt.Printf("Problem with DNS %s : %s\n", d, err.Error())
			go incrementStats("DNS Problems "+d, 1)
			continue
		} else {
			go incrementStats(d, 1)
			in.SetReply(query)
			in.Authoritative = true
			DomainCache(lfqdn, in)
			return in

		}

	}

}

func init() {

	fmt.Println("DNS client engine starting")
	m := new(dns.Msg)
	// RFC2606 test domain, should always work, unless internet is down.
	m.SetQuestion(dns.Fqdn("example.com"), dns.TypeA)
	ForwardQuery(m)
	fmt.Println("DNS client tested")
	printCache()
}

func oneTimeDNS() (dns string) {

	rand.Seed(time.Now().Unix())

	upl := fileByLines(ZabovUpDNS)

	n := rand.Int() % len(upl)

	dns = upl[n]

	return

}
