package mockserver

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

const (
	tlsCertificateValidity = 365 * 24 * time.Hour
	tlsSerialNumberBits    = 128
)

// newTLSConfig creates a TLS configuration. If certFile and keyFile are empty,
// it generates a self-signed certificate for localhost.
func newTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	var certificate tls.Certificate
	var err error

	if certFile != "" && keyFile != "" {
		certificate, err = tls.LoadX509KeyPair(certFile, keyFile)
	} else {
		certificate, err = generateSelfSignedCertificate()
	}
	if err != nil {
		return nil, fmt.Errorf("load TLS certificate: %w", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{certificate},
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// generateSelfSignedCertificate creates a self-signed TLS certificate and key
// for localhost usage. Returns a [tls.Certificate].
func generateSelfSignedCertificate() (tls.Certificate, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("generate key: %w", err)
	}

	serialNumber, err := rand.Int(
		rand.Reader,
		new(big.Int).Lsh(big.NewInt(1), tlsSerialNumberBits),
	)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("generate serial: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"swag2mcp-mock"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(tlsCertificateValidity),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
		DNSNames:              []string{"localhost", "*.localhost"},
	}

	certificateBytes, certError := x509.CreateCertificate(
		rand.Reader,
		&template,
		&template,
		&privateKey.PublicKey,
		privateKey,
	)
	if certError != nil {
		return tls.Certificate{}, fmt.Errorf("create certificate: %w", certError)
	}

	certificatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificateBytes,
	})

	keyBytes, keyMarshalError := x509.MarshalECPrivateKey(privateKey)
	if keyMarshalError != nil {
		return tls.Certificate{}, fmt.Errorf("marshal key: %w", keyMarshalError)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: keyBytes,
	})

	return tls.X509KeyPair(certificatePEM, keyPEM)
}
