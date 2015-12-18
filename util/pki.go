package util

import (
	// "crypto/rsa"
	// "crypto/x509"
	// "crypto/x509/pkix"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
)

func (d *Deployer) GeneratePki() (pkiProperties *PkiProperties, err error) {
	// generate pki

	// generate kubeconfigs for each client
	// store those???

	// that's basically it

	return nil, nil
}

func createCertificateAuthority() (err error) {
	csr := &csr.CertificateRequest{
		CN:         "",
		Names:      []csr.Name{},
		Hosts:      []string{},
		KeyRequest: csr.NewBasicKeyRequest(),
		CA:         &csr.CAConfig{},
	}

	cert, csrPEM, key, err := initca.New(csr)

	_ = cert
	_ = csrPEM
	_ = key
	_ = err

	return nil
}
