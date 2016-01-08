package cmd

import (
	"crypto/x509"
	"log"
	"reflect"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	createPkiLongDescription = "long desc"
)

func NewCreatePkiCmd() *cobra.Command {
	var statePath string

	var createPkiCmd = &cobra.Command{
		Use:   "create-pki",
		Short: "creates the public key infrastructure to be used by the Kubernetes cluster",
		Long:  createPkiLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting create-pki command")

			state := &util.State{}
			var err error
			state, err = ReadAndValidateState(statePath,
				[]reflect.Type{
					reflect.TypeOf(state.Common),
				},
				[]reflect.Type{
					reflect.TypeOf(state.Pki),
					reflect.TypeOf(state.Myriad),
				},
			)
			if err != nil {
				panic(err)
			}

			RunCreatePkiCmd(state)

			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished create-pki command")
		},
	}

	createPkiCmd.Flags().StringVarP(&statePath, "state", "s", "./state.json", "path to load state from, and to persist state into")

	return createPkiCmd
}

func RunCreatePkiCmd(state *util.State) {
	pki := &util.PkiProperties{}
	var err error

	pki.CA, err = util.CreateCertificateAuthority(*state.Common)
	if err != nil {
		panic(err)
	}

	pki.ApiServer, err = util.CreateCertificate(*pki.CA, "kube-apiserver", x509.ExtKeyUsageServerAuth)
	if err != nil {
		panic(err)
	}

	pki.Kubelet, err = util.CreateCertificate(*pki.CA, "kube-kubelet", x509.ExtKeyUsageClientAuth)
	if err != nil {
		panic(err)
	}

	pki.Kubeproxy, err = util.CreateCertificate(*pki.CA, "kube-kubeproxy", x509.ExtKeyUsageClientAuth)
	if err != nil {
		panic(err)
	}

	pki.Scheduler, err = util.CreateCertificate(*pki.CA, "kube-scheduler", x509.ExtKeyUsageClientAuth)
	if err != nil {
		panic(err)
	}

	pki.ReplicationController, err = util.CreateCertificate(*pki.CA, "kube-replication-controller", x509.ExtKeyUsageClientAuth)
	if err != nil {
		panic(err)
	}

	state.Pki = pki
}
