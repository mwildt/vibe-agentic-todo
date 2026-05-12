package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "vibe-cli",
		Short: "Vibe CLI tool for administrative tasks",
		Long:  "A command-line interface for managing the Vibe application",
		Run: func(cmd *cobra.Command, args []string) {
			// Show help if no subcommand is provided
			cmd.Help()
		},
	}

	// Add subcommands
	userCmd := &cobra.Command{
		Use:   "user",
		Short: "User management commands",
		Long:  "Commands for managing users in the YAML configuration",
	}
	
	userCmd.AddCommand(NewUserAddCommand())
	rootCmd.AddCommand(userCmd)

	return rootCmd
}

func checkFileExists(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Create the file if it doesn't exist
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		
		// Write initial empty YAML structure
		_, err = file.WriteString("users: []\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
