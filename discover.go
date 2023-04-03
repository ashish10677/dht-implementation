package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	discoveryUtil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"log"
	"time"
)

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
