package p2pnode

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	crouting "github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
)

type P2PNode struct {
	host.Host
	*dht.IpfsDHT
	*routing.RoutingDiscovery
	*pubsub.PubSub
}

func (n *P2PNode) MyAddrs() (addrs []string) {
	for _, addr := range n.Addrs() {
		addrs = append(addrs, fmt.Sprintf("%s/p2p/%s", addr, n.ID().Pretty()))
	}
	return
}

func NewP2PNode(privKey crypto.PrivKey, port uint16) (*P2PNode, error) {
	var dhtNode *dht.IpfsDHT

	h, err := libp2p.New(
		libp2p.Identity(privKey),
		libp2p.ListenAddrStrings(listenAddrStrings(port)...),
		libp2p.EnableAutoRelay(), libp2p.EnableHolePunching(), libp2p.NATPortMap(),
		libp2p.Routing(func(h host.Host) (crouting.PeerRouting, error) {
			var err error
			dhtNode, err = dht.New(context.Background(), h)
			if err != nil {
				return nil, err
			}
			return dhtNode, nil
		}))
	if err != nil {
		return nil, err
	}

	err = dhtNode.Bootstrap(context.Background())
	if err != nil {
		return nil, err
	}

	ps, err := pubsub.NewGossipSub(context.Background(), h)
	if err != nil {
		return nil, err
	}

	return &P2PNode{
		Host:             h,
		IpfsDHT:          dhtNode,
		RoutingDiscovery: routing.NewRoutingDiscovery(dhtNode),
		PubSub:           ps,
	}, nil
}

func listenAddrStrings(port uint16) []string {
	var strs []string
	strs = append(strs, fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	strs = append(strs, fmt.Sprintf("/ip6/::/tcp/%d", port))
	return strs
}

func LoadPrivKey(keyPath string) (crypto.PrivKey, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPrivateKey(data)
}

func GenPrivKey(keyPath string) (crypto.PrivKey, error) {
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, err
	}

	data, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	err = os.WriteFile(keyPath, data, 0644)
	return priv, err
}
