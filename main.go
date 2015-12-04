package main

import (
	"log"

	"github.com/colemickens/azkube/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		log.Fatalln(err)
	}
}
