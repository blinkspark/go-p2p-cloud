package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/blinkspark/go-p2p-cloud/client"
)

const (
	PROTOCOL = "nealfree.ml/test/v0.1.0"
)

var (
	keyPath string
)

func init() {
	flag.StringVar(&keyPath, "k", "priv.key", "-k {KEY_PATH}")
	flag.Parse()
}

func main() {
	c, err := client.NewClient(keyPath)
	if err != nil {
		log.Panic(err)
	}
	c.Advertise(PROTOCOL)
	go func() {
		for {
			time.Sleep(5 * time.Second)
			conns := c.Network().Conns()
			log.Printf("connected:%d", len(conns))
			// for _, con := range conns {
			// 	log.Println("conn: ", con.RemotePeer().Pretty())
			// }
		}
	}()
	go func() {
		for {
			time.Sleep(5 * time.Second)
			pc, err := c.RoutingDiscovery.FindPeers(context.Background(), PROTOCOL)
			if err != nil {
				log.Println(err)
			}
			for p := range pc {
				log.Println(p.String())
			}
		}
	}()
	select {}
}
