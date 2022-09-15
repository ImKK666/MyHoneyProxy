// Package certs implements a certificate signing authority implementation
// to sign MITM'ed hosts certificates using a self-signed authority.
//
// It has uses an LRU-based certificate caching implementation for
// caching the generated certificates for frequently accessed hosts.
package HoneyProxy

import (
	"bytes"
	"crypto/tls"
	"encoding/pem"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"
)

// Manager implements a certificate signing authority for TLS Mitm.
type Manager struct {
	cert  *tls.Certificate
	cache *lru.Cache
}

const (
	caKeyName  = "cakey.pem"
	caCertName = "cacert.pem"
)

var(
	gCertManager *Manager
)

// New creates a new certificate manager signing client instance
func New(cacheSize int) (*Manager, error) {
	manager := &Manager{}
	certFile := "./" + caCertName
	keyFile := "./" + caKeyName
	_, certFileErr := os.Stat(certFile)
	_, keyFileErr := os.Stat(keyFile)
	if os.IsNotExist(certFileErr) || os.IsNotExist(keyFileErr) {
		if err := manager.createAuthority(certFile, keyFile); err != nil {
			return nil, errors.Wrap(err, "could not create certificate authority")
		}
	}
retryRead:
	cert, err := manager.readCertificateDisk(certFile, keyFile)
	if err != nil {
		// Check if we have an expired cert and regenerate
		if err == errExpiredCert {
			if err := manager.createAuthority(certFile, keyFile); err != nil {
				return nil, errors.Wrap(err, "could not create certificate authority")
			}
			goto retryRead
		}
		return nil, errors.Wrap(err, "could not read certificate authority")
	}

	cache, err := lru.New(cacheSize)
	if err != nil {
		return nil, errors.Wrap(err, "could not create lru cache")
	}
	return &Manager{cert: cert, cache: cache}, nil
}

// GetCA returns the CA certificate in PEM Encoded format.
func (m *Manager) GetCA() (tls.Certificate, []byte) {
	buffer := &bytes.Buffer{}

	_ = pem.Encode(buffer, &pem.Block{Type: "CERTIFICATE", Bytes: m.cert.Certificate[0]})
	return *m.cert, buffer.Bytes()
}

// Get returns a certificate for the current host.
func (m *Manager) Get(host string) (*tls.Certificate, error) {
	if value, ok := m.cache.Get(host); ok {
		return value.(*tls.Certificate), nil
	}
	cert, err := m.signCertificate(host)
	if err != nil {
		return nil, err
	}
	m.cache.Add(host, cert)
	return cert, nil
}

func TLSConfigFromCA()func(host string, ctx *ProxyCtx) (*tls.Config, error)  {
	return gCertManager.tlsConfigFromCA()
}

// TLSConfigFromCA generates a spoofed TLS certificate for a host
func (m *Manager) tlsConfigFromCA() func(host string, ctx *ProxyCtx) (*tls.Config, error) {
	return func(host string, ctx *ProxyCtx) (c *tls.Config, err error) {
		hostname := stripPort(host)
		value, ok := m.cache.Get(host)
		if !ok {
			certificate, err := m.signCertificate(hostname)
			if err != nil {
				return nil, err
			}
			value = certificate
			m.cache.Add(host, certificate)
		}
		return &tls.Config{InsecureSkipVerify: true, Certificates: []tls.Certificate{*value.(*tls.Certificate)}}, nil
	}
}

func stripPort(s string) string {
	ix := strings.IndexRune(s, ':')
	if ix == -1 {
		return s
	}
	return s[:ix]
}

func init()  {
	var err error
	gCertManager,err = New(1024)
	if err != nil{
		log.Panicln("init certManager:",err)
	}
}