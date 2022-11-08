package server

import (
	"crypto/rand"
	"strconv"
	"strings"

	"github.com/libp2p/go-libp2p"

	// "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

type Server struct {
	protocol string
	host.Host
	rel *relay.Relay
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

func NewServer(keyPath string, port int) (s *Server, err error) {
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
	return &Server{Host: h}, nil
}

// func (s *Server) Advertise() error {

// }
