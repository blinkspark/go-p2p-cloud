package main

import (
	"log"

	"github.com/blinkspark/go-p2p-cloud/server"
)

func main() {
	s, err := server.NewServer("test.key")
	if err != nil {
		log.Panic(err)
	}
	log.Printf("%#+v,%s\n", s, s.ID())
	for _, addr := range s.Addrs() {
		log.Printf("%#+v\n", addr.String())
	}
}
