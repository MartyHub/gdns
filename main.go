package main

import (
	"log"

	"github.com/MartyHub/gdns/dns"
)

func main() {
	for _, domain := range []string{
		"www.example.com",
		"twitter.com",
		"www.facebook.com",
		"google.com",
		"www.metafilter.com",
	} {
		ip, err := dns.NewResolver(domain).Resolve()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("=> %s = %s\n\n", domain, ip.String())
	}
}
