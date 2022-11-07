package server

import (
	"context"
	"crypto/rand"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

type Server struct {
	protocol string
	host.Host
	*dht.IpfsDHT
	rel *relay.Relay
	*routing.RoutingDiscovery
}

const (
	RELAY_PROTOCOL = "/nealfree.ml/relay/v0.1.0"
)

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

func NewServer(keyPath string, port int, protocol string) (s *Server, err error) {
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

	// // github.com/libp2p/go-libp2p/p2p/net/connmgr
	// cm, err := connmgr.NewConnManager(100, 150)
	// if err != nil {
	// 	return nil, err
	// }
	h, err := libp2p.New(libp2p.Identity(priv), libp2p.ListenAddrStrings(makeAddrs(port)...),
		libp2p.EnableNATService(), libp2p.EnableRelayService(), libp2p.ForceReachabilityPublic())
	if err != nil {
		return nil, err
	}
	dhtNode, err = dht.New(context.Background(), h)
	if err != nil {
		return nil, err
	}
	disc := routing.NewRoutingDiscovery(dhtNode)
	rel, err := relay.New(h)
	if err != nil {
		return nil, err
	}
	return &Server{Host: h, IpfsDHT: dhtNode, rel: rel, RoutingDiscovery: disc, protocol: protocol}, nil
}

func (s *Server) Bootstrap() error {
	err := s.IpfsDHT.Bootstrap(context.Background())
	if err != nil {
		return err
	}
	pis, err := peer.AddrInfosFromP2pAddrs(dht.DefaultBootstrapPeers...)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, pi := range pis {
		wg.Add(1)
		go func(pi peer.AddrInfo) {
			err = s.Connect(context.Background(), pi)
			if err != nil {
				log.Println(err)
			}
			wg.Done()
		}(pi)
	}
	wg.Wait()
	go func() {
		for {
			var ttl time.Duration
			ttl, err = s.RoutingDiscovery.Advertise(context.Background(), s.protocol)
			time.Sleep(ttl)
		}
	}()
	return err
}

// func (s *Server) Advertise() error {

// }
