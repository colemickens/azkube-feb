package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"math/big"
	"net"
	"time"
)

const (
	ValidityDuration = time.Hour * 1
)

type PkiKeyCertPair struct {
	PrivateKey  *rsa.PrivateKey
	Certificate *x509.Certificate
}

func main() {
	pair, err := CreateCertificateAuthority()
	if err != nil {
		panic(err)
	}

	value, err := json.Marshal(pair)
	if err != nil {
		panic(err)
	}

	newPair := &PkiKeyCertPair{}

	err = json.Unmarshal(value, &newPair)
	if err != nil {
		panic(err)
	}
}

func CreateCertificateAuthority() (*PkiKeyCertPair, error) {
	var err error

	now := time.Now()

	ip := net.ParseIP("10.0.0.1")

	template := x509.Certificate{
		Subject:   pkix.Name{CommonName: "foo.com"},
		NotBefore: now,
		NotAfter:  now.Add(ValidityDuration),

		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,

		SerialNumber: big.NewInt(3),

		IPAddresses: []net.IP{ip},

		IsCA: true,
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)

	derCertificate, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(derCertificate)
	if err != nil {
		// this would be weird
		return nil, err
	}

	return &PkiKeyCertPair{Certificate: cert, PrivateKey: privateKey}, nil
}
