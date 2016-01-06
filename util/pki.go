package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
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

func GeneratePki(common CommonProperties) (*PkiProperties, error) {
	pki := &PkiProperties{}
	var err error

	pki.CA, err = createCertificateAuthority(common)
	if err != nil {
		return nil, err
	}

	pki.ApiServer, err = createCertificate(*pki.CA, "kube-apiserver", x509.ExtKeyUsageServerAuth)
	if err != nil {
		return nil, err
	}

	pki.Kubelet, err = createCertificate(*pki.CA, "kube-kubelet", x509.ExtKeyUsageClientAuth)
	if err != nil {
		return nil, err
	}

	pki.Kubeproxy, err = createCertificate(*pki.CA, "kube-kubeproxy", x509.ExtKeyUsageClientAuth)
	if err != nil {
		return nil, err
	}

	pki.Scheduler, err = createCertificate(*pki.CA, "kube-scheduler", x509.ExtKeyUsageClientAuth)
	if err != nil {
		return nil, err
	}

	pki.ReplicationController, err = createCertificate(*pki.CA, "kube-replication-controller", x509.ExtKeyUsageClientAuth)
	if err != nil {
		return nil, err
	}

	return pki, nil
}

func createCertificateAuthority(common CommonProperties) (*PkiKeyCertPair, error) {
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

	snMax := new(big.Int).Lsh(big.NewInt(0), 128)
	template.SerialNumber, err = rand.Int(rand.Reader, snMax)
	if err != nil {
		return nil, err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, PkiKeySize)

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

func createCertificate(ca PkiKeyCertPair, subjectName string, usage x509.ExtKeyUsage) (*PkiKeyCertPair, error) {
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

	snMax := new(big.Int).Lsh(big.NewInt(0), 128)
	template.SerialNumber, err = rand.Int(rand.Reader, snMax)
	if err != nil {
		return nil, err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, PkiKeySize)

	derCertificate, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, ca.PrivateKey)
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

func (pair *PkiKeyCertPair) Kubeconfig(common CommonProperties) {
	// render the private key and cert into a kubeconfig file
	// requires `common` so that we know the master endpoint

}
