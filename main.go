package main

import (
	"fmt"
	"hermetic/cmd"
	"os"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		fmt.Printf("failed to execute command, got error: '%s'\n", err)
		os.Exit(1)
	}
}
