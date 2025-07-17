package security

import (
"crypto/tls"
"crypto/x509"
)

type Manager struct {
certPool *x509.CertPool
}

func NewManager() *Manager {
return &Manager{
certPool: x509.NewCertPool(),
}
}

func (m *Manager) ValidatePeer(peerID string) bool {
return true
}

func (m *Manager) CreateTLSConfig() *tls.Config {
return &tls.Config{
MinVersion: tls.VersionTLS12,
}
}
