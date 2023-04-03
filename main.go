package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/discovery"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	discoveryUtil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
	"log"
	"sync"
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

func Announce(ctx context.Context, node host.Host, bootstrapPeers []multiaddr.Multiaddr) (*routing.RoutingDiscovery, error) {
	kademliaDHT, err := dht.New(ctx, node, dht.Mode(dht.ModeServer))
	if err != nil {
		return nil, fmt.Errorf("fail to create DHT: %w", err)
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		return nil, fmt.Errorf("fail to bootstrap DHT: %w", err)
	}
	var wg sync.WaitGroup
	for _, peerAddr := range bootstrapPeers {
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := node.Connect(ctx, *peerInfo); err != nil {
				fmt.Printf("Error while connecting to node %q: %-v\n", peerInfo, err)
			} else {
				fmt.Printf("Connection established with bootstrap node: %q\n", *peerInfo)
			}
		}()
	}
	wg.Wait()
	routingDiscovery := routing.NewRoutingDiscovery(kademliaDHT)
	fmt.Println("Announcing rendezvous as: ", Rendezvous)
	discoveryUtil.Advertise(ctx, routingDiscovery, Rendezvous, discovery.TTL(5*time.Second))
	return routingDiscovery, nil
}

func Discover(ctx context.Context, node host.Host, routingDiscovery *routing.RoutingDiscovery) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			peers, err := discoveryUtil.FindPeers(ctx, routingDiscovery, Rendezvous)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Found %d peers\n\n", len(peers))
			for i, p := range peers {
				if p.ID == node.ID() {
					continue
				}
				fmt.Printf("Peer %d: %s\n", i+1, p.ID)
				fmt.Printf("Addresses: %s\n\n\n\n", p.Addrs)
				if node.Network().Connectedness(p.ID) != network.Connected {
					_, err = node.Network().DialPeer(ctx, p.ID)
					if err != nil {
						fmt.Println("Error:", err)
						continue
					}
				}
			}
			fmt.Printf("Connected to %d peers\n", len(node.Network().Peers()))

			if len(node.Network().Peers()) == TotalNumberOfPeers-1 {
				fmt.Println("Connected to all peers")
				return
			}
		}
	}
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
