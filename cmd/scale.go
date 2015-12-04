package cmd

import (
	"log"

	//"github.com/colemickens/azkube/util"
	"github.com/spf13/cobra"
)

const (
	scaleLongDescription = "long desc"
)

func NewScaleCmd() *cobra.Command {
	var configPath string

	var scaleCmd = &cobra.Command{
		Use:   "scale",
		Short: "scale a kubernetes deployment",
		Long:  scaleLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("starting scale command")

			log.Println("finished scale command")
		},
	}

	scaleCmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config")

	return scaleCmd
}
