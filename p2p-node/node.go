package p2pnode

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	leveldb "github.com/ipfs/go-ds-leveldb"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	crouting "github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-peerstore/pstoreds"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
)

type P2PNode struct {
	host.Host
	*dht.IpfsDHT
	*routing.RoutingDiscovery
	*pubsub.PubSub
	advertiseTicker *time.Ticker
}

func (n *P2PNode) MyAddrs() (addrs []string) {
	for _, addr := range n.Addrs() {
		addrs = append(addrs, fmt.Sprintf("%s/p2p/%s", addr, n.ID().Pretty()))
	}
	return
}

func (n *P2PNode) ShowMyAddrs() {
	addrs := n.MyAddrs()
	for _, addr := range addrs {
		log.Println(addr)
	}
}

func (n *P2PNode) CloseAll() {
	n.Host.Close()
	n.IpfsDHT.Close()
	n.Peerstore().Close()
}

// BootstrapDefaultDHT bootstraps the DHT with the default bootstrap nodes.
// func (n *P2PNode) BootstrapDefaultDHT(h host.Host) error {
// 	ctx := context.Background()
// 	dht.BootstrapPeers()
// }
func NewP2PNode(privKey crypto.PrivKey, port uint16) (*P2PNode, error) {
	var dhtNode *dht.IpfsDHT

	cm, err := connmgr.NewConnManager(50, 100)
	if err != nil {
		return nil, err
	}

	ds, err := leveldb.NewDatastore("peerstore", nil)
	if err != nil {
		return nil, err
	}
	// defer ds.Close()
	pstore, err := pstoreds.NewPeerstore(context.Background(), ds, pstoreds.DefaultOpts())
	if err != nil {
		return nil, err
	}
	// defer pstore.Close()

	h, err := libp2p.New(
		libp2p.Identity(privKey),
		libp2p.Peerstore(pstore),
		libp2p.ListenAddrStrings(listenAddrStrings(port)...),
		libp2p.EnableNATService(),
		libp2p.ConnectionManager(cm),
		libp2p.EnableAutoRelay(), libp2p.EnableHolePunching(), libp2p.NATPortMap(),
		libp2p.Routing(func(h host.Host) (crouting.PeerRouting, error) {
			var err error
			dhtPIs := dht.GetDefaultBootstrapPeerAddrInfos()

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

func (n *P2PNode) TestShowPeerCount() {
	for range time.Tick(time.Second * 10) {
		log.Println("peer count:", len(n.Peerstore().Peers()))
	}
}

func (n *P2PNode) TestShowConnectionCount() {
	for range time.Tick(time.Second * 10) {
		log.Println("connection count:", len(n.Network().Conns()))
	}
}

func (n *P2PNode) TestPings() {
	log.Println("pinging...")
	for range time.Tick(time.Second * 30) {
		for _, conn := range n.Network().Conns() {
			log.Println("pinging:", conn.RemotePeer())
			res := <-ping.Ping(context.Background(), n.Host, conn.RemotePeer())
			log.Println("ping:", conn.RemotePeer(), res.RTT)
		}
	}
}

func (n *P2PNode) AdvertiseService(service string) {
	var (
		ttl   time.Duration
		err   error
		start = make(chan struct{})
	)
	go func() {
		for {
			time.Sleep(time.Second * 5)
			ttl, err = n.RoutingDiscovery.Advertise(context.Background(), service)
			log.Println("adv:", ttl, err)
			if err == nil {
				break
			}
		}
		n.advertiseTicker = time.NewTicker(ttl)
		close(start)
	}()
	go func() {
		<-start
		for range n.advertiseTicker.C {
			ttl, err = n.RoutingDiscovery.Advertise(context.Background(), service)
			log.Println("adv:", ttl, err)
		}
	}()
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
