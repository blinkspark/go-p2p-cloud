package main

import (
	"flag"
	"log"

	"github.com/blinkspark/go-p2p-cloud/server"
)

var (
	port    int
	keyPath string
	topic   string
)

func init() {
	flag.IntVar(&port, "p", 62233, "-p {PORT}")
	flag.StringVar(&keyPath, "k", "priv.key", "-k {KEY_PATH}")
	flag.StringVar(&topic, "t", "/nealfree.ml/ps", "-t {TOPIC}")
	flag.Parse()
}

func main() {
	s, err := server.NewServer(keyPath, port, topic)
	if err != nil {
		log.Panic(err)
	}
	log.Println(s.ID().Pretty())
	for _, addr := range s.Addrs() {
		log.Println(addr.String())
	}
	select {}
}
