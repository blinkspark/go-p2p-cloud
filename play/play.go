package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/blinkspark/go-p2p-cloud/server"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
)

func main() {
	play_relay()
}

func play_relay() {
	// /ip4/192.243.120.164/tcp/62233/p2p/12D3KooWQPZwkxdHSE8sKkokkGDX8nSAHD9RjiCc1WKf8vg6e7MW
	addr, err := peer.AddrInfoFromString("/ip4/192.243.120.164/tcp/62233/p2p/12D3KooWQPZwkxdHSE8sKkokkGDX8nSAHD9RjiCc1WKf8vg6e7MW")
	if err != nil {
		log.Panic(err)
	}
	h, err := libp2p.New(
		libp2p.EnableHolePunching(), libp2p.EnableNATService(), libp2p.EnableRelayService(),
		libp2p.EnableAutoRelay(autorelay.WithStaticRelays([]peer.AddrInfo{*addr})))
	if err != nil {
		log.Panic(err)
	}
	err = h.Connect(context.Background(), *addr)
	if err != nil {
		log.Panic(err)
	}
	// _, err = client.Reserve(context.Background(), h, *addr)
	// if err != nil {
	// 	log.Panic(err)
	// }
	var dhtNode *dht.IpfsDHT
	h, err := libp2p.New(libp2p.NATPortMap(), libp2p.EnableHolePunching(), libp2p.EnableAutoRelay(autorelay.WithDefaultStaticRelays()), libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
		var err error
		dhtNode, err = dht.New(context.Background(), h)
		return dhtNode, err
	}))
	if err != nil {
		log.Panic(err)
	}

	pis, err := peer.AddrInfosFromP2pAddrs(dht.DefaultBootstrapPeers...)
	if err != nil {
		log.Panic(err)
	}
	var wg sync.WaitGroup
	for _, pi := range pis {
		wg.Add(1)
		go func(pi *peer.AddrInfo) {
			err := h.Connect(context.Background(), *pi)
			if err != nil {
				log.Println(err)
			}
			wg.Done()
		}(&pi)
	}
	wg.Wait()

	err = dhtNode.Bootstrap(context.Background())
	if err != nil {
		log.Panic(err)
	}
	disc := drouting.NewRoutingDiscovery(dhtNode)
	peers, err := disc.FindPeers(context.Background(), "/nealfree.ml/relay/v0.1.0")
	if err != nil {
		log.Panic(err)
	}
	for p := range peers {
		log.Println(p.ID.Pretty())
		log.Println(p.Addrs)
		cnss := h.Network().Connectedness(p.ID)
		if cnss == network.Connected || cnss == network.CannotConnect {
			continue
		}
		err = h.Connect(context.Background(), p)
		if err != nil {
			log.Println("--connect error:", err)
			continue
		}
		log.Println(p.ID.Pretty(), " connected")
	}
}

func play_s() {
	s, err := server.NewServer("test.key", 62233)
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
		log.Println(len(s.Network().Conns()))
		for _, p := range s.Network().Conns() {
			log.Printf("%#+v\n", p.RemotePeer().Pretty())
		}
		for _, addr := range s.Addrs() {
			log.Printf("%#+v\n", addr.String())
		}
	}
}
