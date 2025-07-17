# Ollama-Nova Phase 1 MVP

A decentralized P2P LLM inference platform extending Ollama with enterprise-grade features.

## 🚀 Features

- **Decentralized P2P Architecture** - libp2p-based networking
- **Ollama-Compatible API** - Drop-in replacement for Ollama
- **Enterprise Security** - TLS encryption and certificate management
- **Comprehensive Monitoring** - Prometheus metrics and health checks
- **Production-Ready** - Docker deployment with security best practices

## 📦 Installation

### Prerequisites
- Go 1.21+
- Docker (optional)
- Ollama running locally

### Quick Start
```bash
git clone https://github.com/khryptorgraphics/ollama-nova.git
cd ollama-nova
go mod tidy
go run cmd/novacron/main.go
```

### Docker
```bash
docker build -t ollama-nova .
docker run -p 8080:8080 -p 9090:9090 ollama-nova
```

## 🔧 Configuration

Edit `configs/prod.yaml` to customize:
- P2P networking settings
- Inference parameters
- Security configuration
- Monitoring endpoints

## 📊 Monitoring

Access metrics at:
- **Prometheus**: http://localhost:9090/metrics
- **Health Check**: http://localhost:9090/health
- **Readiness**: http://localhost:9090/ready

## 🔒 Security

- TLS encryption support
- Certificate management
- Peer validation
- Secure defaults

## 🧪 Testing

```bash
go test ./...
```

## 📈 Roadmap

- **Phase 2**: Advanced P2P networking
- **Phase 3**: Federated learning
- **Phase 4**: Enterprise features
- **Phase 5**: Global edge computing

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
