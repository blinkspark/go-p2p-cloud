package main

import (
	"log"
	"time"

	"github.com/blinkspark/go-p2p-cloud/client"
)

func main() {
	c, err := client.NewClient("priv.key")
	if err != nil {
		log.Panic(err)
	}
	go func() {
		for {
			time.Sleep(5 * time.Second)
			conns := c.Network().Conns()
			log.Printf("connected:%d", len(conns))
			// for _, con := range conns {
			// 	log.Println(con.RemotePeer().Pretty())
			// }
		}
	}()
	select {}
}
