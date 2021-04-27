package util

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	_ "crypto/sha512" // Needed for RegisterHash in init
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

// NewTLSCertPair returns a new PEM-encoded x.509 certificate pair based on a 521-bit ECDSA private key. The machine's
// local interface addresses and all variants of IPv4 and IPv6 localhost are included as valid IP addresses.
func NewTLSCertPair(organization string, validUntil time.Time, extraHosts []string) (cert, key []byte, e error) {
	now := time.Now()
	if validUntil.Before(now) {
		return nil, nil, errors.New("validUntil would create an already-expired certificate")
	}
	priv, e := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if e != nil {
		return nil, nil, e
	}
	// end of ASN.1 time
	endOfTime := time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC)
	if validUntil.After(endOfTime) {
		validUntil = endOfTime
	}
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, e := rand.Int(rand.Reader, serialNumberLimit)
	if e != nil {
		return nil, nil, fmt.Errorf("failed to generate serial number: %s", e)
	}
	host, e := os.Hostname()
	if e != nil {
		return nil, nil, e
	}
	ipAddresses := []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}
	dnsNames := []string{host}
	if host != "localhost" {
		dnsNames = append(dnsNames, "localhost")
	}
	addIP := func(ipAddr net.IP) {
		for _, ip := range ipAddresses {
			if ip.Equal(ipAddr) {
				return
			}
		}
		ipAddresses = append(ipAddresses, ipAddr)
	}
	addHost := func(host string) {
		for _, dnsName := range dnsNames {
			if host == dnsName {
				return
			}
		}
		dnsNames = append(dnsNames, host)
	}
	addrs, e := interfaceAddrs()
	if e != nil {
		return nil, nil, e
	}
	for _, a := range addrs {
		var ipAddr net.IP
		ipAddr, _, e = net.ParseCIDR(a.String())
		if e == nil {
			addIP(ipAddr)
		}
	}
	for _, hostStr := range extraHosts {
		host, _, e = net.SplitHostPort(hostStr)
		if e != nil {
			host = hostStr
		}
		if ip := net.ParseIP(host); ip != nil {
			addIP(ip)
		} else {
			addHost(host)
		}
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{organization},
			CommonName:   host,
		},
		NotBefore: now.Add(-time.Hour * 24),
		NotAfter:  validUntil,
		KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature |
			x509.KeyUsageCertSign,
		IsCA:                  true, // so can sign self.
		BasicConstraintsValid: true,
		DNSNames:              dnsNames,
		IPAddresses:           ipAddresses,
	}
	derBytes, e := x509.CreateCertificate(
		rand.Reader, &template,
		&template, &priv.PublicKey, priv,
	)
	if e != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %v", e)
	}
	certBuf := &bytes.Buffer{}
	e = pem.Encode(certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if e != nil {
		return nil, nil, fmt.Errorf("failed to encode certificate: %v", e)
	}
	keybytes, e := x509.MarshalECPrivateKey(priv)
	if e != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %v", e)
	}
	keyBuf := &bytes.Buffer{}
	e = pem.Encode(keyBuf, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keybytes})
	if e != nil {
		return nil, nil, fmt.Errorf("failed to encode private key: %v", e)
	}
	return certBuf.Bytes(), keyBuf.Bytes(), nil
}
