package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	log "github.com/Sirupsen/logrus"
)

const (
	ValidityDuration = time.Hour * 24 * 365 * 2
	PkiKeySize       = 4096
)

// TODO(colemickens): could surely refactor/dedupe the two functions below, similar to x509's CreateCertificate
// TODO(Colemickens): potential options, duration, alt subj names, etc

func CreateKubeCertificates(masterFQDN string, extraFQDNs []string) (*PkiKeyCertPair, *PkiKeyCertPair, *PkiKeyCertPair, error) {
	log.Info("pki: generating certificate authority")
	ca, err := createCertificate("ca", "Azkube Certificate Authority", true, false, "")
	if err != nil {
		return nil, nil, nil, err
	}
	log.Info("pki: generating apiserver server certificate")
	apiserver, err := createCertificate("apiserver", "apiserver", false, true, masterFQDN, extraFQDNs...)
	if err != nil {
		return nil, nil, nil, err
	}
	log.Info("pki: generating client certificate")
	client, err := createCertificate("client", "client", false, false, "")
	if err != nil {
		return nil, nil, nil, err
	}

	return ca, apiserver, client, nil
}

func createCertificate(filenamePrefix string, commonName string, isCA bool, isServer bool, FQDN string, extraFQDNs ...string) (*PkiKeyCertPair, error) {
	var err error

	now := time.Now()

	template := x509.Certificate{
		Subject:   pkix.Name{CommonName: commonName},
		NotBefore: now,
		NotAfter:  now.Add(ValidityDuration),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	if isCA {
		template.KeyUsage |= x509.KeyUsageCertSign
		template.IsCA = true
	}
	if isServer {
		template.DNSNames = extraFQDNs
	}

	snMax := new(big.Int).Lsh(big.NewInt(1), 128)
	template.SerialNumber, err = rand.Int(rand.Reader, snMax)
	if err != nil {
		return nil, err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, PkiKeySize)

	certDerBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, err
	}

	certificatePem := CertificateToPem(certDerBytes)
	privateKeyPem := PrivateKeyToPem(privateKey)

	return &PkiKeyCertPair{
		CertificatePem: string(certificatePem),
		PrivateKeyPem:  string(privateKeyPem),
	}, nil
}
