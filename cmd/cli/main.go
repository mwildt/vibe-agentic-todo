package main

import (
	"log"
	"os"
	"vibe-agentic/cmd/cli/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
