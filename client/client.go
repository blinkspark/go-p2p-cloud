package client

import (
	"context"
	"crypto/rand"
	"log"
	"sync"

	"github.com/blinkspark/go-p2p-cloud/key"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
)

type Client struct {
	host.Host
	*dht.IpfsDHT
	*drouting.RoutingDiscovery
	relays []peer.AddrInfo
}

func NewClient(keyPath string) (*Client, error) {
	// TODO add a function to get static relays
	addr, err := peer.AddrInfoFromString("/ip4/192.243.120.164/tcp/62233/p2p/12D3KooWQPZwkxdHSE8sKkokkGDX8nSAHD9RjiCc1WKf8vg6e7MW")
	if err != nil {
		log.Panic(err)
	}
	staticRelays := []peer.AddrInfo{*addr}
	priv, err := key.ReadKey(keyPath)
	if err != nil {
		priv, _, err = crypto.GenerateEd25519Key(rand.Reader)
		if err != nil {
			return nil, err
		}
		err = key.WriteKey(keyPath, priv)
		if err != nil {
			return nil, err
		}
	}

	// create host
	h, err := libp2p.New(
		libp2p.Identity(priv), libp2p.EnableHolePunching(), libp2p.NATPortMap(),
		libp2p.EnableAutoRelay(autorelay.WithStaticRelays(staticRelays)),
		libp2p.EnableNATService(), libp2p.EnableRelayService())
	if err != nil {
		return nil, err
	}

	// create dht and routing
	dhtNode, err := dht.New(context.Background(), h)
	if err != nil {
		return nil, err
	}
	routing := drouting.NewRoutingDiscovery(dhtNode)

	// make client
	c := &Client{Host: h, IpfsDHT: dhtNode, RoutingDiscovery: routing, relays: staticRelays}
	// bootstrap
	err = c.bootstrap()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) bootstrap() error {
	bsPeers, err := peer.AddrInfosFromP2pAddrs(dht.DefaultBootstrapPeers...)
	if err != nil {
		return err
	}
	bsPeers = append(bsPeers, c.relays...)
	var wg sync.WaitGroup
	for _, p := range bsPeers {
		wg.Add(1)
		go func(p *peer.AddrInfo) {
			defer wg.Done()
			err := c.Connect(context.Background(), *p)
			if err != nil {
				log.Println(err)
			}
		}(&p)
	}
	err = c.Bootstrap(context.Background())
	if err != nil {
		return err
	}
	return nil
}
