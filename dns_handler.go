package main

import (
	"net"
	"strings"

	"github.com/miekg/dns"
)

func (mydns *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	go incrementStats("TotalQueries", 1)

	remIP, _, e := net.SplitHostPort(w.RemoteAddr().String())
	if e != nil {
		go incrementStats("CLIENT: "+remIP, 1)
	}

	msg := dns.Msg{}
	msg.SetReply(r)

	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		fqdn := strings.TrimRight(domain, ".")

		if domainInKillfile(fqdn) {
			go incrementStats("Killed", 1)

			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(ZabovAddBL),
			})
		} else {
			ret := ForwardQuery(r)
			w.WriteMsg(ret)
		}
	default:
		ret := ForwardQuery(r)
		w.WriteMsg(ret)
	}
	w.WriteMsg(&msg)

}
