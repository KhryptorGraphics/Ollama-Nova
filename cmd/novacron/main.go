package main

import (
"context"
"flag"
"log"
"os"
"os/signal"
"syscall"

"github.com/khryptorgraphics/ollama-nova/internal/p2p"
"github.com/khryptorgraphics/ollama-nova/internal/inference"
"github.com/khryptorgraphics/ollama-nova/internal/security"
"github.com/khryptorgraphics/ollama-nova/internal/monitoring"
)

func main() {
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Initialize components
monitor := monitoring.NewMonitor()
securityManager := security.NewManager()

// Start P2P node
p2pNode, _, err := p2p.NewP2PNode(ctx, []string{})
if err != nil {
log.Fatal("P2P initialization failed: ", err)
}
defer p2pNode.Close()

// Start inference engine
engine := inference.NewEngine()

// Handle shutdown
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
<-sigChan
log.Println("Shutting down...")
}
