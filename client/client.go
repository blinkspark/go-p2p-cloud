package client

import (
	"crypto/rand"
	"log"

	"github.com/blinkspark/go-p2p-cloud/key"
	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
)

type Client struct {
	host.Host
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
	h, err := libp2p.New(
		libp2p.Identity(priv), libp2p.EnableHolePunching(), libp2p.NATPortMap(),
		libp2p.EnableAutoRelay(autorelay.WithStaticRelays(staticRelays)),
		libp2p.EnableNATService(), libp2p.EnableRelayService())
	if err != nil {
		return nil, err
	}
	return &Client{Host: h}, nil
}
