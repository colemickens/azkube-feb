package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"net"
	"time"
)

const (
	ValidityDuration = time.Hour * 24 * 365 * 2
	PkiKeySize       = 4096
)

// TODO(colemickens): could surely refactor/dedupe the two functions below, similar to x509's CreateCertificate
// TODO(Colemickens): potential options, duration, alt subj names, etc

func CreateCertificateAuthority(common CommonProperties) (*PkiKeyCertPair, error) {
	var err error

	now := time.Now()

	template := x509.Certificate{
		Subject:   pkix.Name{CommonName: common.MasterFQDN},
		NotBefore: now,
		NotAfter:  now.Add(ValidityDuration),

		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,

		IPAddresses: []net.IP{common.MasterIP},

		IsCA: true,
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

func CreateCertificate(ca PkiKeyCertPair, subjectName string, usage x509.ExtKeyUsage) (*PkiKeyCertPair, error) {
	var err error

	now := time.Now()

	template := x509.Certificate{
		Subject:   pkix.Name{CommonName: subjectName},
		NotBefore: now,
		NotAfter:  now.Add(ValidityDuration),

		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{usage},
		BasicConstraintsValid: true,
	}

	snMax := new(big.Int).Lsh(big.NewInt(1), 128)
	template.SerialNumber, err = rand.Int(rand.Reader, snMax)
	if err != nil {
		return nil, err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, PkiKeySize)

	caPrivateKey, err := PemToPrivateKey(ca.PrivateKeyPem)
	if err != nil {
		log.Println("!!!! 1")
		return nil, err
	}

	certDerBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, caPrivateKey)
	if err != nil {
		log.Println("!!!! 2")
		return nil, err
	}

	certificatePem := CertificateToPem(certDerBytes)
	privateKeyPem := PrivateKeyToPem(privateKey)

	return &PkiKeyCertPair{
		CertificatePem: string(certificatePem),
		PrivateKeyPem:  string(privateKeyPem),
	}, nil
}

func (pair *PkiKeyCertPair) Kubeconfig(common CommonProperties) {
	// render the private key and cert into a kubeconfig file
	// requires `common` so that we know the master endpoint

}
