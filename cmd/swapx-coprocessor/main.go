package main

import (
	"log"
	"os"

	"github.com/henriquemarlon/swapx/cmd/swapx-coprocessor/root"
)

func main() {
	if err := root.Cmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
		os.Exit(1)
	}
}