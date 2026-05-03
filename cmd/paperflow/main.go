package main

import (
	"os"

	"github.com/fcarvajalbrown/agentic-cs-paper-makers/internal/cli"
)

func main() {
	if err := cli.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
