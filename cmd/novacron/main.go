package main

import (
"context"
"log"
"os"
"os/signal"
"syscall"

"github.com/khryptorgraphics/ollama-nova/api"
"github.com/khryptorgraphics/ollama-nova/internal/config"
"github.com/khryptorgraphics/ollama-nova/internal/inference"
"github.com/khryptorgraphics/ollama-nova/internal/monitoring"
"github.com/khryptorgraphics/ollama-nova/internal/p2p"
"github.com/khryptorgraphics/ollama-nova/internal/security"
)

func main() {
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Load configuration
cfg, err := config.LoadConfig("configs/prod.yaml")
if err != nil {
log.Fatal("Failed to load config:", err)
}

// Initialize components
monitor := monitoring.NewMonitor()
security := security.NewManager()
engine := inference.NewEngine()

// Start P2P node
p2pNode, dht, err := p2p.NewP2PNode(ctx, cfg.P2P.Bootstrap)
if err != nil {
log.Fatal("P2P initialization failed:", err)
}
defer p2pNode.Close()
defer dht.Close()

// Start API server
server := api.NewServer(engine)
go func() {
if err := server.Start(":8080"); err != nil {
log.Fatal("Server failed:", err)
}
}()

// Start monitoring
go monitor.StartMetricsServer(cfg.Monitoring.MetricsPort)

log.Println("Phase 1 MVP started on :8080")

// Handle shutdown
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
<-sigChan
log.Println("Shutting down Phase 1 MVP...")
}
