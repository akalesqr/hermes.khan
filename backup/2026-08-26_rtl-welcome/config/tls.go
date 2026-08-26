package config

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// EnsureSelfSignedCert generates a self-signed TLS cert+key into the data dir
// if they don't already exist. Returns the cert and key paths.
// Used when TLS is enabled but no cert/key paths are configured.
func (c *Config) EnsureSelfSignedCert() (certPath, keyPath string, err error) {
	if c.Server.CertPath != "" && c.Server.KeyPath != "" {
		return c.Server.CertPath, c.Server.KeyPath, nil
	}
	certPath = filepath.Join(c.Server.DataDir, "tls-cert.pem")
	keyPath = filepath.Join(c.Server.DataDir, "tls-key.pem")

	// If both already exist, reuse them (stable across restarts).
	if _, e1 := os.Stat(certPath); e1 == nil {
		if _, e2 := os.Stat(keyPath); e2 == nil {
			c.Server.CertPath = certPath
			c.Server.KeyPath = keyPath
			return certPath, keyPath, nil
		}
	}

	if err := os.MkdirAll(c.Server.DataDir, 0o755); err != nil {
		return "", "", err
	}

	// Generate an ECDSA P-256 self-signed cert valid for 10 years.
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", err
	}
	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	now := time.Now()
	tmpl := x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: "Khan LAN Chat", Organization: []string{"Khan"}},
		NotBefore:    now.Add(-time.Hour),
		NotAfter:     now.Add(10 * 365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"localhost", "khan.local", "lan"},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1)},
	}
	der, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	if err != nil {
		return "", "", err
	}

	certOut, err := os.Create(certPath)
	if err != nil {
		return "", "", err
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: der}); err != nil {
		return "", "", err
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return "", "", err
	}
	keyOut, err := os.Create(keyPath)
	if err != nil {
		return "", "", err
	}
	defer keyOut.Close()
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes}); err != nil {
		return "", "", err
	}

	c.Server.CertPath = certPath
	c.Server.KeyPath = keyPath
	return certPath, keyPath, nil
}