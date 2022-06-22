package main

func main() {
	// var dhtNode *dht.IpfsDHT
	// h, err := libp2p.New(libp2p.EnableAutoRelay(), libp2p.EnableHolePunching(), libp2p.NATPortMap(), libp2p.Routing(func(h host.Host) (crouting.PeerRouting, error) {
	// 	var err error
	// 	dhtNode, err = dht.New(context.Background(), h)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return dhtNode, nil
	// }))
	// if err != nil {
	// 	log.Panic(err)
	// }

	// err = dhtNode.Bootstrap(context.Background())
	// if err != nil {
	// 	log.Panic(err)
	// }

	// log.Printf("%s", h.ID())
	// log.Println(dhtNode)
	// disc := routing.NewRoutingDiscovery(dhtNode)
	// pubsub.NewGossipSub(context.Background(), h)
}
