package cmd

import (
	"encoding/hex"
	"io/ioutil"
	"log"
	"time"

	"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	createSshLongDescription = "long desc"
)

func NewCreateSshCmd() *cobra.Command {
	var statePath string
	var deploymentName string

	var createSshCmd = &cobra.Command{
		Use:   "create-ssh",
		Short: "creates ssh configuration needed by subsequent stages",
		Long:  createSshLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting create-ssh command")

			if _, err := os.Stat(statePath); err == nil {
				state, err = ReadState(statePath)
				if err != nil {
					panic(err)
				}
			} else {
				log.Println("no state file, creating empty one")
			}

			// validate state
			// validate deploymentName
			// validate common state object

			state = RunCreateCommonCmd(deployProperties, state)

			err = WriteState(statePath, state)
			if err != nil {
				panic(err)
			}

			log.Println("finished create-ssh command")
		},
	}

	return createSshCmd
}

func RunCreateSshCmd(stateIn util.State, deploymentName string) (stateOut util.State) {
	stateOut = stateIn
}
