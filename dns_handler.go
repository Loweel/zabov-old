package main

import (
	"net"

	"github.com/miekg/dns"
)

func (mydns *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name

		if MyKillfile.Has(domain) {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP("127.0.0.1"),
			})
		} else {
			r := ForwardQuery(&msg)
			w.WriteMsg(r)
		}
	}
	w.WriteMsg(&msg)
}
