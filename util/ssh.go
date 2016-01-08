package util

import (
	"crypto/rand"
	"crypto/rsa"
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

	privateKeyPem := PrivateKeyToPem(privateKey)

	sshProperties = &SshProperties{
		PrivateKeyPem: string(privateKeyPem),
	}

	return sshProperties, nil
}

func (s *SshProperties) OpenSshPublicKey() (string, error) {
	privateKey, err := PemToPrivateKey(s.PrivateKeyPem)
	if err != nil {
		return "", err
	}

	publicKey := privateKey.PublicKey

	sshPublicKey, err := ssh.NewPublicKey(&publicKey)
	if err != nil {
		return "", err
	}
	authorizedKeyBytes := ssh.MarshalAuthorizedKey(sshPublicKey)
	authorizedKey := string(authorizedKeyBytes)
	return authorizedKey, nil
}
