# Kademlia DHT Implementation 

## Objectives
- Create a libp2p node.
- Advertise your rendezvous
- Find and connect to `n` peers. 
- Wait and do not proceed till connected to all `n` peers.

## How to run
Build the project using 
```
$ go build -o dht .
```

For the first node, run 
```
$ ./dht
```
It'll give you an output like this:
```
Node id: 12D3KooWCEsQ9vbkBCgR7gwF1xcqZiNrNkhA57EjUuho48MiypaT
Connect on: 
/ip4/192.168.1.20/udp/53634/quic/p2p/12D3KooWCEsQ9vbkBCgR7gwF1xcqZiNrNkhA57EjUuho48MiypaT
/ip4/127.0.0.1/udp/53634/quic/p2p/12D3KooWCEsQ9vbkBCgR7gwF1xcqZiNrNkhA57EjUuho48MiypaT
/ip4/192.168.1.20/udp/59552/quic-v1/webtransport/certhash/uEiCKymWCAadpwgqyAUa4gj8CXCrFK_eMZuHa4p0zXrKm4A/certhash/uEiA8fIi0sbwy-6osWTT-sB-lzm82iggETvWfVox-L7iwqw/p2p/12D3KooWCEsQ9vbkBCgR7gwF1xcqZiNrNkhA57EjUuho48MiypaT
/ip4/127.0.0.1/udp/59552/quic-v1/webtransport/certhash/uEiCKymWCAadpwgqyAUa4gj8CXCrFK_eMZuHa4p0zXrKm4A/certhash/uEiA8fIi0sbwy-6osWTT-sB-lzm82iggETvWfVox-L7iwqw/p2p/12D3KooWCEsQ9vbkBCgR7gwF1xcqZiNrNkhA57EjUuho48MiypaT
/ip6/::1/udp/53018/quic/p2p/12D3KooWCEsQ9vbkBCgR7gwF1xcqZiNrNkhA57EjUuho48MiypaT
/ip6/::1/udp/53018/quic-v1/p2p/12D3KooWCEsQ9vbkBCgR7gwF1xcqZiNrNkhA57EjUuho48MiypaT
/ip6/::1/udp/54603/quic-v1/webtransport/certhash/uEiCKymWCAadpwgqyAUa4gj8CXCrFK_eMZuHa4p0zXrKm4A/certhash/uEiA8fIi0sbwy-6osWTT-sB-lzm82iggETvWfVox-L7iwqw/p2p/12D3KooWCEsQ9vbkBCgR7gwF1xcqZiNrNkhA57EjUuho48MiypaT
/ip4/192.168.1.20/tcp/58476/p2p/12D3KooWCEsQ9vbkBCgR7gwF1xcqZiNrNkhA57EjUuho48MiypaT
/ip4/127.0.0.1/tcp/58476/p2p/12D3KooWCEsQ9vbkBCgR7gwF1xcqZiNrNkhA57EjUuho48MiypaT
/ip4/192.168.1.20/udp/53634/quic-v1/p2p/12D3KooWCEsQ9vbkBCgR7gwF1xcqZiNrNkhA57EjUuho48MiypaT
/ip4/127.0.0.1/udp/53634/quic-v1/p2p/12D3KooWCEsQ9vbkBCgR7gwF1xcqZiNrNkhA57EjUuho48MiypaT
/ip6/::1/tcp/58479/p2p/12D3KooWCEsQ9vbkBCgR7gwF1xcqZiNrNkhA57EjUuho48MiypaT
```

In the next 3 new terminals run

```
$ ./dht --peer /ip4/192.168.1.20/tcp/58476/p2p/12D3KooWCEsQ9vbkBCgR7gwF1xcqZiNrNkhA57EjUuho48MiypaT --port <port> --externalIp <your_public_ipAddress>
```

You'll see all the nodes are connected to each other.