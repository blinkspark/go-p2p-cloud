package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/akamensky/argparse"
	p2pnode "github.com/blinkspark/go-p2p-cloud/p2p-node"
	"github.com/libp2p/go-libp2p-core/crypto"
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

	servercmd := parser.NewCommand("start", "Start a p2p-cloud service hub")
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
		// log.Println(*startcmd_key, *startcmd_keyfile, *startcmd_port)
		start(*startcmd_key, *startcmd_keyfile, uint16(*startcmd_port))
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

func start(key string, keyPath string, port uint16) {
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

	log.Println(node.MyAddrs())
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

	log.Println(node.MyAddrs())
	<-sigChan
}
