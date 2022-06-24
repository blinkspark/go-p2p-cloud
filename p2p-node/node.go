package p2pnode

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
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

// BootstrapDefaultDHT bootstraps the DHT with the default bootstrap nodes.
// func (n *P2PNode) BootstrapDefaultDHT(h host.Host) error {
// 	ctx := context.Background()
// 	dht.BootstrapPeers()
// }
func NewP2PNode(privKey crypto.PrivKey, port uint16) (*P2PNode, error) {
	var dhtNode *dht.IpfsDHT

	h, err := libp2p.New(
		libp2p.Identity(privKey),
		libp2p.ListenAddrStrings(listenAddrStrings(port)...),
		libp2p.EnableNATService(),
		libp2p.EnableAutoRelay(), libp2p.EnableHolePunching(), libp2p.NATPortMap(),
		libp2p.Routing(func(h host.Host) (crouting.PeerRouting, error) {
			var err error
			dhtPIs, err := peer.AddrInfosFromP2pAddrs(dht.DefaultBootstrapPeers...)
			if err != nil {
				return nil, err
			}
			dhtNode, err = dht.New(context.Background(), h, dht.BootstrapPeers(dhtPIs...))
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
	strs = append(strs, fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", port))
	strs = append(strs, fmt.Sprintf("/ip6/::/udp/%d/quic", port))
	return strs
}

func LoadPrivKey(keyPath string) (crypto.PrivKey, error) {
	data, err := os.ReadFile(keyPath)
	log.Println(keyPath, data, err)
	if err != nil {
		return nil, err
	}
	buf, err := base64.URLEncoding.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPrivateKey(buf)
}

func LoadPrivKeyFromString(key string) (crypto.PrivKey, error) {
	data, err := base64.URLEncoding.DecodeString(key)
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

	buf := make([]byte, base64.URLEncoding.EncodedLen(len(data)))
	base64.URLEncoding.Encode(buf, data)

	log.Println("Generated private key: ", string(buf))

	err = os.WriteFile(keyPath, buf, 0644)
	return priv, err
}
