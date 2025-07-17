package p2p

import (
"context"
"fmt"
"time"

"github.com/libp2p/go-libp2p"
"github.com/libp2p/go-libp2p-kad-dht"
"github.com/libp2p/go-libp2p/core/host"
"github.com/libp2p/go-libp2p/core/peer"
"github.com/libp2p/go-libp2p/p2p/net/swarm"
"github.com/multiformats/go-multiaddr"
)

func NewP2PNode(ctx context.Context, bootstrapPeers []string) (host.Host, *dht.IpfsDHT, error) {
host, err := libp2p.New(
libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
libp2p.EnableAutoRelay(),
libp2p.EnableNATService(),
libp2p.EnableHolePunching(),
)
if err != nil {
return nil, nil, fmt.Errorf("failed to create host: %w", err)
}

dht, err := dht.New(ctx, host, dht.Mode(dht.ModeServer))
if err != nil {
return nil, nil, fmt.Errorf("failed to create DHT: %w", err)
}

// Bootstrap DHT
for _, addr := range bootstrapPeers {
peerInfo, _ := peer.AddrInfoFromString(addr)
if peerInfo != nil {
host.Connect(ctx, *peerInfo)
}
}

return host, dht, nil
}
