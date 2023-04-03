package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"time"
)

func CreateNode() (host.Host, error) {
	node, err := libp2p.New()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Node id: %s\n", node.ID().String())
	fmt.Println("Connect on: ")
	for _, addr := range node.Addrs() {
		fmt.Printf("  %s/p2p/%s", addr, node.ID().String())
	}
	return node, nil
}

func main() {
	var discoveryPeers addrList
	flag.Var(&discoveryPeers, "peer", "Peer multi address for peer discovery")
	flag.Parse()

	node, err := CreateNode()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(120 * time.Second)
		fmt.Println("Cancelling context now!")
		cancel()
	}()
	routingDiscovery, err := Announce(ctx, node, discoveryPeers)
	if err != nil {
		panic(err)
	}
	Discover(ctx, node, routingDiscovery)
	select {
	case <-ctx.Done():
		fmt.Println("Context cancelled!")
		break
	}
}
