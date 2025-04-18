package main

import (
	"context"
	"log"
	"time"

	"github.com/blinkspark/go-p2p-cloud/config"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Panicf("Failed to load config: %v", err)
	}

	log.Printf("%#+v", cfg)

	options, err := config.BuildOptionFromConfig(cfg)
	if err != nil {
		log.Panicf("Failed to build options from config: %v", err)
	}

	h, err := libp2p.New(options...)
	if err != nil {
		log.Panicf("Failed to create host: %v", err)
	}

	dhtNode, err := dht.New(context.Background(), h)
	if err != nil {
		log.Panicf("Failed to create DHT: %v", err)
	}

	// connect to bootstrap peers
	for _, addr := range dht.DefaultBootstrapPeers {
		pi, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			log.Printf("Failed to parse address: %v\n", err)
			continue
		}
		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()
		go func() {
			err = h.Connect(context.Background(), *pi)
			if err != nil {
				log.Printf("Failed to connect to peer %s: %v\n", addr, err)
			} else {
				log.Printf("Connected to peer %s\n", addr)
			}
		}()
	}

	err = dhtNode.Bootstrap(context.Background())
	if err != nil {
		log.Panicf("Failed to bootstrap DHT: %v", err)
	}

	log.Println("Host created: ", h.ID())
	log.Println("Host listening on: ", h.Addrs())

	time.Sleep(5 * time.Second)
	routingDiscovery := routing.NewRoutingDiscovery(dhtNode)
	dutil.Advertise(context.Background(), routingDiscovery, "/p2p/peer-discovery/1.0.0")
	time.Sleep(5 * time.Second)
	peers, err := dutil.FindPeers(context.Background(), routingDiscovery, "/p2p/peer-discovery/1.0.0")
	if err != nil {
		log.Printf("Failed to find peers: %v\n", err)
		return
	}
	log.Printf("Found peers: %v\n", peers)

	select {}
}
