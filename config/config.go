package config

import (
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
)

type HostConfig struct {
	PrivateKey    []byte `json:"private_key"`
	Port          int    `json:"port"`
	DataStorePath string `json:"data_store_path"`
}

func DefaultConfig() (*HostConfig, error) {
	privKey, _, err := crypto.GenerateKeyPair(crypto.Ed25519, 0)
	if err != nil {
		return nil, err
	}
	privKeyBytes, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	return &HostConfig{
		PrivateKey:    privKeyBytes,
		Port:          12233,
		DataStorePath: "./data",
	}, nil
}

func LoadConfig(path string) (*HostConfig, error) {
	config := &HostConfig{}

	configFile, err := os.Open(path)
	var defaultConfig *HostConfig
	if err != nil {
		if os.IsNotExist(err) {
			defaultConfig, err = DefaultConfig()
			if err != nil {
				return nil, err
			}
			configFile, err = os.Create(path)
			if err != nil {
				return nil, err
			}
			defer configFile.Close()
			var configBytes []byte
			configBytes, err = json.MarshalIndent(defaultConfig, "", "  ")
			if err != nil {
				return nil, err
			}
			_, err = configFile.Write(configBytes)
			if err != nil {
				return nil, err
			}
			return defaultConfig, nil
		}
		return nil, err
	}
	defer configFile.Close()

	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(configBytes, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func BuildOptionFromConfig(cfg *HostConfig) ([]libp2p.Option, error) {
	privKey, err := crypto.UnmarshalPrivateKey(cfg.PrivateKey)
	if err != nil {
		return nil, err
	}

	options := []libp2p.Option{
		libp2p.Identity(privKey),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/"+strconv.Itoa(cfg.Port),
			"/ip4/0.0.0.0/udp/"+strconv.Itoa(cfg.Port)+"/quic-v1",
			"/ip4/0.0.0.0/udp/"+strconv.Itoa(cfg.Port)+"/quic-v1/webtransport",
			"/ip4/0.0.0.0/udp/"+strconv.Itoa(cfg.Port)+"/webrtc-direct",
			"/ip6/::/tcp/"+strconv.Itoa(cfg.Port),
			"/ip6/::/udp/"+strconv.Itoa(cfg.Port)+"/quic-v1",
			"/ip6/::/udp/"+strconv.Itoa(cfg.Port)+"/quic-v1/webtransport",
			"/ip6/::/udp/"+strconv.Itoa(cfg.Port)+"/webrtc-direct",
		),
		libp2p.EnableAutoNATv2(),
		libp2p.EnableHolePunching(),
		libp2p.NATPortMap(),
	}
	return options, nil
}
