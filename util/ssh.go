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
	SshKeySize = 4096
)

func GenerateSsh() (sshProperties *SshProperties, err error) {
	log.Println("generating rsa key")

	privateKey, err := rsa.GenerateKey(rand.Reader, SshKeySize)
	if err != nil {
		return nil, err
	}

	sshProperties = &SshProperties{
		PrivateKey: privateKey,
	}

	return sshProperties, nil
}

func (s *SshProperties) OpenSshPublicKey() (string, error) {
	publicKey := s.PrivateKey.PublicKey
	sshPublicKey, err := ssh.NewPublicKey(&publicKey)
	if err != nil {
		return "", err
	}
	authorizedKeyBytes := ssh.MarshalAuthorizedKey(sshPublicKey)
	authorizedKey := string(authorizedKeyBytes)
	return authorizedKey, nil
}

func (s *SshProperties) PrivateKeyPem() string {
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(s.PrivateKey),
	}
	pemBuffer := &bytes.Buffer{}
	pem.Encode(pemBuffer, pemBlock)
	pemString := string(pemBuffer.Bytes())
	return pemString
}
