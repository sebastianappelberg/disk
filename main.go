package main

import (
	"log"

	"github.com/sebastianappelberg/disk/cmd"
)

func main() {
	root := cmd.NewCmdRoot()
	err := root.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
