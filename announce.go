package main

import (
	"context"
	"fmt"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/discovery"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	discoveryUtil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
	"sync"
	"time"
)

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
