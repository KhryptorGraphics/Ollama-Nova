package p2p

import (
"context"
"fmt"

"github.com/libp2p/go-libp2p"
"github.com/libp2p/go-libp2p/core/host"
"github.com/libp2p/go-libp2p/core/peer"
"github.com/libp2p/go-libp2p/p2p/net/swarm"
)

func NewP2PNode(ctx context.Context, bootstrapPeers []string) (host.Host, error) {
host, err := libp2p.New(
libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
libp2p.EnableAutoRelay(),
)
if err != nil {
return nil, fmt.Errorf("failed to create host: %w", err)
}

// Connect to bootstrap peers
for _, addr := range bootstrapPeers {
peerInfo, _ := peer.AddrInfoFromString(addr)
if peerInfo != nil {
host.Connect(ctx, *peerInfo)
}
}

return host, nil
}
