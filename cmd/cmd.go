package cmd

import (
	"github.com/spf13/cobra"
)

const (
	rootLongDescription = "longer description"
)

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "azkube",
		Short: "azure <-> kubernetes tool",
		Long:  rootLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	rootCmd.AddCommand(NewDeployCmd())
	rootCmd.AddCommand(NewScaleCmd())
	rootCmd.AddCommand(NewCertInstallCmd())

	return rootCmd
}
