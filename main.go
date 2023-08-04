package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
	"time"
)

func CreateNode(externalIp string, port string) (host.Host, error) {
	addressFactory, err := getAddressFactory(externalIp, port)
	if err != nil {
		return nil, err
	}
	node, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("%s%s", "/ip4/0.0.0.0/tcp/", port)),
		libp2p.AddrsFactory(addressFactory))
	if err != nil {
		return nil, err
	}
	fmt.Printf("Node id: %s\n", node.ID().String())
	fmt.Println("Connect on: ")
	for _, addr := range node.Addrs() {
		fmt.Printf("  %s/p2p/%s\n", addr, node.ID().String())
	}
	return node, nil
}

func getAddressFactory(externalIP string, port string) (func([]multiaddr.Multiaddr) []multiaddr.Multiaddr, error) {
	var (
		externalAddr multiaddr.Multiaddr
		err          error
	)

	// binding external IP to libp2p node
	if len(externalIP) != 0 {
		externalAddr, err = multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", externalIP, port))
		if err != nil {
			return nil, fmt.Errorf("fail to create listen with given external IP: %v", err)
		}
		fmt.Printf("Binding to external IP: %v", externalAddr.String())
	}
	addressFactory := func(addrs []multiaddr.Multiaddr) []multiaddr.Multiaddr {
		if externalAddr != nil {
			return []multiaddr.Multiaddr{externalAddr}
		}
		return addrs
	}
	return addressFactory, nil
}

func main() {
	var discoveryPeers addrList
	var externalIp string
	var port string
	flag.Var(&discoveryPeers, "peer", "Peer multi address for peer discovery")
	//flag.StringVar(&externalIp, "externalIp", "", "Public IP address of user")
	flag.StringVar(&port, "port", "", "Port of user")
	flag.Parse()

	node, err := CreateNode(externalIp, port)
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
