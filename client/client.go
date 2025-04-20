package client

import (
	"context"
	"log"

	"github.com/blinkspark/go-p2p-cloud/config"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
)

type Client struct {
	host      host.Host
	dhtNode   *dht.IpfsDHT
	discovery *routing.RoutingDiscovery
}

func (c *Client) Addrs() []string {
	id := c.host.ID()
	addrs := c.host.Addrs()
	addrStrings := make([]string, len(addrs))
	for i, addr := range addrs {
		addrStrings[i] = addr.String() + "/p2p/" + id.String()
	}
	return addrStrings
}

func (c *Client) Close() error {
	err := c.host.Close()
	if err != nil {
		return err
	}
	err = c.dhtNode.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewClient(configPath string) (*Client, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	options, err := config.BuildOptionFromConfig(cfg)
	if err != nil {
		return nil, err
	}

	h, err := libp2p.New(options...)
	if err != nil {
		return nil, err
	}

	dhtNode, err := dht.New(context.Background(), h)
	if err != nil {
		return nil, err
	}
	discovery := routing.NewRoutingDiscovery(dhtNode)

	for _, addr := range dht.DefaultBootstrapPeers {
		pi, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			log.Printf("Failed to parse address: %v\n", err)
			continue
		}
		go func() {
			err = h.Connect(context.Background(), *pi)
			if err != nil {
				log.Printf("Failed to connect to peer %s: %v\n", addr, err)
			} else {
				log.Printf("Connected to peer %s\n", addr)
			}
		}()
	}

	return &Client{
		host:      h,
		dhtNode:   dhtNode,
		discovery: discovery,
	}, nil
}
