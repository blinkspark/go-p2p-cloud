package key

import (
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
)

func ReadKey(keyPath string) (crypto.PrivKey, error) {
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPrivateKey(data)
}

func WriteKey(keyPath string, priv crypto.PrivKey) error {
	data, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return err
	}
	return os.WriteFile(keyPath, data, 0666)
}
