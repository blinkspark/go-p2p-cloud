package main

import (
	"context"
	"log"
	"time"

	"github.com/blinkspark/go-p2p-cloud/server"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

func main() {
	play_relay()
}

func play_relay() {
	// /ip4/192.243.120.164/tcp/62233/p2p/12D3KooWQPZwkxdHSE8sKkokkGDX8nSAHD9RjiCc1WKf8vg6e7MW
	h, err := libp2p.New(libp2p.EnableHolePunching())
	if err != nil {
		log.Panic(err)
	}
	addr, err := peer.AddrInfoFromString("/ip4/192.243.120.164/tcp/62233/p2p/12D3KooWQPZwkxdHSE8sKkokkGDX8nSAHD9RjiCc1WKf8vg6e7MW")
	if err != nil {
		log.Panic(err)
	}
	err = h.Connect(context.Background(), *addr)
	if err != nil {
		log.Panic(err)
	}
	ps, err := pubsub.NewGossipSub(context.Background(), h)
	if err != nil {
		log.Panic(err)
	}
	t, err := ps.Join("/nealfree.ml/ps")
	if err != nil {
		log.Panic(err)
	}
	sub, err := t.Subscribe()
	if err != nil {
		log.Panic(err)
	}
	go func() {
		for {
			msg, err := sub.Next(context.Background())
			if err != nil {
				log.Println(err)
				continue
			}
			log.Printf("from:%s, msg:%s\n", msg.GetFrom().Pretty(), string(msg.GetData()))
			from := msg.GetFrom().Pretty()
			if from == h.ID().Pretty() {
				continue
			}
			addr := "/p2p/12D3KooWQPZwkxdHSE8sKkokkGDX8nSAHD9RjiCc1WKf8vg6e7MW/p2p-circuit/p2p/" + from
			pi, err := peer.AddrInfoFromString(addr)
			if err != nil {
				log.Println(err)
				continue
			}
			err = h.Connect(context.Background(), *pi)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("---dial successfully---")
		}
	}()
	go func() {
		for {
			time.Sleep(5 * time.Second)
			err = t.Publish(context.Background(), []byte("hello"))
			if err != nil {
				log.Println(err)
				continue
			}
		}
	}()

	for {
		time.Sleep(time.Second * 5)
		log.Println(len(h.Network().Conns()))
		for _, p := range h.Network().Conns() {
			log.Printf("%#+v\n", p.RemotePeer().Pretty())
		}
		for _, addr := range h.Addrs() {
			log.Printf("%#+v\n", addr.String())
		}
	}
}

func play_s() {
	s, err := server.NewServer("test.key", 62233, "")
	if err != nil {
		log.Panic(err)
	}
	// err = s.Bootstrap()
	// if err != nil {
	// 	log.Panic(err)
	// }
	// log.Printf("%#+v,%s\n", s, s.ID())

	for {
		time.Sleep(time.Second * 5)
		log.Println(len(s.Peerstore().Peers()))
		for _, p := range s.Peerstore().Peers() {
			log.Printf("%#+v\n", p.String())
		}
		for _, addr := range s.Addrs() {
			log.Printf("%#+v\n", addr.String())
		}
	}
}
