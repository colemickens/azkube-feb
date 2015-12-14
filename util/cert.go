package util

import ()

// TODO(colemickens): replace this with cfssl whenever pkcs12 output support lands

func (d *Deployer) GenerateSsh(destination string) (sshProperties *SshProperties, err error) {
	return nil, nil
}

func (d *Deployer) GeneratePki(destination string) (pkiProperties *PkiProperties, err error) {
	return nil, nil
}
