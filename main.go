package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
	"time"
)

func CreateNode(port string) (host.Host, error) {
	peerChan := make(chan peer.AddrInfo)
	node, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", port)),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
		libp2p.EnableAutoRelayWithPeerSource(
			func(ctx context.Context, numPeers int) <-chan peer.AddrInfo {
				r := make(chan peer.AddrInfo)
				go func() {
					defer close(r)
					for ; numPeers != 0; numPeers-- {
						select {
						case v, ok := <-peerChan:
							if !ok {
								return
							}
							select {
							case r <- v:
							case <-ctx.Done():
								return
							}
						case <-ctx.Done():
							return
						}
					}
				}()
				return r
			},
			autorelay.WithMinInterval(0),
		),
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
	routingDiscovery, err := Announce(ctx, modeOpt, node, rendezvous, dht.DefaultBootstrapPeers)
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
