package main

import (
	"fmt"
	"log"
	"strings"
	"strconv"

	"github.com/miekg/dns"
)

var records = map[string]string{
	"game.f1remoon.com.": "dig into x.game.f1remoon.com, where x = [0..9]\n_|_|_\n_|_|_\n_|_|_",

	"1.game.f1remoon.com.": "_|_|_\n_|_|_\nx|_|_",
	"2.game.f1remoon.com.": "_|_|_\n_|_|_\n_|x|_",
	"3.game.f1remoon.com.": "_|_|_\n_|_|_\n_|_|x",
	"4.game.f1remoon.com.": "_|_|_\nx|_|_\n_|_|_",
	"5.game.f1remoon.com.": "_|_|_\n_|x|_\n_|_|_",
	"6.game.f1remoon.com.": "_|_|_\n_|_|x\n_|_|_",
	"7.game.f1remoon.com.": "x|_|_\n_|_|_\n_|_|_",
	"8.game.f1remoon.com.": "_|x|_\n_|_|_\n_|_|_",
	"9.game.f1remoon.com.": "_|_|x\n_|_|_\n_|_|_",
}

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeTXT:
			log.Printf("Query for %s\n", q.Name)
			ans := records[q.Name];
			if ans != "" {
				for _, s := range strings.Split(ans, "\n") {
					t := &dns.TXT{
						Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 0},
						Txt: strings.Split(s, "|"),
					}
					m.Answer = append(m.Answer, t)
				}
			}
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			rr, err := dns.NewRR(fmt.Sprintf("%s A 127.0.0.1", q.Name))
			if err == nil {
				m.Answer = append(m.Answer, rr)
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func main() {
	// attach request handler func
	dns.HandleFunc("game.f1remoon.com.", handleDnsRequest)

	// start server
	port := 5353
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
