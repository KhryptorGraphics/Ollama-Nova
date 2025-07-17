package config

type Config struct {
P2P      P2PConfig      `yaml:"p2p"`
Inference InferenceConfig `yaml:"inference"`
Security SecurityConfig  `yaml:"security"`
Monitoring MonitoringConfig `yaml:"monitoring"`
}

type P2PConfig struct {
Port        int      `yaml:"port"`
Bootstrap   []string `yaml:"bootstrap"`
MaxPeers    int      `yaml:"max_peers"`
}

type InferenceConfig struct {
ModelPath   string `yaml:"model_path"`
MaxTokens   int    `yaml:"max_tokens"`
Temperature float64 `yaml:"temperature"`
}

type SecurityConfig struct {
TLS         bool   `yaml:"tls"`
CertPath    string `yaml:"cert_path"`
KeyPath     string `yaml:"key_path"`
}

type MonitoringConfig struct {
MetricsPort int    `yaml:"metrics_port"`
LogLevel    string `yaml:"log_level"`
}

func LoadConfig(path string) (*Config, error) {
// Implementation
return &Config{}, nil
}
