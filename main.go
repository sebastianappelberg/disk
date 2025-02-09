package main

import (
	"github.com/sebastianappelberg/disk/cmd"
	"log"
)

func main() {
	root := cmd.NewCmdRoot()
	err := root.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
