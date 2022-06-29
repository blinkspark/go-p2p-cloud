package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/akamensky/argparse"
	p2pnode "github.com/blinkspark/go-p2p-cloud/p2p-node"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

func main() {
	var err error
	parser := argparse.NewParser("service-hub", "A p2p-cloud service hub")

	genkeycmd := parser.NewCommand("genkey", "Generate a new private key")
	genkeycmd_key := genkeycmd.String("f", "fpath", &argparse.Options{Help: "Key path", Default: "key.txt"})

	startcmd := parser.NewCommand("start", "Start a p2p-cloud service hub")
	startcmd_port := startcmd.Int("p", "port", &argparse.Options{Help: "Port", Default: 32233})
	startcmd_key := startcmd.String("k", "key", &argparse.Options{Help: "Key file path"})
	startcmd_keyfile := startcmd.String("f", "keyfile", &argparse.Options{Help: "Key file path", Default: "key.txt"})
	startcmd_boot := startcmd.StringList("b", "bootstrap", &argparse.Options{Help: "Bootstrap nodes"})

	err = parser.Parse(os.Args)
	if err != nil {
		log.Panic(err)
	}

	if genkeycmd.Happened() {
		genkey(*genkeycmd_key)
	} else if startcmd.Happened() {
		start(*startcmd_key, *startcmd_keyfile, *startcmd_boot, uint16(*startcmd_port))
	}
}

func genkey(keyPath string) {
	_, err := p2pnode.GenPrivKey(keyPath)
	if err != nil {
		log.Panic(err)
	}
}

func start(key string, keyPath string, bootstrap []string, port uint16) {
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
	log.Println("my id:", node.ID())

	wg := sync.WaitGroup{}
	for _, b := range bootstrap {
		wg.Add(1)
		go func(b string) {
			defer wg.Done()
			log.Println("connecting to bootstrap:", b)
			pi, err := peer.AddrInfoFromString(b)
			if err != nil {
				log.Println(err)
				return
			}
			err = node.Connect(context.Background(), *pi)
			if err != nil {
				log.Println(err)
			}
		}(b)
	}
	wg.Wait()

	node.AdvertiseService("nealfree.ml/test/v0.1.1")
	node.SetStreamHandler("nealfree.ml/test/v0.1.1", func(s network.Stream) {
		log.Println("received a stream ssssssssssssssssssssssssssssssssssssssssss")
		s.Close()
	})

	go func() {
		for {
			pic, err := node.FindPeers(context.Background(), "nealfree.ml/test/v0.1.1")
			log.Println("found:")
			if err != nil {
				log.Println(err)
			}
			for p := range pic {
				log.Println(p)
				err = node.Connect(context.Background(), p)
				if err != nil {
					log.Println(err)
				}

				if p.ID == node.ID() || len(p.Addrs) == 0 {
					log.Println("skip self or empty addr")
					continue
				}
				s, err := node.NewStream(context.Background(), p.ID, "nealfree.ml/test/v0.1.1")
				if err != nil {
					log.Println(s, err)
				}
				// s.Close()
			}
			time.Sleep(time.Second * 5)
		}
	}()

	// go node.TestShowPeerCount()
	// go node.TestShowConnectionCount()
	// go node.TestPings()

	<-sigChan
	node.CloseAll()
}
