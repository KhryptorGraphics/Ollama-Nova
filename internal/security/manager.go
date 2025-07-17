package security

import (
"crypto/rand"
"crypto/rsa"
"crypto/tls"
"crypto/x509"
"crypto/x509/pkix"
"encoding/pem"
"fmt"
"math/big"
"os"
"time"
)

type Manager struct {
certPool    *x509.CertPool
privateKey  *rsa.PrivateKey
certificate *x509.Certificate
config      *Config
}

type Config struct {
TLSEnabled bool   `yaml:"tls_enabled"`
CertPath   string `yaml:"cert_path"`
KeyPath    string `yaml:"key_path"`
CAPath     string `yaml:"ca_path"`
}

type PeerInfo struct {
ID        string
Address   string
PublicKey []byte
Verified  bool
}

func NewManager() *Manager {
return &Manager{
certPool: x509.NewCertPool(),
config: &Config{
TLSEnabled: false,
CertPath:   "/certs/server.crt",
KeyPath:    "/certs/server.key",
CAPath:     "/certs/ca.crt",
},
}
}

func (m *Manager) GenerateSelfSignedCert() error {
privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
if err != nil {
return fmt.Errorf("failed to generate private key: %w", err)
}

template := x509.Certificate{
SerialNumber: big.NewInt(1),
Subject: pkix.Name{
Organization: []string{"Ollama-Nova"},
CommonName:   "ollama-nova.local",
},
NotBefore:             time.Now(),
NotAfter:              time.Now().AddDate(1, 0, 0),
KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
BasicConstraintsValid: true,
}

certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
if err != nil {
return fmt.Errorf("failed to create certificate: %w", err)
}

// Ensure directory exists
os.MkdirAll("/certs", 0755)

// Save certificate
certFile, err := os.Create(m.config.CertPath)
if err != nil {
return fmt.Errorf("failed to create cert file: %w", err)
}
defer certFile.Close()

certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
if _, err := certFile.Write(certPEM); err != nil {
return fmt.Errorf("failed to write certificate: %w", err)
}

// Save private key
keyFile, err := os.Create(m.config.KeyPath)
if err != nil {
return fmt.Errorf("failed to create key file: %w", err)
}
defer keyFile.Close()

keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
if _, err := keyFile.Write(keyPEM); err != nil {
return fmt.Errorf("failed to write private key: %w", err)
}

m.privateKey = privateKey
return nil
}

func (m *Manager) CreateTLSConfig() (*tls.Config, error) {
if !m.config.TLSEnabled {
return nil, nil
}

cert, err := tls.LoadX509KeyPair(m.config.CertPath, m.config.KeyPath)
if err != nil {
return nil, fmt.Errorf("failed to load certificate: %w", err)
}

return &tls.Config{
Certificates: []tls.Certificate{cert},
MinVersion:   tls.VersionTLS12,
NextProtos:   []string{"h2", "http/1.1"},
}, nil
}

func (m *Manager) ValidatePeer(peerID string, cert *x509.Certificate) bool {
// Implement peer certificate validation
opts := x509.VerifyOptions{
Roots: m.certPool,
}

if _, err := cert.Verify(opts); err != nil {
return false
}

return true
}

func (m *Manager) LoadCA(caPath string) error {
caData, err := os.ReadFile(caPath)
if err != nil {
return fmt.Errorf("failed to read CA file: %w", err)
}

if !m.certPool.AppendCertsFromPEM(caData) {
return fmt.Errorf("failed to parse CA certificate")
}

return nil
}
