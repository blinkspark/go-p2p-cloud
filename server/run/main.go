package main

import (
	"flag"
	"log"

	"github.com/blinkspark/go-p2p-cloud/server"
)

var (
	port     int
	keyPath  string
	protocol string
)

func init() {
	flag.IntVar(&port, "p", 62233, "-p {PORT}")
	flag.StringVar(&keyPath, "k", "priv.key", "-k {KEY_PATH}")
	flag.StringVar(&protocol, "P", "/nealfree.ml/relay/v0.1.0", "-P {PROTOCOL}")
	flag.Parse()
}

func main() {
	s, err := server.NewServer(keyPath, port)
	if err != nil {
		log.Panic(err)
	}
	sid := s.ID()
	for _, addr := range s.Addrs() {
		log.Printf("%s/p2p/%s", addr.String(), sid.Pretty())
	}
	select {}
}
