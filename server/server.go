package server

import (
	"context"
	"crypto/rand"
	"strconv"
	"strings"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
)

type Server struct {
	host.Host
	*dht.IpfsDHT
	*pubsub.PubSub
}

func makeAddrs(port int) []string {
	temps := []string{
		"/ip4/0.0.0.0/tcp/{PORT}",
		"/ip4/0.0.0.0/udp/{PORT}/quic",
		"/ip6/::/tcp/{PORT}",
		"/ip6/::/udp/{PORT}/quic",
	}
	var addrs []string
	for _, addr := range temps {
		addr = strings.Replace(addr, "{PORT}", strconv.Itoa(port), 1)
		addrs = append(addrs, addr)
	}
	return addrs
}

func NewServer(keyPath string, port int, topic string) (s *Server, err error) {
	priv, err := readKey(keyPath)
	if err != nil {
		priv, _, err = crypto.GenerateEd25519Key(rand.Reader)
		if err != nil {
			return nil, err
		}
		err = writeKey(keyPath, priv)
		if err != nil {
			return nil, err
		}
	}

	var dhtNode *dht.IpfsDHT

	// h, err := libp2p.New(libp2p.Identity(priv), libp2p.NATPortMap(), libp2p.EnableAutoRelay(autorelay.WithDefaultStaticRelays()), libp2p.EnableHolePunching(), libp2p.EnableRelayService(), libp2p.EnableNATService(), libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
	// 	dhtNode, err = dht.New(context.Background(), h)
	// 	return dhtNode, err
	// }))
	h, err := libp2p.New(libp2p.Identity(priv), libp2p.ListenAddrStrings(makeAddrs(port)...), libp2p.NATPortMap(), libp2p.EnableNATService(), libp2p.EnableRelayService(), libp2p.EnableHolePunching(), libp2p.ForceReachabilityPublic())
	if err != nil {
		return nil, err
	}
	ps, err := pubsub.NewGossipSub(context.Background(), h)
	if err != nil {
		return nil, err
	}
	t, err := ps.Join(topic)
	if err != nil {
		return nil, err
	}
	_, err = t.Subscribe()
	if err != nil {
		return nil, err
	}
	return &Server{Host: h, IpfsDHT: dhtNode, PubSub: ps}, nil
}

func (s *Server) Bootstrap() error {
	err := s.IpfsDHT.Bootstrap(context.Background())
	if err != nil {
		return err
	}
	return nil
}
