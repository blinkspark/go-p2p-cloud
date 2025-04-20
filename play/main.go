package main

import (
	"log"

	"github.com/blinkspark/go-p2p-cloud/client"
)

func main() {
	p2pClient, err := client.NewClient("config.json")
	if err != nil {
		log.Panic(err)
	}

	log.Println("Listening On:")
	for _, addr := range p2pClient.Addrs() {
		log.Printf("%s\n", addr)
	}

	select {}
}
