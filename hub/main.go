package main

import (
	"log"
	"os"

	"github.com/akamensky/argparse"
	p2pnode "github.com/blinkspark/go-p2p-cloud/p2p-node"
)

func main() {
	var err error
	parser := argparse.NewParser("service-hub", "A p2p-cloud service hub")

	genkeycmd := parser.NewCommand("genkey", "Generate a new private key")
	genkeycmd_key := genkeycmd.String("k", "key", &argparse.Options{Help: "Key path", Default: "key.bin"})

	startcmd := parser.NewCommand("start", "Start a p2p-cloud service hub")
	startcmd_port := startcmd.Int("p", "port", &argparse.Options{Help: "Port", Default: 32233})
	startcmd_key := startcmd.String("k", "key", &argparse.Options{Help: "Key path", Default: "key.bin"})
	err = parser.Parse(os.Args)
	if err != nil {
		log.Panic(err)
	}

	if genkeycmd.Happened() {
		genkey(*genkeycmd_key)
	} else if startcmd.Happened() {
		start(*startcmd_key, uint16(*startcmd_port))
	}
}

func genkey(keyPath string) {
	priv, err := p2pnode.GenPrivKey(keyPath)
	if err != nil {
		log.Panic(err)
	}
	log.Println(priv)
}

func start(keyPath string, port uint16) {
	priv, err := p2pnode.LoadPrivKey(keyPath)
	if err != nil {
		log.Panic(err)
	}

	node, err := p2pnode.NewP2PNode(priv, port)
	if err != nil {
		log.Panic(err)
	}

	log.Println(node.MyAddrs())
}
