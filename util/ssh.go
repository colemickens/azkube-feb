package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"

	"golang.org/x/crypto/ssh"
)

const (
	KeySize = 2048 // TODO : bump this back up when I'm not lazy to wait on it
)

func (d *Deployer) GenerateSsh() (sshProperties *SshProperties, err error) {
	log.Println("generating rsa key")

	privateKey, err := rsa.GenerateKey(rand.Reader, KeySize)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.PublicKey

	log.Println("generating openssh public key")
	sshPublicKey, err := ssh.NewPublicKey(&publicKey)
	if err != nil {
		return nil, err
	}
	authorizedKeyBytes := ssh.MarshalAuthorizedKey(sshPublicKey)
	authorizedKey := string(authorizedKeyBytes)

	log.Println("generating private key pem")
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	pemBuffer := &bytes.Buffer{}
	pem.Encode(pemBuffer, pemBlock)

	sshProperties = &SshProperties{
		OpenSshPublicKey: authorizedKey,
		PrivateKeyPem:    pemBuffer.Bytes(),
	}

	return sshProperties, nil
}
