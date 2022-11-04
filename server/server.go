package server

import (
	"crypto/rand"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
)

type Server struct {
	host.Host
}

func NewServer(keyPath string) (s *Server, err error) {
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

	h, err := libp2p.New(libp2p.Identity(priv))
	if err != nil {
		return nil, err
	}
	return &Server{Host: h}, nil
}

func readKey(keyPath string) (crypto.PrivKey, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPrivateKey(data)
}

func writeKey(keyPath string, priv crypto.PrivKey) error {
	data, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return err
	}
	return os.WriteFile(keyPath, data, 0666)
}
