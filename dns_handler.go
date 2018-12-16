package main

import (
	"net"

	"github.com/miekg/dns"
)

// tmp. This is a test, just to experiment with dns libraries
var domainsToAddresses = map[string]string{
	"google.com.": "1.2.3.4",
	"puppa.qui.":  "151.152.153.154",
}

func (mydns *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		address, ok := domainsToAddresses[domain]
		if ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(address),
			})
		}
	}
	w.WriteMsg(&msg)
}
