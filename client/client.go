package client

import (
	"context"
	"log"
	"time"

	"github.com/blinkspark/go-p2p-cloud/config"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Client struct {
	host       host.Host
	dhtNode    *dht.IpfsDHT
	discovery  *routing.RoutingDiscovery
	pubsub     *pubsub.PubSub
	ds         DataStore
	configPath string
	config     *config.HostConfig
}

func (c *Client) Advertise(serviceName string) {
	go func() {
		_, err := c.discovery.Advertise(context.Background(), serviceName)
		if err != nil {
			log.Printf("Failed to advertise service %s: %v\n", serviceName, err)
		}
	}()
}

func (c *Client) FindPeers(serviceName string) (<-chan peer.AddrInfo, error) {
	return c.discovery.FindPeers(context.Background(), serviceName)
}

func (c *Client) Peers() []peer.ID {
	return c.host.Network().Peers()
}

func (c *Client) Bootstrap() error {
	for _, addr := range dht.DefaultBootstrapPeers {
		pi, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			log.Printf("Failed to parse address: %v\n", err)
			continue
		}
		go func() {
			err = c.host.Connect(context.Background(), *pi)
			if err != nil {
				log.Printf("Failed to connect to peer %s: %v\n", addr, err)
			} else {
				log.Printf("Connected to peer %s\n", addr)
			}
		}()
		time.Sleep(time.Second)
	}

	err := c.dhtNode.Bootstrap(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetConfig() *config.HostConfig {
	return c.config
}

func (c *Client) ID() string {
	return c.host.ID().String()
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
	err = c.ds.Close()
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

	ps, err := pubsub.NewGossipSub(context.Background(), h)
	if err != nil {
		return nil, err
	}

	ds, err := NewBadgerDataStore(cfg.DataStorePath, nil)
	if err != nil {
		return nil, err
	}

	client := &Client{
		host:       h,
		dhtNode:    dhtNode,
		discovery:  discovery,
		pubsub:     ps,
		ds:         ds,
		config:     cfg,
		configPath: configPath,
	}

	log.Println("Bootstraping...")
	err = client.Bootstrap()
	if err != nil {
		return nil, err
	}

	return client, nil
}
