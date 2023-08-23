package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"time"
)

func CreateNode(port string) (host.Host, error) {
	//peerChan := make(chan peer.AddrInfo)
	node, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", port)),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
		libp2p.EnableHolePunching(),
	)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Node id: %s\n", node.ID().String())
	fmt.Println("Connect on: ")
	for _, addr := range node.Addrs() {
		fmt.Printf("  %s/p2p/%s\n", addr, node.ID().String())
	}
	fmt.Println(node.Peerstore())
	fmt.Println(node.Addrs())
	return node, nil
}

func main() {
	var (
		discoveryPeers addrList
		port           string
		mode           string
		rendezvous     string
	)
	flag.Var(&discoveryPeers, "peer", "Peer multi address for peer discovery")
	flag.StringVar(&port, "port", "", "Port of user")
	flag.StringVar(&mode, "mode", "", "server/client mode")
	flag.StringVar(&rendezvous, "rendezvous", "", "rendezvous point")
	flag.Parse()

	node, err := CreateNode(port)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(600 * time.Second)
		fmt.Println("Cancelling context now!")
		cancel()
	}()
	modeOpt := dht.ModeClient
	if mode == "server" {
		fmt.Println("Running this node in server mode!")
		modeOpt = dht.ModeServer
	}
	routingDiscovery, err := Announce(ctx, modeOpt, node, rendezvous, discoveryPeers)
	if err != nil {
		panic(err)
	}
	Discover(ctx, node, rendezvous, routingDiscovery)
	select {
	case <-ctx.Done():
		fmt.Println("Context cancelled!")
		break
	}
}
