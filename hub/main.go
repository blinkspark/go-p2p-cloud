package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/akamensky/argparse"
	p2pnode "github.com/blinkspark/go-p2p-cloud/p2p-node"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

func main() {
	var err error
	parser := argparse.NewParser("service-hub", "A p2p-cloud service hub")

	genkeycmd := parser.NewCommand("genkey", "Generate a new private key")
	genkeycmd_key := genkeycmd.String("k", "key", &argparse.Options{Help: "Key path", Default: "key.txt"})

	startcmd := parser.NewCommand("start", "Start a p2p-cloud service hub")
	startcmd_port := startcmd.Int("p", "port", &argparse.Options{Help: "Port", Default: 32233})
	startcmd_key := startcmd.String("k", "key", &argparse.Options{Help: "Key file path"})
	startcmd_keyfile := startcmd.String("f", "keyfile", &argparse.Options{Help: "Key file path", Default: "key.txt"})
	startcmd_dial := startcmd.String("d", "dial", &argparse.Options{Help: "Dial address"})

	servercmd := parser.NewCommand("server", "Start a p2p-cloud service hub")
	servercmd_port := servercmd.Int("p", "port", &argparse.Options{Help: "Port", Default: 32233})
	servercmd_key := servercmd.String("k", "key", &argparse.Options{Help: "Key file path"})
	servercmd_keyfile := servercmd.String("f", "keyfile", &argparse.Options{Help: "Key file path", Default: "key.txt"})

	err = parser.Parse(os.Args)
	if err != nil {
		log.Panic(err)
	}

	if genkeycmd.Happened() {
		genkey(*genkeycmd_key)
	} else if startcmd.Happened() {
		start(*startcmd_key, *startcmd_keyfile, *startcmd_dial, uint16(*startcmd_port))
	} else if servercmd.Happened() {
		server(*servercmd_key, *servercmd_keyfile, uint16(*servercmd_port))
	}
}

func genkey(keyPath string) {
	_, err := p2pnode.GenPrivKey(keyPath)
	if err != nil {
		log.Panic(err)
	}
}

func start(key string, keyPath string, dial string, port uint16) {
	var (
		priv    crypto.PrivKey
		err     error
		sigChan chan os.Signal = make(chan os.Signal, 1)
	)
	signal.Notify(sigChan, os.Interrupt)

	if key != "" {
		priv, err = p2pnode.LoadPrivKeyFromString(key)
		if err != nil {
			log.Panic(err)
		}
	} else {
		log.Println("Loading key from", keyPath)
		priv, err = p2pnode.LoadPrivKey(keyPath)
		if err != nil {
			log.Panic(err)
		}
	}

	node, err := p2pnode.NewP2PNode(priv, port)
	if err != nil {
		log.Panic(err)
	}

	pi, err := peer.AddrInfoFromString(dial)
	if err != nil {
		log.Panic(err)
	}
	err = node.Connect(context.Background(), *pi)
	if err != nil {
		log.Panic(err)
	}
	log.Println("connected to", dial)

	log.Println(node.MyAddrs())

	topic, err := node.Join("nealfree.ml/p2p-cloud/service-hub/pubusb/v0.1.0")
	if err != nil {
		log.Panic(err)
	}
	log.Println("joined topic", topic)

	log.Println("pubsub peers:", topic.ListPeers())
	sub, err := topic.Subscribe()
	if err != nil {
		log.Panic(err)
	}
	go func() {
		for {
			msg, err := sub.Next(context.Background())
			if err != nil {
				log.Panic(err)
			}
			log.Println("from:", msg.GetFrom(), "got message:", string(msg.GetData()))
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second)
			err = topic.Publish(context.Background(), []byte("hello"))
			if err != nil {
				log.Println(err)
			}
		}
	}()

	<-sigChan
}

func server(key string, keyPath string, port uint16) {
	var (
		priv    crypto.PrivKey
		err     error
		sigChan chan os.Signal = make(chan os.Signal, 1)
	)
	signal.Notify(sigChan, os.Interrupt)

	if key != "" {
		priv, err = p2pnode.LoadPrivKeyFromString(key)
		if err != nil {
			log.Panic(err)
		}
	} else {
		priv, err = p2pnode.LoadPrivKey(keyPath)
		if err != nil {
			log.Panic(err)
		}
	}

	node, err := p2pnode.NewP2PNode(priv, port)
	if err != nil {
		log.Panic(err)
	}

	topic, err := node.Join("nealfree.ml/p2p-cloud/service-hub/pubusb/v0.1.0")
	if err != nil {
		log.Panic(err)
	}
	log.Println("joined topic", topic)
	log.Println("my addrs:", node.MyAddrs())

loop:
	for {
		select {
		case <-sigChan:
			break loop
		default:
			time.Sleep(time.Second * 15)
			peers := node.Peerstore().Peers()
			for _, peer := range peers {
				log.Println("peer:", peer)
				log.Println(node.Peerstore().Addrs(peer))
			}
		}
	}

}
